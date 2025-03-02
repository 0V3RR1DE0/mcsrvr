package main

import (
	"fmt"
	"os"

	"github.com/0v3rr1de0/mcsrvr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
