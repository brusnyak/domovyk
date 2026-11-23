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

package generate

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brusnyak/domovyk/internal/inventory"
	"github.com/brusnyak/domovyk/internal/score"
)

func TestGenerateHetzner(t *testing.T) {
	inv := inventory.NewInventory("eu-central-1")
	inv.Resources = []inventory.Resource{
		{ID: "i-1", Service: "ec2", Type: "t3.medium", Region: "eu-central-1", State: "running", Details: map[string]string{"kind": "ec2_instance", "size_gb": "80"}},
		{ID: "vol-1", Service: "ec2", Type: "gp3", Region: "eu-central-1", State: "in-use", Details: map[string]string{"kind": "ebs_volume", "size_gb": "100"}},
		{ID: "db-1", Service: "rds", Type: "db.t3.micro", Region: "eu-central-1", State: "available", Details: map[string]string{"engine": "postgres"}},
		{ID: "b-1", Service: "s3", Type: "bucket", Region: "eu-central-1", State: "active", Details: map[string]string{}},
	}
	rows := []score.ScoreResult{
		{ResourceID: "i-1", Class: "portable", Source: "aws-ec2-instance"},
		{ResourceID: "vol-1", Class: "portable", Source: "aws-ec2-instance"},
		{ResourceID: "db-1", Class: "rework", Source: "aws-rds-postgresql"},
		{ResourceID: "b-1", Class: "rework", Source: "aws-s3-bucket"},
	}

	dir := t.TempDir()
	out := filepath.Join(dir, "gen")
	plan, err := Generate(inv, rows, Hetzner, out)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Servers) != 1 || plan.Servers[0].Name != "i-1" || plan.Servers[0].Type != "cx42" {
		t.Fatalf("unexpected server plan: %+v", plan.Servers)
	}
	if len(plan.Volumes) != 1 || plan.Volumes[0].Name != "vol-1" || plan.Volumes[0].SizeGB != 100 {
		t.Fatalf("unexpected volume plan: %+v", plan.Volumes)
	}
	if len(plan.Databases) != 1 || plan.Databases[0].Name != "db-1" {
		t.Fatalf("unexpected database plan: %+v", plan.Databases)
	}
	if len(plan.Buckets) != 1 || plan.Buckets[0].Name != "b-1" {
		t.Fatalf("unexpected bucket plan: %+v", plan.Buckets)
	}

	// Verify HCL files exist and reference hcloud provider.
	for _, f := range []string{"server_0.tf", "volume_0.tf", "database_0.tf", "bucket_0.tf", "plan.tf"} {
		content, err := os.ReadFile(filepath.Join(out, f))
		if err != nil {
			t.Fatalf("missing generated file %s: %v", f, err)
		}
		if f != "plan.tf" && !strings.Contains(string(content), "hcloud_") {
			t.Fatalf("%s does not reference hcloud provider: %s", f, string(content))
		}
	}
}

func TestGenerateOVH(t *testing.T) {
	inv := inventory.NewInventory("eu-central-1")
	inv.Resources = []inventory.Resource{
		{ID: "i-1", Service: "ec2", Type: "t3.medium", Region: "eu-central-1", State: "running", Details: map[string]string{"kind": "ec2_instance"}},
	}
	rows := []score.ScoreResult{
		{ResourceID: "i-1", Class: "portable", Source: "aws-ec2-instance"},
	}

	dir := t.TempDir()
	out := filepath.Join(dir, "gen")
	plan, err := Generate(inv, rows, OVHcloud, out)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Servers) != 1 || plan.Servers[0].Type != "b3-32" || plan.Servers[0].Image != "ubuntu-2204" {
		t.Fatalf("unexpected ovh server plan: %+v", plan.Servers)
	}
	content, err := os.ReadFile(filepath.Join(out, "server_0.tf"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "ovh_instance") {
		t.Fatalf("server HCL does not reference ovh provider: %s", string(content))
	}
}

func TestGenerateRejectsInvalidTarget(t *testing.T) {
	inv := inventory.NewInventory("eu-central-1")
	_, err := Generate(inv, nil, Target("aws"), t.TempDir())
	if err == nil || !strings.Contains(err.Error(), "unsupported target") {
		t.Fatalf("expected unsupported target error, got %v", err)
	}
}

func TestGenerateSkipsRewriteAndUnknown(t *testing.T) {
	inv := inventory.NewInventory("eu-central-1")
	inv.Resources = []inventory.Resource{
		{ID: "i-1", Service: "ec2", Type: "t3.medium", Region: "eu-central-1", State: "running", Details: map[string]string{"kind": "ec2_instance"}},
		{ID: "f-1", Service: "lambda", Type: "python3.12", Region: "eu-central-1", State: "active", Details: map[string]string{}},
		{ID: "g-1", Service: "glue", Type: "job", Region: "eu-central-1", State: "active", Details: map[string]string{}},
	}
	rows := []score.ScoreResult{
		{ResourceID: "i-1", Class: "portable", Source: "aws-ec2-instance"},
		{ResourceID: "f-1", Class: "rewrite", Source: "aws-lambda"},
		{ResourceID: "g-1", Class: "unknown", Source: "unmapped"},
	}

	dir := t.TempDir()
	out := filepath.Join(dir, "gen")
	plan, err := Generate(inv, rows, Hetzner, out)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Servers) != 1 {
		t.Fatalf("expected 1 server (rewrite/unknown skipped), got %d", len(plan.Servers))
	}
}
