package db

import "database/sql"

// Queryer reads projection rows through a database or a caller's transaction.
type Queryer interface {
	Query(string, ...any) (*sql.Rows, error)
}
