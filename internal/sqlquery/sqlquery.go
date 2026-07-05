package sqlquery

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/mbell-nzl/fid/internal/config"
	"github.com/mbell-nzl/fid/internal/journal"

	_ "modernc.org/sqlite"
)

const sqliteDriverName = "sqlite"

var requiredReadPragmas = []string{
	"foreign_keys(1)",
	"recursive_triggers(1)",
	"query_only(1)",
}

var readPragmasWithArgs = map[string]struct{}{
	"foreign_key_check": {},
	"foreign_key_list":  {},
	"function_list":     {},
	"index_info":        {},
	"index_list":        {},
	"index_xinfo":       {},
	"module_list":       {},
	"pragma_list":       {},
	"table_info":        {},
	"table_list":        {},
	"table_xinfo":       {},
}

type Options struct {
	WorkspaceRoot string
	MigrationsDir string
	Query         string
}

type Row map[string]any

type Result struct {
	Columns []string
	Rows    []Row
}

type RejectionError struct {
	Statement string
}

func (err *RejectionError) Error() string {
	if err == nil {
		return "<nil>"
	}
	return "fid sql accepts only SELECT, EXPLAIN, PRAGMA query, and WITH statements"
}

func Classify(query string) error {
	statements := splitStatements(query)
	if len(statements) == 0 {
		return &RejectionError{Statement: strings.TrimSpace(query)}
	}
	if len(statements) > 1 {
		return &RejectionError{Statement: statements[1]}
	}

	return classifyStatement(statements[0])
}

func Run(options Options) (Result, error) {
	workspaceRoot := strings.TrimSpace(options.WorkspaceRoot)
	if workspaceRoot == "" {
		return Result{}, fmt.Errorf("sql query: workspace root is required")
	}
	migrationsDir := strings.TrimSpace(options.MigrationsDir)
	if migrationsDir == "" {
		return Result{}, fmt.Errorf("sql query: migrations dir is required")
	}

	query := strings.TrimSpace(options.Query)
	if err := Classify(query); err != nil {
		return Result{}, err
	}

	if _, err := journal.EnsureProjectionFresh(workspaceRoot, migrationsDir); err != nil {
		return Result{}, err
	}

	conn, err := openReadOnly(filepath.Join(workspaceRoot, config.WorkspaceDirName, config.DatabaseFileName))
	if err != nil {
		return Result{}, err
	}
	defer closeSQLQueryConn(conn)

	rows, err := conn.Query(query)
	if err != nil {
		return Result{}, fmt.Errorf("sql query: %w", err)
	}
	defer closeSQLQueryRows(rows)

	return scanResult(rows)
}

func classifyStatement(statement string) error {
	tokens := tokenize(statement)
	if len(tokens) == 0 {
		return &RejectionError{Statement: statement}
	}

	switch tokens[0].upper {
	case "SELECT", "EXPLAIN":
		return nil
	case "PRAGMA":
		if pragmaReadOnly(tokens) {
			return nil
		}
	case "WITH":
		index, ok := withTerminalKeyword(tokens)
		if ok && index < len(tokens) {
			switch tokens[index].upper {
			case "SELECT", "EXPLAIN":
				return nil
			case "PRAGMA":
				if pragmaReadOnly(tokens[index:]) {
					return nil
				}
			}
			return &RejectionError{Statement: statement[tokens[index].start:]}
		}
	}

	return &RejectionError{Statement: statement}
}

func pragmaReadOnly(tokens []sqlToken) bool {
	name, nextIndex, ok := pragmaName(tokens)
	if !ok {
		return false
	}
	if nextIndex >= len(tokens) {
		return true
	}
	if tokens[nextIndex].text == "=" {
		return false
	}
	if tokens[nextIndex].text != "(" {
		return true
	}
	_, ok = readPragmasWithArgs[name]
	return ok
}

func pragmaName(tokens []sqlToken) (string, int, bool) {
	if len(tokens) < 2 || tokens[0].upper != "PRAGMA" || tokens[1].kind != tokenIdentifier {
		return "", 0, false
	}

	index := 2
	name := strings.ToLower(tokens[1].text)
	if index+1 < len(tokens) && tokens[index].text == "." && tokens[index+1].kind == tokenIdentifier {
		name = strings.ToLower(tokens[index+1].text)
		index += 2
	}
	return name, index, true
}

func withTerminalKeyword(tokens []sqlToken) (int, bool) {
	if len(tokens) == 0 || tokens[0].upper != "WITH" {
		return 0, false
	}

	index := 1
	if index < len(tokens) && tokens[index].upper == "RECURSIVE" {
		index++
	}

	for index < len(tokens) {
		nextIndex, ok := skipCTE(tokens, index)
		if !ok || nextIndex >= len(tokens) {
			return 0, false
		}
		if tokens[nextIndex].text == "," {
			index = nextIndex + 1
			continue
		}
		return nextIndex, true
	}

	return 0, false
}

