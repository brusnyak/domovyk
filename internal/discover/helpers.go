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
