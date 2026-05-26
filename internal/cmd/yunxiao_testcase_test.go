package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestYunxiaoCLITestcaseRepoListReturnsToolError(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"--disable-domains", "projex", "testcase", "repo-list"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected tool error")
	}
	if !strings.Contains(err.Error(), `unknown Yunxiao tool "list_testcase_repositories"`) {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLITestcaseViewRequiresRepoID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"testcase", "view", "tc-1"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repo-id error")
	}
	if !strings.Contains(err.Error(), "repo-id is required") {
		t.Fatalf("error = %v", err)
	}
}

func TestYunxiaoCLITestcaseSearchRequiresRepoID(t *testing.T) {
	command := NewYunxiaoCLI(IOStreams{Out: &bytes.Buffer{}, ErrOut: &bytes.Buffer{}})
	command.SetArgs([]string{"testcase", "search"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() expected repo-id error")
	}
	if !strings.Contains(err.Error(), "repo-id is required") {
		t.Fatalf("error = %v", err)
	}
}
