package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/futuretea/yunxiao-mcp-server/internal/cmd"
	"github.com/futuretea/yunxiao-mcp-server/pkg/core/logging"
)

func init() {
	_ = logging.Initialize("error", os.Stderr)
}

func main() {
	command := cmd.NewMCPServer(cmd.IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	})

	if err := command.Execute(); err != nil {
		log.Fatal().Err(err).Msg("failed to execute command")
	}
}
