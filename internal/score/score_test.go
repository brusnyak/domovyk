	"strings"
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
