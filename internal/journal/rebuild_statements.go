package journal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// Bound memory even when a journal uses many different valid column sets.
const rebuildStatementLimit = 128

// Rebuild owns a single connection. Statements execute on that connection
// inside the current Event transaction; they never escape into live writes or
// another workspace. Pinning also supports a DB pool limited to one connection.
type rebuildStatements struct {
	conn       *sql.Conn
	statements map[string]*sql.Stmt
}

func replayRebuild(db *sql.DB, journalDir string) (stats ReplayStats, err error) {
	conn, err := db.Conn(context.Background())
	if err != nil {
		return stats, fmt.Errorf("acquire rebuild connection: %w", err)
	}
	writes := &rebuildStatements{conn: conn, statements: make(map[string]*sql.Stmt)}
	defer func() { err = errors.Join(err, writes.close()) }()
	return replayJournal(journalDir, nil, writes.project)
}

func (writes *rebuildStatements) project(input ProjectionInput) error {
	if err := validateProjectionInput(input); err != nil {
		return err
	}
	tx, err := writes.conn.BeginTx(context.Background(), nil)
	if err != nil {
		return fmt.Errorf("project event %q: begin transaction: %w", input.Event.EventID, err)
	}
	return projectEventTransaction(tx, input, writes)
}

func (writes *rebuildStatements) Exec(query string, args ...any) (sql.Result, error) {
	statement := writes.statements[query]
	if statement == nil {
		if len(writes.statements) >= rebuildStatementLimit {
			return writes.conn.ExecContext(context.Background(), query, args...)
		}
		var err error
		statement, err = writes.conn.PrepareContext(context.Background(), query)
		if err != nil {
			return nil, err
		}
		writes.statements[query] = statement
	}
	return statement.Exec(args...)
}

func (writes *rebuildStatements) close() error {
	var result error
	for _, statement := range writes.statements {
		result = errors.Join(result, statement.Close())
	}
	result = errors.Join(result, writes.conn.Close())
	if result != nil {
		return fmt.Errorf("close rebuild statements and connection: %w", result)
	}
	return nil
}
