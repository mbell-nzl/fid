package main

import (
	"os"

	"github.com/mbell-nzl/fid/internal/checkrun"
	"github.com/mbell-nzl/fid/internal/cli"
)

func main() {
	if child, err := checkrun.LaunchChild(); child {
		if err != nil {
			os.Exit(125)
		}
		return
	}
	if err := execute(); err != nil {
		os.Exit(cli.ExitCode(err))
	}
}
