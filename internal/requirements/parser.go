package requirements

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	headingPattern       = regexp.MustCompile(`^#{1,6}\s+(.+?)\s*$`)
	requirementIDPattern = regexp.MustCompile(`^(PRD|TRD)-[A-Z0-9]+-[0-9]+[a-z]?$`)
)

func parseSourceDocument(workspaceRoot, sourceDoc, typeFilter, requirementID string) ([]requirementRow, error) {
	filter := strings.TrimSpace(typeFilter)
	if filter != "" && filter != "PRD" && filter != "TRD" {
		return nil, &UnknownTypeError{Type: filter}
	}

	path := sourceDocPath(workspaceRoot, sourceDoc)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer closeSourceFile(file)

	var (
		currentSection *string
		rows           []requirementRow
		seen           = make(map[string]struct{})
	)

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if heading := parseHeading(line); heading != nil {
			currentSection = heading
			continue
		}

		cells, ok := parseTableRow(line)
		if !ok || len(cells) < 2 {
			continue
		}
		if !requirementIDPattern.MatchString(cells[0]) {
			continue
		}
		if _, exists := seen[cells[0]]; exists {
			return nil, &DuplicateIDError{ID: cells[0]}
		}
		seen[cells[0]] = struct{}{}

		rowType := requirementType(cells[0])
		if filter != "" && rowType != filter {
			continue
		}
		if strings.TrimSpace(requirementID) != "" && cells[0] != strings.TrimSpace(requirementID) {
			continue
		}

		rows = append(rows, requirementRow{
			ID:             cells[0],
			Type:           rowType,
			Summary:        cells[1],
			SourceDoc:      filepath.ToSlash(filepath.Clean(strings.TrimSpace(sourceDoc))),
			SourceSection:  cloneOptionalString(currentSection),
			LastSyncedHash: hashSummary(cells[1]),
			Status:         "active",
		})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if trimmedID := strings.TrimSpace(requirementID); trimmedID != "" && len(rows) == 0 {
		return nil, &RequirementNotFoundInSourceError{ID: trimmedID, SourceDoc: strings.TrimSpace(sourceDoc)}
	}
	if len(rows) == 0 {
		return nil, &NoParseableIDsError{SourceDoc: strings.TrimSpace(sourceDoc)}
	}
	return rows, nil
}

func sourceDocPath(workspaceRoot, sourceDoc string) string {
	if filepath.IsAbs(sourceDoc) {
		return sourceDoc
	}
	return filepath.Join(workspaceRoot, sourceDoc)
}

func parseHeading(line string) *string {
	matches := headingPattern.FindStringSubmatch(strings.TrimSpace(line))
	if len(matches) != 2 {
		return nil
	}
	heading := strings.TrimSpace(matches[1])
	if heading == "" {
		return nil
	}
	return &heading
}

func parseTableRow(line string) ([]string, bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, "|") || !strings.HasSuffix(trimmed, "|") {
		return nil, false
	}

	parts := strings.Split(trimmed, "|")
	if len(parts) < 3 {
		return nil, false
	}

	cells := make([]string, 0, len(parts)-2)
	for _, part := range parts[1 : len(parts)-1] {
		cells = append(cells, strings.TrimSpace(part))
	}
	return cells, true
}

func requirementType(requirementID string) string {
	if strings.HasPrefix(requirementID, "TRD-") {
		return "TRD"
	}
	return "PRD"
}

func hashSummary(summary string) string {
	sum := sha256.Sum256([]byte(summary))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func cloneOptionalString(value *string) *string {
	if value == nil {
		return nil
	}
	copy := *value
	return &copy
}

func closeSourceFile(file *os.File) {
	_ = file.Close()
}
