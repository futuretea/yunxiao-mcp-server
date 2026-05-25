package main

import (
	"fmt"
	"os"

	"github.com/futuretea/yunxiao-mcp-server/internal/cmd"
)

func main() {
	command := cmd.NewYunxiaoCLI(cmd.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	})

	if err := command.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
