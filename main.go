package main

import (
	"fmt"
	"os"

	"batsig/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		if err == cli.ErrUsage {
			cli.PrintUsage(os.Stderr)
			os.Exit(2)
		}
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
