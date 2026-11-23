/*
 * Copyright 2026 Yegor
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Command domovyk is an AWS-to-EU cloud exit planner.
// Read-only against the customer's AWS account. Local-first.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/brusnyak/domovyk/internal/discover"
	"github.com/brusnyak/domovyk/internal/generate"
	"github.com/brusnyak/domovyk/internal/inventory"
	"github.com/brusnyak/domovyk/internal/price"
	"github.com/brusnyak/domovyk/internal/score"
)

const usage = `domovyk - AWS to EU exit planner (v0.1)

Usage:
  domovyk discover [--region REGION] [--profile PROFILE] [--out FILE]
                   Discover read-only AWS inventory into inventory JSON.
                   Credentials come from the AWS SDK default chain
                   (env, ~/.aws config, assume-role, IMDS); --profile is optional.
  domovyk score    [--inventory FILE] [--mapping FILE] [--out FILE]
                   Score each resource as portable / rework / rewrite.
  domovyk price    --focus FILE --out FILE [--prices FILE] [--fx USD=0.92]
                   Compare AWS billing against OVHcloud / Hetzner (v0.2).
  domovyk generate [--inventory FILE] [--score FILE] [--out DIR] [--target hetzner|ovhcloud]
                   Generate OpenTofu HCL for the target from scored inventory (v0.3).
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}

	sub := os.Args[1]

	switch sub {
	case "discover":
		runDiscover(os.Args[2:])
	case "score":
		runScore(os.Args[2:])
	case "price":
		runPrice(os.Args[2:])
	case "generate":
		runGenerate(os.Args[2:])
	case "-h", "--help", "help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n%s", sub, usage)
		os.Exit(2)
	}
}

func runDiscover(args []string) {
	fs := flag.NewFlagSet("discover", flag.ExitOnError)
	region := fs.String("region", "eu-central-1", "AWS region to inventory")
	profile := fs.String("profile", "", "AWS shared config profile")
	out := fs.String("out", "inventory.json", "output inventory JSON path")
	fs.Parse(args)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	inv, err := discover.Discover(ctx, discover.Options{
		Region:  *region,
		Profile: *profile,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "discover failed: %v\n", err)
		os.Exit(1)
	}

	if err := inv.Save(*out); err != nil {
		fmt.Fprintf(os.Stderr, "save inventory failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("discovered %d resources across %d services -> %s\n",
		len(inv.Resources), len(inv.Manifest.Services), *out)
	for svc, c := range inv.Manifest.Services {
		fmt.Printf("  %s: %d (attempted=%v completed=%v)", svc, c.ResourceCount, c.Attempted, c.Completed)
		if c.FirstError != "" {
			fmt.Printf(" first_error=%q", c.FirstError)
		}
		fmt.Println()
	}
}

func runScore(args []string) {
	fs := flag.NewFlagSet("score", flag.ExitOnError)
	inventoryPath := fs.String("inventory", "inventory.json", "input inventory JSON path")
	mappingPath := fs.String("mapping", "map/domovyk-map.yaml", "mapping data YAML path")
	out := fs.String("out", "score.json", "output score JSON path")
	fs.Parse(args)

	inv, err := inventory.Load(*inventoryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load inventory failed: %v\n", err)
		os.Exit(1)
	}

	rows, err := score.LoadMapping(*mappingPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load mapping failed: %v\n", err)
		os.Exit(1)
	}

	results, summary := score.Score(inv, rows)

	if err := score.WriteResults(results, summary, *out); err != nil {
		fmt.Fprintf(os.Stderr, "write results failed: %v\n", err)
		os.Exit(1)
	}

	data, _ := json.MarshalIndent(summary, "", "  ")
	fmt.Printf("scored %d resources -> %s\n%s\n", summary.Total, *out, string(data))
}

func runPrice(args []string) {
	fs := flag.NewFlagSet("price", flag.ExitOnError)
	focusPath := fs.String("focus", "", "FOCUS 1.2 billing export CSV for one calendar month (required)")
	pricesPath := fs.String("prices", "map/prices.yaml", "sourced EUR price table")
	mapPath := fs.String("map", "map/domovyk-map.yaml", "mapping table")
	out := fs.String("out", "", "write results JSON (mode 0600) (required)")
	fxFlag := fs.String("fx", "", "explicit currency-to-EUR rate, e.g. USD=0.92")
	fs.Parse(args)

	if *focusPath == "" || *out == "" {
		fmt.Fprintln(os.Stderr, "error: --focus and --out are required for price comparison")
		fmt.Print(usage)
		os.Exit(2)
	}
	fx, err := price.ParseFX(*fxFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "price failed: %v\n", err)
		os.Exit(1)
	}
	rows, err := price.Compare(*focusPath, *pricesPath, *mapPath, fx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "price failed: %v\n", err)
		os.Exit(1)
	}
	if err := price.WriteResults(*out, rows); err != nil {
		fmt.Fprintf(os.Stderr, "write results failed: %v\n", err)
		os.Exit(1)
	}
	if err := price.PrintTable(os.Stdout, rows); err != nil {
		fmt.Fprintf(os.Stderr, "print failed: %v\n", err)
		os.Exit(1)
	}
}

func runGenerate(args []string) {
	fs := flag.NewFlagSet("generate", flag.ExitOnError)
	inventoryPath := fs.String("inventory", "inventory.json", "input inventory JSON path")
	outDir := fs.String("out", "generated", "output directory for OpenTofu HCL")
	target := fs.String("target", "hetzner", "EU target: hetzner or ovhcloud")
	fs.Parse(args)

	inv, err := inventory.Load(*inventoryPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load inventory failed: %v\n", err)
		os.Exit(1)
	}

	rows, err := score.LoadMapping("map/domovyk-map.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "load mapping failed: %v\n", err)
		os.Exit(1)
	}

	results, _ := score.Score(inv, rows)

	plan, err := generate.Generate(inv, results, generate.Target(*target), *outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("generated %d servers, %d volumes, %d databases, %d buckets -> %s\n",
		len(plan.Servers), len(plan.Volumes), len(plan.Databases), len(plan.Buckets), *outDir)
	fmt.Println("Review and right-size the generated HCL before `tofu apply`.")
}
