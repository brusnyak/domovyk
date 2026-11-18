// Command domovyk is an AWS-to-EU cloud exit planner.
// Read-only against the customer's AWS account. Local-first.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/brusnyak/domovyk/internal/discover"
	"github.com/brusnyak/domovyk/internal/inventory"
	"github.com/brusnyak/domovyk/internal/price"
	"github.com/brusnyak/domovyk/internal/score"
)

const usage = `domovyk - AWS to EU exit planner (v0.1)

Usage:
  domovyk discover [--region REGION] [--profile PROFILE] [--out FILE]
                   Discover read-only AWS inventory into inventory JSON.
  domovyk score    [--inventory FILE] [--mapping FILE] [--out FILE]
                   Score each resource as portable / rework / rewrite.
  domovyk price    --focus FILE --out FILE [--prices FILE] [--fx USD=0.92]
                   Compare AWS billing against OVHcloud / Hetzner (v0.2).
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
	accessKey := fs.String("access-key", "", "AWS access key ID (overrides profile/env)")
	secretKey := fs.String("secret-key", "", "AWS secret access key (overrides profile/env)")
	session := fs.String("session", "", "AWS session token")
	out := fs.String("out", "inventory.json", "output inventory JSON path")
	fs.Parse(args)

	inv, err := discover.Discover(discover.Options{
		Region:    *region,
		Profile:   *profile,
		AccessKey: *accessKey,
		SecretKey: *secretKey,
		Session:   *session,
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
		len(inv.Resources), len(inv.Counts), *out)
	for svc, n := range inv.Counts {
		fmt.Printf("  %s: %d\n", svc, n)
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
