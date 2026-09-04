package main

import (
	"context"
	"fmt"
	"os"

	shortycli "github.com/vekio/shorty/internal/cli"
)

func main() {
	command, err := shortycli.New()
	if err == nil {
		err = command.Run(context.Background(), os.Args)
	}
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
