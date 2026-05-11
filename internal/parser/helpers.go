package parser

import (
	"strconv"
	"strings"
)

func nullableString(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" || s == "N/A" {
		return nil
	}

	return &s
}

func parseIntOrZero(s string) int {
	s = strings.TrimSpace(s)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return v
}
