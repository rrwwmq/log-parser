package parser

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/rrwwmq/log-parser/internal/core/domain"
)

type ParsedLog struct {
	Nodes []domain.Node
	Ports []domain.Port
}

func ParseZip(path string) (ParsedLog, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return ParsedLog{}, fmt.Errorf("open zip: %w", err)
	}
	defer r.Close()

	var logFile, infoFile io.ReadCloser

	for _, f := range r.File {
		name := strings.ToLower(f.Name)
		if strings.Contains(name, ".sharp_an_info") {
			rc, err := f.Open()
			if err != nil {
				return ParsedLog{}, fmt.Errorf("open file info in zip: %w", err)
			}
			infoFile = rc
		} else if strings.HasSuffix(name, ".db_csv") {
			rc, err := f.Open()
			if err != nil {
				return ParsedLog{}, fmt.Errorf("open log file in zip: %w", err)
			}
			logFile = rc
		}
	}

	if logFile == nil {
		return ParsedLog{}, fmt.Errorf("no log file found in zip")
	}
	defer logFile.Close()

	result, err := parseLog(logFile)
	if err != nil {
		return ParsedLog{}, fmt.Errorf("parse log: %w", err)
	}

	if infoFile != nil {
		defer infoFile.Close()

		swInfos, err := parseInfoFile(infoFile)
		if err != nil {
			return ParsedLog{}, fmt.Errorf("parse info file: %w", err)
		}

		for i, node := range result.Nodes {
			guidKey := strings.ToLower(strings.TrimPrefix(node.NodeGUID, "0x"))
			if sw, ok := swInfos[guidKey]; ok {
				if result.Nodes[i].Info == nil {
					result.Nodes[i].Info = &domain.NodeInfo{}
				}

				result.Nodes[i].Info.Endianness = sw.Endianness
				result.Nodes[i].Info.EnableEndiannessPerJob = sw.EnableEndiannessPerJob
				result.Nodes[i].Info.ReproducibilityDisable = sw.ReproducibilityDisable
			}
		}
	}

	return result, nil
}

func parseLog(r io.Reader) (ParsedLog, error) {
	sections := map[string][]string{}
	var currentSection string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "START_") {
			currentSection = strings.TrimPrefix(line, "START_")
			continue
		}

		if strings.HasPrefix(line, "END_") {
			currentSection = ""
			continue
		}

		if currentSection != "" && line != "" {
			sections[currentSection] = append(sections[currentSection], line)
		}
	}

	if err := scanner.Err(); err != nil {
		return ParsedLog{}, fmt.Errorf("scan log file: %w", err)
	}

	nodes, err := parseNodes(sections["NODES"])
	if err != nil {
		return ParsedLog{}, fmt.Errorf("parse nodes: %w", err)
	}

	sysInfos, err := parseSysInfo(sections["SYSTEM_GENERAL_INFORMATIONS"])
	if err != nil {
		return ParsedLog{}, fmt.Errorf("parse system info: %w", err)
	}

	guidToInfo := map[string]sysInfoEntry{}
	for _, info := range sysInfos {
		guidToInfo[info.NodeGUID] = info
	}

	for i, node := range nodes {
		if info, ok := guidToInfo[node.NodeGUID]; ok {
			nodes[i].Info = &domain.NodeInfo{
				SerialNumber: info.SerialNumber,
				PartNumber:   info.PartNumber,
				Revision:     info.Revision,
				ProductName:  info.ProductName,
			}
		}
	}

	ports, err := parsePorts(sections["PORTS"])
	if err != nil {
		return ParsedLog{}, fmt.Errorf("parse ports: %w", err)
	}

	return ParsedLog{
		Nodes: nodes,
		Ports: ports,
	}, nil
}