func skipCTE(tokens []sqlToken, index int) (int, bool) {
	asIndex, ok := findKeywordAtDepth(tokens, index, "AS")
	if !ok {
		return 0, false
	}

	index = asIndex + 1
	if index < len(tokens) && tokens[index].upper == "NOT" {
		index++
	}
	if index < len(tokens) && tokens[index].upper == "MATERIALIZED" {
		index++
	}
	if index >= len(tokens) || tokens[index].text != "(" {
		return 0, false
	}

	return skipBalanced(tokens, index)
}

func findKeywordAtDepth(tokens []sqlToken, index int, keyword string) (int, bool) {
	depth := 0
	for index < len(tokens) {
		switch tokens[index].text {
		case "(":
			depth++
		case ")":
			if depth > 0 {
				depth--
			}
		default:
			if depth == 0 && tokens[index].upper == keyword {
				return index, true
			}
		}
		index++
	}
	return 0, false
}

func skipBalanced(tokens []sqlToken, index int) (int, bool) {
	depth := 0
	for index < len(tokens) {
		switch tokens[index].text {
		case "(":
			depth++
		case ")":
			depth--
			if depth == 0 {
				return index + 1, true
			}
		}
		index++
	}
	return 0, false
}

func openReadOnly(path string) (*sql.DB, error) {
	dsn, err := readOnlyDSN(path)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open(sqliteDriverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}
	return conn, nil
}

func readOnlyDSN(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve sqlite database path: %w", err)
	}

	dsnPath := filepath.ToSlash(absolutePath)
	if !strings.HasPrefix(dsnPath, "/") {
		dsnPath = "/" + dsnPath
	}

	query := url.Values{}
	query.Set("mode", "ro")
	for _, pragma := range requiredReadPragmas {
		query.Add("_pragma", pragma)
	}

	return (&url.URL{
		Scheme:   "file",
		Path:     dsnPath,
		RawQuery: query.Encode(),
	}).String(), nil
}

