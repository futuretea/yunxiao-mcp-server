package cmd

import (
	"bytes"
	"context"
	"net"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestVersionCommandPrintsVersionInfo(t *testing.T) {
	var out, errOut bytes.Buffer
	command := NewMCPServer(IOStreams{Out: &out, ErrOut: &errOut})
	command.SetArgs([]string{"version"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if !strings.Contains(out.String(), "yunxiao-mcp-server version=") {
		t.Fatalf("stdout = %q", out.String())
	}
	if errOut.Len() != 0 {
		t.Fatalf("stderr = %q", errOut.String())
	}
}

func TestRootCommandHelp(t *testing.T) {
	var out bytes.Buffer
	command := NewMCPServer(IOStreams{Out: &out, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--help"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	for _, want := range []string{"Yunxiao MCP Server", "--port", "--access-token"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("stdout = %q, missing %q", out.String(), want)
		}
	}
}

func TestRootCommandReportsMissingConfigFile(t *testing.T) {
	var errOut bytes.Buffer
	command := NewMCPServer(IOStreams{Out: &bytes.Buffer{}, ErrOut: &errOut})
	command.SetArgs([]string{"--config", "/path/does/not/exist.yaml"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected missing config error")
	}
	if !strings.Contains(err.Error(), "load config") {
		t.Fatalf("error = %v", err)
	}
}

func TestRootCommandValidatesEnabledToolsBeforeServing(t *testing.T) {
	restoreLogger := preserveLogger()
	t.Cleanup(restoreLogger)

	command := NewMCPServer(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--port", "1",
		"--enabled-tools", "not_a_tool",
		"--access-token", "token",
	})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected unknown tool error")
	}
	if !strings.Contains(err.Error(), `unknown MCP tool "not_a_tool"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestRootCommandStartsAndShutsDownHTTPServer(t *testing.T) {
	restoreLogger := preserveLogger()
	t.Cleanup(restoreLogger)

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	command := NewMCPServer(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{
		"--port", strconv.Itoa(port),
		"--access-token", "token",
	})
	command.SetContext(ctx)

	err = command.Execute()
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
}

func preserveLogger() func() {
	globalLevel := zerolog.GlobalLevel()
	logger := log.Logger
	timeFieldFormat := zerolog.TimeFieldFormat
	return func() {
		zerolog.SetGlobalLevel(globalLevel)
		log.Logger = logger
		zerolog.TimeFieldFormat = timeFieldFormat
	}
}
