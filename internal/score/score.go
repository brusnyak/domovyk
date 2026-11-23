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

package score

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/brusnyak/domovyk/internal/inventory"
	"gopkg.in/yaml.v3"
)

// ValidClass is the set of allowed mapping classes.
var ValidClass = map[string]bool{
	"portable": true,
	"rework":   true,
	"rewrite":  true,
}

// MappingRow is one row from map/domovyk-map.yaml.
type MappingRow struct {
	ID               string `yaml:"id"`
	AWSService       string `yaml:"aws_service"`
	Class            string `yaml:"class"`
	Caveat           string `yaml:"caveat"`
	CheckedOn        string `yaml:"checked_on"`
	Status           string `yaml:"status"`
	TypicalEUTarget  string `yaml:"typical_eu_target"`
}

// MappingFile is the top-level structure of domovyk-map.yaml.
type MappingFile struct {
	Rows []MappingRow `yaml:"rows"`
}

// MappingError records a schema/validation problem found in the mapping file.
type MappingError struct {
	Row   int
	Field string
	Msg   string
}

func (e *MappingError) Error() string {
	return fmt.Sprintf("mapping row %d field %s: %s", e.Row, e.Field, e.Msg)
}

// ScoreResult is the scoring output for one resource.
type ScoreResult struct {
	ResourceID string `json:"resource_id"`
	Service    string `json:"service"`
	Type       string `json:"type"`
	Region     string `json:"region"`
	Class      string `json:"class"`
	Caveat     string `json:"caveat,omitempty"`
	Source     string `json:"source"` // mapping row id or "unmapped"
}

// Summary aggregates scoring results by class.
type Summary struct {
	Portable int `json:"portable"`
	Rework   int `json:"rework"`
	Rewrite  int `json:"rewrite"`
	Unknown  int `json:"unknown"`
	Total    int `json:"total"`
}

// LoadMapping reads map/domovyk-map.yaml, validates each row against the schema,
// and returns rows indexed by AWS service prefix.
func LoadMapping(path string) ([]MappingRow, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read mapping file: %w", err)
	}
	var mf MappingFile
	if err := yaml.Unmarshal(data, &mf); err != nil {
		return nil, fmt.Errorf("parse mapping yaml: %w", err)
	}
	if err := validateMapping(mf.Rows); err != nil {
		return nil, err
	}
	return mf.Rows, nil
}

// validateMapping checks every row for required fields and valid enum values.
func validateMapping(rows []MappingRow) error {
	seen := map[string]bool{}
	for i, r := range rows {
		row := i + 1
		if strings.TrimSpace(r.ID) == "" {
			return &MappingError{Row: row, Field: "id", Msg: "required"}
		}
		if seen[r.ID] {
			return &MappingError{Row: row, Field: "id", Msg: "duplicate id"}
		}
		seen[r.ID] = true
		if strings.TrimSpace(r.AWSService) == "" {
			return &MappingError{Row: row, Field: "aws_service", Msg: "required"}
		}
		if !ValidClass[r.Class] {
			return &MappingError{Row: row, Field: "class", Msg: fmt.Sprintf("must be one of portable, rework, rewrite; got %q", r.Class)}
		}
		if strings.TrimSpace(r.Caveat) == "" {
			return &MappingError{Row: row, Field: "caveat", Msg: "required"}
		}
	}
	return nil
}

// classify matches a resource to a mapping row. Exact service-name match takes
// priority over substring matches; longer substring matches win.
func classify(res inventory.Resource, rows []MappingRow) ScoreResult {
	result := ScoreResult{
		ResourceID: res.ID,
		Service:    res.Service,
		Type:       res.Type,
		Region:     res.Region,
		Class:      "unknown",
		Source:     "unmapped",
	}

	// Exact service-name match takes priority.
	for _, row := range rows {
		if strings.EqualFold(row.AWSService, res.Service) {
			result.Class = row.Class
			result.Caveat = row.Caveat
			result.Source = row.ID
			return result
		}
	}

	bestLen := 0
	for _, row := range rows {
		serviceLower := strings.ToLower(row.AWSService)
		resServiceLower := strings.ToLower(res.Service)

		match := false
		matchLen := 0
		if strings.Contains(serviceLower, resServiceLower) {
			match = true
			matchLen = len(resServiceLower)
		} else if strings.Contains(resServiceLower, serviceLower) {
			match = true
			matchLen = len(serviceLower)
		}

		// Prefer longer, more specific matches.
		if match && matchLen > bestLen {
			bestLen = matchLen
			result.Class = row.Class
			result.Caveat = row.Caveat
			result.Source = row.ID
		}
	}

	return result
}

// Score classifies every resource in the inventory against the mapping data.
func Score(inv *inventory.Inventory, rows []MappingRow) ([]ScoreResult, Summary) {
	results := make([]ScoreResult, 0, len(inv.Resources))
	summary := Summary{}

	for _, res := range inv.Resources {
		r := classify(res, rows)
		results = append(results, r)

		switch r.Class {
		case "portable":
			summary.Portable++
		case "rework":
			summary.Rework++
		case "rewrite":
			summary.Rewrite++
		default:
			summary.Unknown++
		}
		summary.Total++
	}

	// Sort results for stable output: unknown first, then rewrite, rework, portable.
	sort.Slice(results, func(i, j int) bool {
		order := map[string]int{"unknown": 0, "rewrite": 1, "rework": 2, "portable": 3}
		return order[results[i].Class] < order[results[j].Class]
	})

	return results, summary
}

// WriteResults writes scoring results as JSON to path.
func WriteResults(results []ScoreResult, summary Summary, path string) error {
	out := struct {
		Summary  Summary      `json:"summary"`
		Results  []ScoreResult `json:"results"`
	}{
		Summary: summary,
		Results: results,
	}
	data, err := json.MarshalIndent(out, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o600)
}
