// Package migrations supplies the SQL history shipped with this fid build.
package migrations

import (
	"embed"
	"io/fs"
)

//go:embed *.sql
var files embed.FS

// FS contains all up and down migrations at the filesystem root.
var FS fs.FS = files
