package main

import (
	"os"

	"github.com/mbell-nzl/fid/internal/cli"
)

func main() {
	if err := execute(); err != nil {
		os.Exit(cli.ExitCode(err))
	}
}
