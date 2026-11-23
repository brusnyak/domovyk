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

// Package generate turns a scored inventory into OpenTofu HCL modules for EU
// targets (Hetzner Cloud, OVHcloud Public Cloud). It is a transliteration,
// not a verified capacity match: instance types and disk sizes are best-effort
// approximations and must be right-sized by the operator before apply.
package generate

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/brusnyak/domovyk/internal/inventory"
	"github.com/brusnyak/domovyk/internal/score"
)

// Target is the EU provider to generate OpenTofu for.
type Target string

const (
	Hetzner  Target = "hetzner"
	OVHcloud Target = "ovhcloud"
)

// validTarget reports whether t is a supported generation target.
func validTarget(t Target) bool {
	switch t {
	case Hetzner, OVHcloud:
		return true
	}
	return false
}

// instanceMap maps an AWS instance type to the closest Hetzner/OVHcloud type
// plus a monthly price (EUR, excl. VAT, from map/prices.yaml checked 2026-11-23).
// Price is advisory only; verify before relying on it.
type instanceMap struct {
	target   string
	priceEUR float64
}

// hetznerTypes maps AWS vCPU/RAM class to Hetzner Cloud shared vCPU types.
var hetznerTypes = map[string]instanceMap{
	"nano":   {"cx11", 4.59},
	"micro":  {"cx22", 12.79},
	"small":  {"cx32", 23.49},
	"medium": {"cx42", 40.49},
	"large":  {"cx52", 79.49},
	"xlarge": {"cx62", 157.49},
}

// ovhTypes maps AWS vCPU/RAM class to OVHcloud b3-series (shared vCPU).
var ovhTypes = map[string]instanceMap{
	"nano":   {"b3-4", 16.68},
	"micro":  {"b3-8", 33.36},
	"small":  {"b3-16", 66.65},
	"medium": {"b3-32", 133.23},
	"large":  {"c3-16", 133.23},
	"xlarge": {"c3-32", 266.55},
}

// volumeClass maps EBS volume types to target block volume classes.
var volumeClass = map[string]string{
	"gp2":     "standard",
	"gp3":     "standard",
	"io1":     "high-performance",
	"io2":     "high-performance",
	"st1":     "standard",
	"sc1":     "standard",
	"standard": "standard",
}

// instanceClass buckets an EC2 instance type (e.g. "t3.medium") into a size class.
func instanceClass(instType string) string {
	parts := strings.SplitN(instType, ".", 2)
	if len(parts) != 2 {
		return "small"
	}
	size := parts[1]
	switch size {
	case "nano":
		return "nano"
	case "micro":
		return "micro"
	case "small":
		return "small"
	case "medium":
		return "medium"
	case "large":
		return "large"
	case "xlarge":
		return "xlarge"
	default:
		return "small"
	}
}

// targetInstance returns the target instance type and price for an AWS instance type and target.
func targetInstance(instType string, target Target) instanceMap {
	cls := instanceClass(instType)
	switch target {
	case Hetzner:
		if m, ok := hetznerTypes[cls]; ok {
			return m
		}
	case OVHcloud:
		if m, ok := ovhTypes[cls]; ok {
			return m
		}
	}
	// Fallback: cheapest shared instance.
	if target == Hetzner {
		return hetznerTypes["micro"]
	}
	return ovhTypes["micro"]
}

// Plan is the set of OpenTofu resources to generate.
type Plan struct {
	Servers   []Server   `json:"servers"`
	Volumes   []Volume   `json:"volumes"`
	Databases []Database `json:"databases"`
	Buckets   []Bucket   `json:"buckets"`
}

// Server is an OpenTofu compute instance.
type Server struct {
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Image    string  `json:"image"`
	Region   string  `json:"region"`
	PriceEUR float64 `json:"price_eur"`
	SourceID string  `json:"source_id"`
}

