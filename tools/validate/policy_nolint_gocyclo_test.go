package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindNolintGocycloViolations(t *testing.T) {
	tests := []struct {
		name      string
		setupDirs []string
		files     map[string]string
		wantCount int
		wantFirst string
	}{
		{
			name:      "empty tree",
			setupDirs: []string{"cmd", "internal", "tests"},
			files:     nil,
			wantCount: 0,
		},
		{
			name:      "detects directive",
			setupDirs: []string{"internal"},
			files: map[string]string{
				"internal/sample.go": "package x\n\nfunc f() {\n\t_ = 1\n}\n//nolint:gocyclo // reason\n",
			},
			wantCount: 1,
			wantFirst: "internal/sample.go:6",
		},
		{
			name:      "ignores occurrence inside string literal",
			setupDirs: []string{"internal"},
			files: map[string]string{
				"internal/policy.go": "package x\n\nfunc f() {\n\t_ = strings.Contains(s, \"//nolint:gocyclo\")\n}\n",
			},
			wantCount: 0,
		},
		{
			name:      "space after comment slashes",
			setupDirs: []string{"cmd"},
			files: map[string]string{
				"cmd/x.go": "// nolint:gocyclo\npackage main\n",
			},
			wantCount: 1,
			wantFirst: "cmd/x.go:1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			for _, sub := range tt.setupDirs {
				if err := os.MkdirAll(filepath.Join(dir, sub), 0o755); err != nil {
					t.Fatal(err)
				}
			}
			for rel, content := range tt.files {
				abs := filepath.Join(dir, rel)
				if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			got, err := findNolintGocycloViolations(dir)
			if err != nil {
				t.Fatalf("findNolintGocycloViolations: %v", err)
			}
			if len(got) != tt.wantCount {
				t.Fatalf("want %d violations, got %v", tt.wantCount, got)
			}
			if tt.wantFirst != "" && got[0] != tt.wantFirst {
				t.Fatalf("want first violation %q, got %q", tt.wantFirst, got[0])
			}
		})
	}
}
