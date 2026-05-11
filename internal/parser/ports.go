package parser

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

func parsePorts(lines []string) ([]domain.Port, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("ports section is empty or missing header")
	}

	var ports []domain.Port
	for _, line := range lines[1:] {
		r := csv.NewReader(strings.NewReader(line))
		fields, err := r.Read()
		if err != nil {
			return nil, fmt.Errorf("parse port line: %w", err)
		}

		if len(fields) < 7 {
			return nil, fmt.Errorf("port line has too few fields: %s", line)
		}

		portNum, err := strconv.Atoi(strings.TrimSpace(fields[2]))
		if err != nil {
			return nil, fmt.Errorf("parse port num: %w", err)
		}

		portState := parseIntOrZero(fields[20])
		lid := parseIntOrZero(fields[6])

		ports = append(ports, domain.NewUninitializedPort(
			domain.UninitializedID,
			strings.TrimSpace(fields[1]),
			portNum,
			portState,
			lid,
		))
	}

	return ports, nil
}
