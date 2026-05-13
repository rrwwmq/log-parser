package parser

import (
	"encoding/csv"
	"fmt"
	"strconv"
	"strings"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

type sysInfoEntry struct {
	NodeGUID     string
	SerialNumber *string
	PartNumber   *string
	Revision     *string
	ProductName  *string
}

func parseNodes(lines []string) ([]domain.Node, error) {
	if len(lines) < 2 {
		return nil, fmt.Errorf("nodes section is empty or missing header")
	}

	var nodes []domain.Node
	for _, line := range lines[1:] {
		r := csv.NewReader(strings.NewReader(line))
		fields, err := r.Read()
		if err != nil {
			return nil, fmt.Errorf("parse node line: %w", err)
		}

		if len(fields) < 7 {
			return nil, fmt.Errorf("node line has too few fields: %s", line)
		}

		numPorts, err := strconv.Atoi(strings.TrimSpace(fields[1]))
		if err != nil {
			return nil, fmt.Errorf("parse num ports: %w", err)
		}

		nodeType, err := strconv.Atoi(strings.TrimSpace(fields[2]))
		if err != nil {
			return nil, fmt.Errorf("parse node type: %w", err)
		}

		nodes = append(nodes, domain.NewUninitializedNode(
			domain.UninitializedID,
			strings.TrimSpace(fields[6]),
			strings.TrimSpace(fields[0]),
			domain.NodeType(nodeType),
			numPorts,
		))
	}

	return nodes, nil
}

func parseSysInfo(lines []string) ([]sysInfoEntry, error) {
	if len(lines) < 2 {
		return nil, nil
	}

	var entries []sysInfoEntry
	for _, line := range lines[1:] {
		r := csv.NewReader(strings.NewReader(line))
		fields, err := r.Read()
		if err != nil {
			return nil, fmt.Errorf("parse sysInfo line: %w", err)
		}

		if len(fields) < 5 {
			return nil, fmt.Errorf("sysInfo line has too few fields: %s", line)
		}

		entries = append(entries, sysInfoEntry{
			NodeGUID:     strings.TrimSpace(fields[0]),
			SerialNumber: nullableString(fields[1]),
			PartNumber:   nullableString(fields[2]),
			Revision:     nullableString(fields[3]),
			ProductName:  nullableString(fields[4]),
		})
	}

	return entries, nil
}