func scanResult(rows *sql.Rows) (Result, error) {
	columns, err := rows.Columns()
	if err != nil {
		return Result{}, fmt.Errorf("sql query: columns: %w", err)
	}

	items := make([]Row, 0)
	for rows.Next() {
		item, err := scanRow(rows, columns)
		if err != nil {
			return Result{}, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return Result{}, fmt.Errorf("sql query: iterate rows: %w", err)
	}

	return Result{Columns: columns, Rows: items}, nil
}

func scanRow(rows *sql.Rows, columns []string) (Row, error) {
	values := make([]any, len(columns))
	destinations := make([]any, len(columns))
	for index := range values {
		destinations[index] = &values[index]
	}

	if err := rows.Scan(destinations...); err != nil {
		return nil, fmt.Errorf("sql query: scan row: %w", err)
	}

	item := make(Row, len(columns))
	for index, column := range columns {
		item[column] = normalizeValue(values[index])
	}
	return item, nil
}

func normalizeValue(value any) any {
	if bytes, ok := value.([]byte); ok {
		return string(bytes)
	}
	return value
}

func splitStatements(query string) []string {
	segments := make([]string, 0, 2)
	var builder strings.Builder

	state := scanStateNormal
	for index := 0; index < len(query); index++ {
		next := byte(0)
		if index+1 < len(query) {
			next = query[index+1]
		}

		switch state {
		case scanStateNormal:
			switch {
			case query[index] == '-' && next == '-':
				builder.WriteByte(query[index])
				builder.WriteByte(next)
				index++
				state = scanStateLineComment
			case query[index] == '/' && next == '*':
				builder.WriteByte(query[index])
				builder.WriteByte(next)
				index++
				state = scanStateBlockComment
			case query[index] == '\'':
				builder.WriteByte(query[index])
				state = scanStateSingleQuote
			case query[index] == '"':
				builder.WriteByte(query[index])
				state = scanStateDoubleQuote
			case query[index] == '`':
				builder.WriteByte(query[index])
				state = scanStateBacktick
			case query[index] == '[':
				builder.WriteByte(query[index])
				state = scanStateBracket
			case query[index] == ';':
				appendStatement(&segments, builder.String())
				builder.Reset()
			default:
				builder.WriteByte(query[index])
			}
		case scanStateLineComment:
			builder.WriteByte(query[index])
			if query[index] == '\n' {
				state = scanStateNormal
			}
		case scanStateBlockComment:
			builder.WriteByte(query[index])
			if query[index] == '*' && next == '/' {
				builder.WriteByte(next)
				index++
				state = scanStateNormal
			}
		case scanStateSingleQuote:
			builder.WriteByte(query[index])
			if query[index] == '\'' {
				if next == '\'' {
					builder.WriteByte(next)
					index++
				} else {
					state = scanStateNormal
				}
			}
		case scanStateDoubleQuote:
			builder.WriteByte(query[index])
			if query[index] == '"' {
				if next == '"' {
					builder.WriteByte(next)
					index++
				} else {
					state = scanStateNormal
				}
			}
		case scanStateBacktick:
			builder.WriteByte(query[index])
			if query[index] == '`' {
				state = scanStateNormal
			}
		case scanStateBracket:
			builder.WriteByte(query[index])
			if query[index] == ']' {
				state = scanStateNormal
			}
		}
	}

	appendStatement(&segments, builder.String())
	return segments
}

func appendStatement(segments *[]string, raw string) {
	trimmed := strings.TrimSpace(trimLeadingTrivia(raw))
	if trimmed != "" {
		*segments = append(*segments, trimmed)
	}
}

func trimLeadingTrivia(value string) string {
	trimmed := value
	for {
		trimmed = strings.TrimLeft(trimmed, " \t\r\n")
		switch {
		case strings.HasPrefix(trimmed, "--"):
			newline := strings.IndexByte(trimmed, '\n')
			if newline < 0 {
				return ""
			}
			trimmed = trimmed[newline+1:]
		case strings.HasPrefix(trimmed, "/*"):
			end := strings.Index(trimmed, "*/")
			if end < 0 {
				return ""
			}
			trimmed = trimmed[end+2:]
		default:
			return trimmed
		}
	}
}

func tokenize(statement string) []sqlToken {
	tokens := make([]sqlToken, 0, 16)
	for index := 0; index < len(statement); {
		index = skipTrivia(statement, index)
		if index >= len(statement) {
			break
		}

		switch char := statement[index]; {
		case isIdentifierStart(char):
			start := index
			index++
			for index < len(statement) && isIdentifierPart(statement[index]) {
				index++
			}
			text := statement[start:index]
			tokens = append(tokens, sqlToken{text: text, upper: strings.ToUpper(text), start: start, kind: tokenIdentifier})
		case char == '\'':
			index = scanQuoted(statement, index, '\'', '\'')
		case char == '"':
			start := index
			index = scanQuoted(statement, index, '"', '"')
			text := statement[start:index]
			tokens = append(tokens, sqlToken{text: text, upper: strings.ToUpper(text), start: start, kind: tokenIdentifier})
		case char == '`':
			start := index
			index = scanUntil(statement, index+1, '`') + 1
			text := statement[start:index]
			tokens = append(tokens, sqlToken{text: text, upper: strings.ToUpper(text), start: start, kind: tokenIdentifier})
		case char == '[':
			start := index
			index = scanUntil(statement, index+1, ']') + 1
			text := statement[start:index]
			tokens = append(tokens, sqlToken{text: text, upper: strings.ToUpper(text), start: start, kind: tokenIdentifier})
		default:
			tokens = append(tokens, sqlToken{text: string(char), upper: strings.ToUpper(string(char)), start: index, kind: tokenSymbol})
			index++
		}
	}
	return tokens
}

func skipTrivia(value string, index int) int {
	for index < len(value) {
		switch {
		case value[index] == ' ' || value[index] == '\t' || value[index] == '\r' || value[index] == '\n':
			index++
		case index+1 < len(value) && value[index] == '-' && value[index+1] == '-':
			index += 2
			for index < len(value) && value[index] != '\n' {
				index++
			}
		case index+1 < len(value) && value[index] == '/' && value[index+1] == '*':
			index += 2
			for index+1 < len(value) && (value[index] != '*' || value[index+1] != '/') {
				index++
			}
			if index+1 < len(value) {
				index += 2
			} else {
				index = len(value)
			}
		default:
			return index
		}
	}
	return index
}

func scanQuoted(value string, index int, quote byte, escaped byte) int {
	index++
	for index < len(value) {
		if value[index] == quote {
			if index+1 < len(value) && value[index+1] == escaped {
				index += 2
				continue
			}
			return index + 1
		}
		index++
	}
	return len(value)
}

func scanUntil(value string, index int, terminator byte) int {
	for index < len(value) && value[index] != terminator {
		index++
	}
	if index >= len(value) {
		return len(value) - 1
	}
	return index
}

func isIdentifierStart(char byte) bool {
	return char == '_' || (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
}

func isIdentifierPart(char byte) bool {
	return isIdentifierStart(char) || (char >= '0' && char <= '9')
}

func closeSQLQueryConn(conn interface{ Close() error }) {
	_ = conn.Close()
}

func closeSQLQueryRows(rows interface{ Close() error }) {
	_ = rows.Close()
}

type scanState int

const (
	scanStateNormal scanState = iota
	scanStateLineComment
	scanStateBlockComment
	scanStateSingleQuote
	scanStateDoubleQuote
	scanStateBacktick
	scanStateBracket
)

type tokenKind int

const (
	tokenIdentifier tokenKind = iota
	tokenSymbol
)

type sqlToken struct {
	text  string
	upper string
	start int
	kind  tokenKind
}
