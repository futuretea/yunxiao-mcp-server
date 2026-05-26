package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestYunxiaoCLIGroupListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "platform", "group", "list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_organization_groups"`) {
		t.Fatalf("error = %v", err)
	}
}
