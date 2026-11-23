package inventory

import (
	"encoding/json"
	"os"
	"time"
)

// SchemaVersion is the inventory.json format version.
// Bump it when the shape of Inventory or Resource changes.
const SchemaVersion = "1"

// Resource is a single AWS resource discovered in the customer account.
type Resource struct {
	ID        string            `json:"id"`
	Service   string            `json:"service"`
	Type      string            `json:"type"`
	Region    string            `json:"region"`
	State     string            `json:"state"`
	Tags      map[string]string `json:"tags,omitempty"`
	Arn       string            `json:"arn,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}

// ServiceCoverage records what discovery observed for one AWS service.
type ServiceCoverage struct {
	// Attempted is true if discovery tried this service at least once.
	Attempted bool `json:"attempted"`
	// Completed is true if the service scan finished without a hard error.
	Completed bool `json:"completed"`
	// ResourceCount is how many resources were collected for this service.
	ResourceCount int `json:"resource_count"`
	// FirstError is the first hard error encountered, if any.
	FirstError string `json:"first_error,omitempty"`
	// PaginationComplete is true when all pages were read (false if scan was cut off).
	PaginationComplete bool `json:"pagination_complete"`
}

// Manifest records how, when, and how completely an inventory was produced.
type Manifest struct {
	// Schema is the inventory format version.
	Schema string `json:"schema"`
	// GeneratedAt is when the inventory was produced (RFC3339 UTC).
	GeneratedAt string `json:"generated_at"`
	// AccountID is the AWS account ID discovered, if sts:GetCallerIdentity succeeded.
	AccountID string `json:"account_id,omitempty"`
	// Region is the primary region scanned.
	Region string `json:"region"`
	// Regions is every region discovery was attempted against.
	Regions []string `json:"regions,omitempty"`
	// Services describes per-service scan coverage.
	Services map[string]ServiceCoverage `json:"services"`
	// TotalResources is the sum of all resources collected.
	TotalResources int `json:"total_resources"`
}

// Inventory is the read-only snapshot of an AWS account.
type Inventory struct {
	Manifest Manifest     `json:"manifest"`
	Resources []Resource  `json:"resources"`
}

// NewInventory creates an empty inventory with a versioned manifest.
func NewInventory(region string) *Inventory {
	return &Inventory{
		Manifest: Manifest{
			Schema:      SchemaVersion,
			GeneratedAt: time.Now().UTC().Format(time.RFC3339),
			Region:      region,
			Services:    map[string]ServiceCoverage{},
		},
		Resources: []Resource{},
	}
}

// Save writes the inventory to a JSON file with mode 0600.
func (i *Inventory) Save(path string) error {
	data, err := json.MarshalIndent(i, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}

// Load reads an inventory from a JSON file.
func Load(path string) (*Inventory, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var inv Inventory
	if err := json.Unmarshal(data, &inv); err != nil {
		return nil, err
	}
	return &inv, nil
}

// CountResources updates the per-service resource counts and totals.
func (i *Inventory) CountResources() {
	counts := map[string]int{}
	total := 0
	for _, r := range i.Resources {
		counts[r.Service]++
		total++
	}
	for svc, n := range counts {
		if c, ok := i.Manifest.Services[svc]; ok {
			c.ResourceCount = n
			i.Manifest.Services[svc] = c
		}
	}
	i.Manifest.TotalResources = total
}
