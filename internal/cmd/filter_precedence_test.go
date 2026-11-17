package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestCLIBlankFiltersDoNotReexposeYAML(t *testing.T) {
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	content := []byte(`
read_only: true
compact: true
disabled_tools:
  - get_current_user
enabled_domains:
  - platform
`)
	if err := os.WriteFile(configPath, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	tests := []struct {
		name       string
		args       []string
		env        map[string]string
		wantUser   bool
		wantDomain string
	}{
		{
			name:       "file",
			wantUser:   false,
			wantDomain: "platform",
		},
		{
			name:       "empty env preserves file",
			env:        map[string]string{"YUNXIAO_MCP_DISABLED_TOOLS": "", "YUNXIAO_MCP_ENABLED_DOMAINS": ""},
			wantUser:   false,
			wantDomain: "platform",
		},
		{
			name:       "blank env preserves file",
			env:        map[string]string{"YUNXIAO_MCP_DISABLED_TOOLS": " ", "YUNXIAO_MCP_ENABLED_DOMAINS": ","},
			wantUser:   false,
			wantDomain: "platform",
		},
		{
			name:       "explicit empty flags clear file",
			args:       []string{"--disabled-tools=", "--enable-domains="},
			wantUser:   true,
			wantDomain: "",
		},
		{
			name:       "blank flags preserve file",
			args:       []string{"--disabled-tools= ", "--enable-domains=,"},
			wantUser:   false,
			wantDomain: "platform",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("YUNXIAO_MCP_ENABLED_TOOLS", "")
			t.Setenv("YUNXIAO_MCP_DISABLED_TOOLS", "")
			t.Setenv("YUNXIAO_MCP_ENABLED_DOMAINS", "")
			t.Setenv("YUNXIAO_MCP_DISABLED_DOMAINS", "")
			for key, value := range tt.env {
				t.Setenv(key, value)
			}

			names, domains := listCLITools(t, configPath, tt.args)
			_, hasUser := names["get_current_user"]
			if hasUser != tt.wantUser {
				t.Fatalf("get_current_user present = %v, want %v", hasUser, tt.wantUser)
			}
			if tt.wantDomain == "" {
				if _, onlyPlatform := domains["platform"]; onlyPlatform && len(domains) == 1 {
					t.Fatal("explicit empty domain flag should clear the platform allowlist")
				}
				return
			}
			if len(domains) != 1 || !domains[tt.wantDomain] {
				t.Fatalf("domains = %#v, want only %s", domains, tt.wantDomain)
			}
		})
	}
}

func listCLITools(t *testing.T, configPath string, extra []string) (map[string]bool, map[string]bool) {
	t.Helper()
	var out, errOut bytes.Buffer
	command := NewYunxiaoCLI(IOStreams{Out: &out, ErrOut: &errOut})
	args := append([]string{"--config", configPath}, extra...)
	args = append(args, "tools", "list", "--json")
	command.SetArgs(args)
	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v\nstderr=%s", err, errOut.String())
	}

	var summaries []toolSummary
	if err := json.Unmarshal(out.Bytes(), &summaries); err != nil {
		t.Fatalf("unmarshal stdout: %v\n%s", err, out.String())
	}
	names := make(map[string]bool, len(summaries))
	domains := make(map[string]bool)
	for _, summary := range summaries {
		names[summary.Name] = true
		domains[summary.Domain] = true
	}
	return names, domains
}
