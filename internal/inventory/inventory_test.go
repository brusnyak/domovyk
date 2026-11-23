package inventory

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoadRoundtrip(t *testing.T) {
	inv := NewInventory("eu-central-1")
	inv.Resources = []Resource{
		{ID: "i-1", Service: "ec2", Type: "t3.medium", Region: "eu-central-1", State: "running"},
		{ID: "b-1", Service: "s3", Type: "bucket", Region: "eu-central-1", State: "active"},
	}
	inv.Manifest.Services = map[string]ServiceCoverage{
		"ec2_instance": {Attempted: true, Completed: true, ResourceCount: 1, PaginationComplete: true},
		"s3_bucket":    {Attempted: true, Completed: true, ResourceCount: 1, PaginationComplete: true},
	}
	inv.Manifest.TotalResources = 2
	inv.Manifest.AccountID = "123456789012"

	dir := t.TempDir()
	path := filepath.Join(dir, "inventory.json")
	if err := inv.Save(path); err != nil {
		t.Fatal(err)
	}

	loaded, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(loaded.Resources) != 2 {
		t.Fatalf("expected 2 resources, got %d", len(loaded.Resources))
	}
	if loaded.Resources[0].ID != "i-1" || loaded.Resources[0].State != "running" {
		t.Errorf("resource 0 unexpected: %+v", loaded.Resources[0])
	}
	if loaded.Manifest.Region != "eu-central-1" {
		t.Errorf("region = %q, want eu-central-1", loaded.Manifest.Region)
	}
	if loaded.Manifest.Schema != SchemaVersion {
		t.Errorf("schema = %q, want %q", loaded.Manifest.Schema, SchemaVersion)
	}
	if loaded.Manifest.AccountID != "123456789012" {
		t.Errorf("account_id = %q, want 123456789012", loaded.Manifest.AccountID)
	}
	if loaded.Manifest.TotalResources != 2 {
		t.Errorf("total_resources = %d, want 2", loaded.Manifest.TotalResources)
	}
	if loaded.Manifest.Services["ec2_instance"].Completed != true {
		t.Errorf("ec2_instance not marked completed")
	}

	// file mode should be 0600
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Errorf("file mode = %o, want 0600", info.Mode().Perm())
	}
}

func TestNewInventoryDefaults(t *testing.T) {
	inv := NewInventory("eu-west-1")
	if inv.Manifest.Schema != SchemaVersion {
		t.Errorf("schema = %q, want %q", inv.Manifest.Schema, SchemaVersion)
	}
	if inv.Manifest.Region != "eu-west-1" {
		t.Errorf("region = %q, want eu-west-1", inv.Manifest.Region)
	}
	if inv.Resources == nil {
		t.Error("resources should be initialized to empty slice")
	}
	if inv.Manifest.Services == nil {
		t.Error("services map should be initialized")
	}
}
