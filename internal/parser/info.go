package parser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type swInfoEntry struct {
	Endianness             *int
	EnableEndiannessPerJob *int
	ReproducibilityDisable *int
}

func parseInfoFile(r io.Reader) (map[string]swInfoEntry, error) {
	result := map[string]swInfoEntry{}
	var currentGUID string
	var current swInfoEntry

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "---") {
			if currentGUID != "" {
				result[currentGUID] = current
				current = swInfoEntry{}
			}
			continue
		}

		if strings.HasPrefix(line, "SW_GUID=") {
			currentGUID = strings.ToLower(strings.TrimPrefix(line, "SW_GUID="))
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])

		v, err := strconv.Atoi(val)
		if err != nil {
			continue
		}

		switch key {
		case "endianness":
			current.Endianness = &v
		case "enable_endianness_per_job":
			current.EnableEndiannessPerJob = &v
		case "reproducibility_disable":
			current.ReproducibilityDisable = &v
		}
	}

	if currentGUID != "" {
		result[currentGUID] = current
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan info file: %w", err)
	}

	return result, nil
}
