package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestYunxiaoCLITaskTypeListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "projex", "task", "type-list", "--project-id", "123", "--category", "Task"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_work_item_types"`) {
		t.Fatalf("error = %v", err)
	}
}
