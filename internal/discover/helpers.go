package discover

import (
	"fmt"
	"time"
)

// strVal returns *s or "" if nil.
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// stringPtr returns a pointer to s.
func stringPtr(s string) *string {
	return &s
}

// boolVal formats *b or "" if nil.
func boolVal(b *bool) string {
	if b == nil {
		return ""
	}
	return fmt.Sprintf("%t", *b)
}

// strInt32 returns *int32 as int32 or 0 if nil.
func strInt32(i *int32) int32 {
	if i == nil {
		return 0
	}
	return *i
}

// timeVal formats *time.Time or "" if nil.
func timeVal(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
