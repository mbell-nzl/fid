package db

import (
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

const driverName = "sqlite"

var requiredPragmas = []string{
	"foreign_keys(1)",
	"recursive_triggers(1)",
}

func Open(path string) (*sql.DB, error) {
	dsn, err := sqliteDSN(path)
	if err != nil {
		return nil, err
	}

	conn, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}

	if err := conn.Ping(); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("ping sqlite database: %w", err)
	}

	return conn, nil
}

func sqliteDSN(path string) (string, error) {
	absolutePath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve sqlite database path: %w", err)
	}

	dsnPath := filepath.ToSlash(absolutePath)
	if !strings.HasPrefix(dsnPath, "/") {
		dsnPath = "/" + dsnPath
	}

	query := url.Values{}
	for _, pragma := range requiredPragmas {
		query.Add("_pragma", pragma)
	}

	return (&url.URL{
		Scheme:   "file",
		Path:     dsnPath,
		RawQuery: query.Encode(),
	}).String(), nil
}
