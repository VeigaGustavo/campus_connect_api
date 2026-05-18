package horario

import "testing"

func TestParseISO8601(t *testing.T) {
	for _, s := range []string{
		"2026-08-31T23:59:59.000Z",
		"2026-08-31T23:59:59.000",
		"2026-06-30T00:00:00.000",
		"2026-08-31",
	} {
		if _, err := ParseISO8601(s); err != nil {
			t.Fatalf("parse %q: %v", s, err)
		}
	}
}
