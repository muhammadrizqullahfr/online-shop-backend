package models

import (
	"database/sql/driver"
	"encoding/csv"
	"fmt"
	"strings"
)

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if s == nil {
		return "{}", nil
	}
	escaped := make([]string, len(s))
	for i, v := range s {
		escaped[i] = `"` + strings.ReplaceAll(strings.ReplaceAll(v, `\`, `\\`), `"`, `\"`) + `"`
	}
	return "{" + strings.Join(escaped, ",") + "}", nil
}

func (s *StringSlice) Scan(src any) error {
	if src == nil {
		*s = nil
		return nil
	}

	var raw string
	switch v := src.(type) {
	case string:
		raw = v
	case []byte:
		raw = string(v)
	default:
		return fmt.Errorf("unsupported scan type for StringSlice: %T", src)
	}

	raw = strings.TrimSpace(raw)
	if raw == "" || raw == "{}" {
		*s = StringSlice{}
		return nil
	}
	if !strings.HasPrefix(raw, "{") || !strings.HasSuffix(raw, "}") {
		return fmt.Errorf("invalid postgres array literal: %q", raw)
	}

	reader := csv.NewReader(strings.NewReader(raw[1 : len(raw)-1]))
	reader.LazyQuotes = true
	parts, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to parse postgres array literal %q: %w", raw, err)
	}
	*s = StringSlice(parts)
	return nil
}
