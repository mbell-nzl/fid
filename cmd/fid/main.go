package main

import "os"

func main() {
	if err := execute(); err != nil {
		os.Exit(1)
	}
}