// Volume is an OpenTofu block volume.
type Volume struct {
	Name     string  `json:"name"`
	SizeGB   int     `json:"size_gb"`
	Region   string  `json:"region"`
	PriceEUR float64 `json:"price_eur"`
	SourceID string  `json:"source_id"`
}

// Database is an OpenTofu managed database.
type Database struct {
	Name     string `json:"name"`
	Engine   string `json:"engine"`
	Region   string `json:"region"`
	SourceID string `json:"source_id"`
}

// Bucket is an OpenTofu object-storage bucket.
type Bucket struct {
	Name     string `json:"name"`
	Region   string `json:"region"`
	SourceID string `json:"source_id"`
}

// Generate writes OpenTofu HCL for every portable/rework scored resource into outDir.
// rewrite resources are skipped (no direct equivalent). Returns the plan. The operator
// must review and right-size before `tofu apply`.
func Generate(inv *inventory.Inventory, results []score.ScoreResult, target Target, outDir string) (*Plan, error) {
	if !validTarget(target) {
		return nil, fmt.Errorf("unsupported target %q; use hetzner or ovhcloud", target)
	}
	if err := os.MkdirAll(outDir, 0o755); err != nil {
		return nil, err
	}

	byID := map[string]score.ScoreResult{}
	for _, r := range results {
		byID[r.ResourceID] = r
	}

	plan := &Plan{}
	serverIdx := 0
	volumeIdx := 0

	for _, res := range inv.Resources {
		r, ok := byID[res.ID]
		if !ok || r.Class == "rewrite" || r.Class == "unknown" {
			continue
		}
		switch res.Service {
		case "ec2":
			kind := res.Details["kind"]
			switch kind {
			case "ec2_instance":
				s := serverForEC2(res, target, r)
				plan.Servers = append(plan.Servers, s)
				if err := writeServerHCL(filepath.Join(outDir, fmt.Sprintf("server_%d.tf", serverIdx)), s, target); err != nil {
					return nil, err
				}
				serverIdx++
			case "ebs_volume":
				v := volumeForEBS(res, target, r)
				plan.Volumes = append(plan.Volumes, v)
				if err := writeVolumeHCL(filepath.Join(outDir, fmt.Sprintf("volume_%d.tf", volumeIdx)), v, target); err != nil {
					return nil, err
				}
				volumeIdx++
			}
		case "rds":
			d := databaseForRDS(res, target)
			plan.Databases = append(plan.Databases, d)
			if err := writeDatabaseHCL(filepath.Join(outDir, fmt.Sprintf("database_%d.tf", len(plan.Databases)-1)), d, target); err != nil {
				return nil, err
			}
		case "s3":
			b := bucketForS3(res)
			plan.Buckets = append(plan.Buckets, b)
			if err := writeBucketHCL(filepath.Join(outDir, fmt.Sprintf("bucket_%d.tf", len(plan.Buckets)-1)), b, target); err != nil {
				return nil, err
			}
		}
	}

	// Sort for stable output.
	sort.Slice(plan.Servers, func(i, j int) bool { return plan.Servers[i].Name < plan.Servers[j].Name })
	sort.Slice(plan.Volumes, func(i, j int) bool { return plan.Volumes[i].Name < plan.Volumes[j].Name })
	sort.Slice(plan.Databases, func(i, j int) bool { return plan.Databases[i].Name < plan.Databases[j].Name })
	sort.Slice(plan.Buckets, func(i, j int) bool { return plan.Buckets[i].Name < plan.Buckets[j].Name })

	if err := writePlanHCL(filepath.Join(outDir, "plan.tf"), plan); err != nil {
		return nil, err
	}
	return plan, nil
}

func serverForEC2(res inventory.Resource, target Target, r score.ScoreResult) Server {
	m := targetInstance(res.Type, target)
	image := "ubuntu-22.04"
	if target == OVHcloud {
		image = "ubuntu-2204"
	}
	return Server{
		Name:     res.ID,
		Type:     m.target,
		Image:    image,
		Region:   res.Region,
		PriceEUR: m.priceEUR,
		SourceID: res.ID,
	}
}

