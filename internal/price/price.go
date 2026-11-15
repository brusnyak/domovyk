// Package price implements the v0.2 price comparison step.
// It parses a FOCUS 1.2 billing export and compares AWS costs against
// OVHcloud and Hetzner price lists, with a checked_on date on every row.
package price

import (
	"fmt"
)

// Compare is a placeholder for v0.2. The real implementation parses a
// FOCUS 1.2 billing export and emits a cost table with checked_on dates.
func Compare(focusExportPath string) error {
	return fmt.Errorf("price comparison not implemented yet (v0.2): parse %s and emit cost table", focusExportPath)
}
