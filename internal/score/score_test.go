package score

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/brusnyak/domovyk/internal/inventory"
)

func TestLoadMapping(t *testing.T) {
	content := []byte("rows:\n  - id: aws-ec2-instance\n    aws_service: Amazon EC2 (instances)\n    class: portable\n    caveat: x\n    checked_on: \"\"\n    status: unverified\n  - id: aws-s3-bucket\n    aws_service: Amazon S3 (Standard buckets)\n    class: rework\n    caveat: y\n    checked_on: \"\"\n    status: unverified\n")
	dir := t.TempDir()
	path := filepath.Join(dir, "map.yaml")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatal(err)
	}
	rows, err := LoadMapping(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(rows))
	}
	if rows[0].ID != "aws-ec2-instance" || rows[0].Class != "portable" {
		t.Fatalf("row 0 unexpected: %+v", rows[0])
	}
	if rows[1].ID != "aws-s3-bucket" || rows[1].Class != "rework" {
		t.Fatalf("row 1 unexpected: %+v", rows[1])
	}
}

func TestLoadMappingRejectsInvalidRows(t *testing.T) {
	cases := []struct {
		name  string
		row   string
		want  string
	}{
		{"missing id", "  - aws_service: Amazon EC2\n    class: portable\n    caveat: x\n", "mapping row 1 field id"},
		{"missing service", "  - id: aws-ec2\n    class: portable\n    caveat: x\n", "mapping row 1 field aws_service"},
		{"invalid class", "  - id: aws-ec2\n    aws_service: Amazon EC2\n    class: bogus\n    caveat: x\n", "mapping row 1 field class"},
		{"missing caveat", "  - id: aws-ec2\n    aws_service: Amazon EC2\n    class: portable\n", "mapping row 1 field caveat"},
		{"duplicate id", "  - id: aws-ec2\n    aws_service: Amazon EC2\n    class: portable\n    caveat: x\n  - id: aws-ec2\n    aws_service: Amazon EC2\n    class: portable\n    caveat: y\n", "mapping row 2 field id"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			content := []byte("rows:\n" + tc.row)
			dir := t.TempDir()
			path := filepath.Join(dir, "map.yaml")
			if err := os.WriteFile(path, content, 0o644); err != nil {
				t.Fatal(err)
			}
			_, err := LoadMapping(path)
			if err == nil {
				t.Fatalf("expected validation error for %s", tc.name)
			}
			if !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %q, want substring %q", err.Error(), tc.want)
			}
		})
	}
}

func TestClassify(t *testing.T) {
	rows := []MappingRow{
		{ID: "aws-ec2-instance", AWSService: "Amazon EC2 (instances)", Class: "portable", Caveat: "caveat ec2"},
		{ID: "aws-s3-bucket", AWSService: "Amazon S3 (Standard buckets)", Class: "rework", Caveat: "caveat s3"},
		{ID: "aws-lambda", AWSService: "AWS Lambda", Class: "rewrite", Caveat: "caveat lambda"},
	}

	tests := []struct {
		name     string
		res      inventory.Resource
		want     string
		wantSrc  string
		wantCave bool
	}{
		{"ec2 instance", inventory.Resource{ID: "i-1", Service: "ec2", Type: "t3.medium"}, "portable", "aws-ec2-instance", true},
		{"s3 bucket", inventory.Resource{ID: "b-1", Service: "s3", Type: "bucket"}, "rework", "aws-s3-bucket", true},
		{"lambda fn", inventory.Resource{ID: "f-1", Service: "lambda", Type: "python3.12"}, "rewrite", "aws-lambda", true},
		{"unknown service", inventory.Resource{ID: "x-1", Service: "glue", Type: "job"}, "unknown", "unmapped", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := classify(tt.res, rows)
			if r.Class != tt.want {
				t.Errorf("class = %q, want %q", r.Class, tt.want)
			}
			if r.Source != tt.wantSrc {
				t.Errorf("source = %q, want %q", r.Source, tt.wantSrc)
			}
			if tt.wantCave && r.Caveat == "" {
				t.Error("expected a caveat, got empty")
			}
			if !tt.wantCave && r.Caveat != "" {
				t.Error("expected no caveat for unmapped, got " + r.Caveat)
			}
		})
	}
}

func TestClassifyExactMatchWins(t *testing.T) {
	rows := []MappingRow{
		{ID: "aws-ec2-instance", AWSService: "ec2", Class: "portable", Caveat: "exact"},
		{ID: "aws-everything", AWSService: "Amazon EC2 (everything)", Class: "rewrite", Caveat: "substring"},
	}
	res := inventory.Resource{ID: "i-1", Service: "ec2", Type: "t3.medium"}
	r := classify(res, rows)
	if r.Class != "portable" || r.Source != "aws-ec2-instance" {
		t.Errorf("expected exact match to win: got class=%q source=%q", r.Class, r.Source)
	}
}

func TestScoreSummary(t *testing.T) {
	rows := []MappingRow{
		{ID: "aws-ec2-instance", AWSService: "Amazon EC2 (instances)", Class: "portable"},
		{ID: "aws-s3-bucket", AWSService: "Amazon S3 (Standard buckets)", Class: "rework"},
		{ID: "aws-lambda", AWSService: "AWS Lambda", Class: "rewrite"},
	}
	inv := inventory.NewInventory("eu-central-1")
	inv.Resources = []inventory.Resource{
		{ID: "i-1", Service: "ec2", Type: "t3.medium"},
		{ID: "b-1", Service: "s3", Type: "bucket"},
		{ID: "f-1", Service: "lambda", Type: "python3.12"},
		{ID: "g-1", Service: "glue", Type: "job"},
	}
	results, summary := Score(inv, rows)
	if len(results) != 4 {
		t.Fatalf("expected 4 results, got %d", len(results))
	}
	if summary.Portable != 1 || summary.Rework != 1 || summary.Rewrite != 1 || summary.Unknown != 1 || summary.Total != 4 {
		t.Errorf("summary = %+v", summary)
	}
	// stable order: unknown, rewrite, rework, portable
	order := []string{"unknown", "rewrite", "rework", "portable"}
	got := []string{results[0].Class, results[1].Class, results[2].Class, results[3].Class}
	for i := range order {
		if got[i] != order[i] {
			t.Errorf("order[%d] = %q, want %q", i, got[i], order[i])
		}
	}
}
