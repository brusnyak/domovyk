package price

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Synthetic test rates, not actual provider quotes.
const testPrices = `rows:
  - provider: hetzner
    sku: synthetic-test-only
    category: compute
    unit: hour
    price_eur: 0.1
    source_url: https://example.invalid/synthetic-price
    checked_on: "2026-10-01"
    status: unverified
`

func writeTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestCompare(t *testing.T) {
	fixture, err := os.ReadFile("testdata/focus.csv")
	if err != nil {
		t.Fatal(err)
	}
	cases := []struct {
		name, prices, csv, fx, wantErr string
		wantCount                      int
		cost, estimate                 float64
		noPrice                        bool
	}{
		{name: "grouping and credits", prices: testPrices, wantCount: 2, cost: 90, estimate: 15},
		{name: "null unverified rate", prices: strings.Replace(testPrices, "price_eur: 0.1", "price_eur: null", 1), wantCount: 2, cost: 90, noPrice: true},
		{name: "missing category", prices: "rows: []\n", wantCount: 2, cost: 90, noPrice: true},
		{name: "USD requires FX", prices: testPrices, csv: strings.ReplaceAll(string(fixture), ",EUR,", ",USD,"), wantErr: "requires explicit --fx"},
		{name: "USD converted", prices: testPrices, csv: strings.ReplaceAll(string(fixture), ",EUR,", ",USD,"), fx: "USD=0.92", wantCount: 2, cost: 82.8, estimate: 15},
		{name: "incompatible units", prices: strings.Replace(testPrices, "unit: hour", "unit: month", 1), wantCount: 2, cost: 90, noPrice: true},
		{name: "multiple units", prices: testPrices, csv: strings.Replace(string(fixture), "50,Hours", "50,GB", 1), wantCount: 2, cost: 90, noPrice: true},
		{name: "mixed EC2 resources", prices: testPrices, csv: strings.Replace(string(fixture), "AWS::EC2::Instance", "AWS::EC2::Volume", 1), wantCount: 2, cost: 90, noPrice: true},
		{name: "multiple months", prices: testPrices, csv: strings.Replace(string(fixture), "2026-10-02", "2026-11-02", 1), wantErr: "multiple months"},
		{name: "invalid number", prices: testPrices, csv: strings.Replace(string(fixture), ",60,", ",NaN,", 1), wantErr: "BilledCost"},
		{name: "missing column", prices: testPrices, csv: strings.Replace(string(fixture), "ResourceId", "OtherColumn", 1), wantErr: "missing FOCUS column ResourceId"},
		{name: "malformed row", prices: testPrices, csv: string(fixture) + "Amazon EC2,Compute\n", wantErr: "FOCUS row"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			csv := tc.csv
			if csv == "" {
				csv = string(fixture)
			}
			focus := writeTestFile(t, dir, "focus.csv", csv)
			prices := writeTestFile(t, dir, "prices.yaml", tc.prices)
			fx, err := ParseFX(tc.fx)
			if err != nil {
				t.Fatal(err)
			}
			results, err := Compare(focus, prices, "../../map/domovyk-map.yaml", fx)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("error = %v, want %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(results) != tc.wantCount {
				t.Fatalf("rows = %+v", results)
			}
			first := results[0]
			if first.AWSService != "Amazon EC2" || first.AWSMonthly != tc.cost || first.Currency != "EUR" {
				t.Fatalf("group = %+v", first)
			}
			if tc.noPrice {
				if first.Status != "no_price" || first.EstimatedMonthly != nil || first.Delta != nil || first.CheckedOn != "not_available" {
					t.Fatalf("no_price = %+v", first)
				}
			} else if first.EstimatedMonthly == nil || *first.EstimatedMonthly != tc.estimate || first.Delta == nil || *first.Delta != tc.cost-tc.estimate {
				t.Fatalf("estimate = %+v", first)
			}
			if results[1].AWSService != "Amazon S3" || results[1].Status != "no_price" {
				t.Fatalf("missing category = %+v", results[1])
			}
			for _, row := range results {
				if row.CheckedOn == "" || row.SourceURL == "" {
					t.Fatalf("missing provenance: %+v", row)
				}
			}
			var table bytes.Buffer
			if err := PrintTable(&table, results); err != nil {
				t.Fatal(err)
			}
			if !tc.noPrice && !strings.Contains(table.String(), "UNVERIFIED") {
				t.Fatal(table.String())
			}
			output := writeTestFile(t, dir, "price.json", "old contents")
			if err := os.Chmod(output, 0644); err != nil {
				t.Fatal(err)
			}
			if err := WriteResults(output, results); err != nil {
				t.Fatal(err)
			}
			info, err := os.Stat(output)
			if err != nil {
				t.Fatal(err)
			}
			if info.Mode().Perm() != 0600 {
				t.Fatalf("permissions = %o", info.Mode().Perm())
			}
			data, err := os.ReadFile(output)
			if err != nil {
				t.Fatal(err)
			}
			var decoded []map[string]any
			if err := json.Unmarshal(data, &decoded); err != nil {
				t.Fatal(err)
			}
			for _, row := range decoded {
				if row["checked_on"] == "" || row["source_url"] == "" {
					t.Fatalf("JSON provenance: %v", row)
				}
			}
		})
	}
}

func TestPriceValidation(t *testing.T) {
	cases := []struct{ name, old, replacement string }{
		{"date", "2026-10-01", "2026-02-30"},
		{"source", "https://example.invalid/synthetic-price", "not-a-url"},
		{"provider", "provider: hetzner", "provider: unknown"},
		{"category", "category: compute", "category: unknown"},
		{"status", "status: verified", "status: unknown"},
		{"negative", "price_eur: 0.1", "price_eur: -1"},
		{"nan", "price_eur: 0.1", "price_eur: .nan"},
		{"missing verified rate", "price_eur: 0.1", "price_eur: null"},
		{"unknown field", "price_eur: 0.1", "price_eur: 0.1\n    typo: value"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			path := writeTestFile(t, t.TempDir(), "prices.yaml", strings.Replace(strings.Replace(testPrices, "status: unverified", "status: verified", 1), tc.old, tc.replacement, 1))
			if _, err := LoadPrices(path); err == nil {
				t.Fatal("expected invalid price error")
			}
		})
	}
}

func TestFXValidation(t *testing.T) {
	for _, input := range []string{"USD", "USD=0", "USD=-1", "USD=NaN", "USD=Inf", "usd=1", "EUR=2", "USD=1=2"} {
		t.Run(input, func(t *testing.T) {
			if _, err := ParseFX(input); err == nil {
				t.Fatal("expected invalid FX error")
			}
		})
	}
}

func TestMappingEligibility(t *testing.T) {
	mapping := writeTestFile(t, t.TempDir(), "map.yaml", "rows:\n  - id: aws-ec2-instance\n    aws_service: Amazon EC2\n    typical_eu_target: ''\n")
	prices := writeTestFile(t, t.TempDir(), "prices.yaml", testPrices)
	rows, err := Compare("testdata/focus.csv", prices, mapping, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 0 {
		t.Fatalf("empty target must be skipped: %+v", rows)
	}
}