func volumeForEBS(res inventory.Resource, target Target, r score.ScoreResult) Volume {
	sizeGB := 20
	if v, ok := res.Details["size_gb"]; ok {
		fmt.Sscanf(v, "%d", &sizeGB)
	}
	return Volume{
		Name:     res.ID,
		SizeGB:   sizeGB,
		Region:   res.Region,
		PriceEUR: 0.044, // advisory GB/month, matches prices.yaml hetzner volumes row
		SourceID: res.ID,
	}
}

func databaseForRDS(res inventory.Resource, target Target) Database {
	engine := "postgres"
	if e := res.Details["engine"]; e != "" {
		engine = e
	}
	return Database{
		Name:     res.ID,
		Engine:   engine,
		Region:   res.Region,
		SourceID: res.ID,
	}
}

func bucketForS3(res inventory.Resource) Bucket {
	return Bucket{
		Name:     res.ID,
		Region:   res.Region,
		SourceID: res.ID,
	}
}

func writeServerHCL(path string, s Server, target Target) error {
	provider := providerName(target)
	content := fmt.Sprintf(`resource "%s_instance" "%s" {
  name  = %q
  type  = %q
  image = %q
  region = %q
  tags = {
    domovyk-ref = "true"
    source_id   = %q
  }
}
`, provider, s.Name, s.Name, s.Type, s.Image, s.Region, s.SourceID)
	return os.WriteFile(path, []byte(content), 0o600)
}

func writeVolumeHCL(path string, v Volume, target Target) error {
	provider := providerName(target)
	content := fmt.Sprintf(`resource "%s_volume" "%s" {
  name   = %q
  size   = %d
  region = %q
  tags = {
    domovyk-ref = "true"
    source_id   = %q
  }
}
`, provider, v.Name, v.Name, v.SizeGB, v.Region, v.SourceID)
	return os.WriteFile(path, []byte(content), 0o600)
}

func writeDatabaseHCL(path string, d Database, target Target) error {
	provider := providerName(target)
	content := fmt.Sprintf(`resource "%s_database" "%s" {
  name   = %q
  engine = %q
  region = %q
  tags = {
    domovyk-ref = "true"
    source_id   = %q
  }
}
`, provider, d.Name, d.Name, d.Engine, d.Region, d.SourceID)
	return os.WriteFile(path, []byte(content), 0o600)
}

func writeBucketHCL(path string, b Bucket, target Target) error {
	provider := providerName(target)
	content := fmt.Sprintf(`resource "%s_bucket" "%s" {
  name   = %q
  region = %q
  tags = {
    domovyk-ref = "true"
    source_id   = %q
  }
}
`, provider, b.Name, b.Name, b.Region, b.SourceID)
	return os.WriteFile(path, []byte(content), 0o600)
}

func writePlanHCL(path string, plan *Plan) error {
	var b strings.Builder
	b.WriteString("# Domovyk-generated OpenTofu plan. Review and right-size before `tofu apply`.\n\n")
	for _, s := range plan.Servers {
		b.WriteString(fmt.Sprintf("# %s -> %s (%s), advisory %.2f EUR/mo\n", s.Name, s.Type, s.Image, s.PriceEUR))
	}
	for _, v := range plan.Volumes {
		b.WriteString(fmt.Sprintf("# %s -> %d GB block volume, advisory %.2f EUR/mo\n", v.Name, v.SizeGB, v.PriceEUR))
	}
	for _, d := range plan.Databases {
		b.WriteString(fmt.Sprintf("# %s -> managed %s database\n", d.Name, d.Engine))
	}
	for _, bkt := range plan.Buckets {
		b.WriteString(fmt.Sprintf("# %s -> object storage bucket\n", bkt.Name))
	}
	b.WriteString("\n")
	return os.WriteFile(path, []byte(b.String()), 0o600)
}

func providerName(target Target) string {
	if target == Hetzner {
		return "hcloud"
	}
	return "ovh"
}
