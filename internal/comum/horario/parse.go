package horario

import (
	"errors"
	"strings"
	"time"
)

func ParseISO8601(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, errors.New("vazio")
	}
	for _, layout := range []string{
		time.RFC3339Nano,
		time.RFC3339,
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		"2006-01-02",
	} {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), nil
		}
	}
	if !strings.ContainsAny(s, "Zz+-") && strings.Contains(s, "T") {
		if t, err := time.Parse(time.RFC3339, s+"Z"); err == nil {
			return t.UTC(), nil
		}
	}
	return time.Time{}, errors.New("formato nao reconhecido")
}
