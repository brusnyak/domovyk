package inventory

import (
	"encoding/json"
	"os"
	"time"
)

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

// Inventory is the read-only snapshot of an AWS account.
type Inventory struct {
	GeneratedAt string      `json:"generated_at"`
	AccountID   string      `json:"account_id,omitempty"`
	Region      string      `json:"region"`
	Resources   []Resource  `json:"resources"`
	Counts      map[string]int `json:"counts,omitempty"`
}

// NewInventory creates an empty inventory.
func NewInventory(region string) *Inventory {
	return &Inventory{
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
		Region:      region,
		Resources:   []Resource{},
		Counts:      map[string]int{},
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
