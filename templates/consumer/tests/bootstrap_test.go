package tests_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

var commitPinRE = regexp.MustCompile(`^[0-9a-f]{40}$`)

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "ai-sdlc.version")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("repository root not found (expected ai-sdlc.version)")
		}
		dir = parent
	}
}

// ensureValidator runs make build only when bin/validate is missing.
// Note: make build is a no-op when bin/validate already exists; delete bin/validate
// manually (or use TestValidatorBuild) to force a rebuild after ai-sdlc changes.
func ensureValidator(t *testing.T, root string) {
	t.Helper()
	bin := filepath.Join(root, "bin", "validate")
	if _, err := os.Stat(bin); err == nil {
		return
	}
	runMake(t, root, "build")
	if _, err := os.Stat(bin); err != nil {
		t.Fatalf("bin/validate missing after make build: %v", err)
	}
}

func runMake(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	cmd := exec.Command("make", args...)
	cmd.Dir = root
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("make %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return out
}

// Covers AC-00.001
func TestPinFilePresent(t *testing.T) {
	root := repoRoot(t)
	pinPath := filepath.Join(root, "ai-sdlc.version")
	data, err := os.ReadFile(pinPath)
	if err != nil {
		t.Fatalf("read ai-sdlc.version: %v", err)
	}
	pin := strings.TrimSpace(string(data))
	if pin == "" {
		t.Fatal("ai-sdlc.version is empty")
	}
	if strings.Contains(pin, " ") || strings.Contains(pin, "\t") {
		t.Fatalf("ai-sdlc.version contains whitespace: %q", pin)
	}

	if commitPinRE.MatchString(pin) {
		gitDir := filepath.Join(root, "ai-sdlc", ".git")
		if _, err := os.Stat(gitDir); err == nil {
			cmd := exec.Command("git", "-C", filepath.Join(root, "ai-sdlc"), "rev-parse", "HEAD")
			out, err := cmd.Output()
			if err != nil {
				t.Fatalf("git rev-parse HEAD in ai-sdlc: %v", err)
			}
			head := strings.TrimSpace(string(out))
			if !strings.EqualFold(head, pin) {
				t.Fatalf("ai-sdlc HEAD %q does not match ai-sdlc.version %q", head, pin)
			}
		}
		return
	}

	if len(pin) < 7 {
		t.Fatalf("ai-sdlc.version pin too short for tag: %q", pin)
	}
}

// Covers AC-00.002
func TestConsumerLayout(t *testing.T) {
	root := repoRoot(t)
	for _, rel := range []string{
		"AGENTS.md",
		".gitignore",
		"Makefile",
		"ai-sdlc.version",
		".golangci.yml",
		"scripts/check-module-boundaries.sh",
		".github/workflows/ci.yml",
	} {
		if _, err := os.Stat(filepath.Join(root, rel)); err != nil {
			t.Fatalf("missing root file %s: %v", rel, err)
		}
	}
	gitignore, err := os.ReadFile(filepath.Join(root, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	if !strings.Contains(string(gitignore), "ai-sdlc/") {
		t.Fatal(".gitignore must list ai-sdlc/")
	}
	if _, err := os.Stat(filepath.Join(root, ".github", "workflows", "ai-sdlc.yml")); err == nil {
		t.Fatal(".github/workflows/ai-sdlc.yml must not exist; use ci.yml for product gates")
	}

	base := filepath.Join(root, "ai-sdlc-artefacts")
	for _, rel := range []string{
		"scope.md",
		"strategy.md",
		"epics/EP-000/ep-scope.md",
		"epics/EP-000/ep-requirements.md",
		"epics/EP-000/ep-acceptance-criteria.md",
		"epics/EP-000/ep-context.md",
	} {
		if _, err := os.Stat(filepath.Join(base, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}
}

// Covers AC-00.003
func TestValidatorBuild(t *testing.T) {
	root := repoRoot(t)
	bin := filepath.Join(root, "bin", "validate")

	t.Run("build_from_missing", func(t *testing.T) {
		_ = os.Remove(bin)
		runMake(t, root, "build")
		if _, err := os.Stat(bin); err != nil {
			t.Fatalf("bin/validate missing after make build: %v", err)
		}
		cmd := exec.Command(bin, "structure", "EP-000")
		cmd.Dir = root
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("bin/validate structure EP-000: %v\n%s", err, out)
		}
	})
}

// Covers AC-00.005
func TestMakeValidateProject(t *testing.T) {
	root := repoRoot(t)
	ensureValidator(t, root)
	runMake(t, root, "validate")
}

// Covers AC-00.006
func TestMakeCheck(t *testing.T) {
	if os.Getenv("BOOTSTRAP_MAKE_CHECK_INVOKED") == "1" {
		// make check already ran vet and the test suite; avoid recursive make check.
		return
	}
	root := repoRoot(t)
	runMake(t, root, "check")
}

// Covers AC-00.004
func TestCIProductWorkflow(t *testing.T) {
	root := repoRoot(t)
	wfPath := filepath.Join(root, ".github", "workflows", "ci.yml")
	data, err := os.ReadFile(wfPath)
	if err != nil {
		t.Fatalf("read workflow: %v", err)
	}
	content := string(data)

	for _, needle := range []string{
		"ai-sdlc.version",
		"ai-sdlc",
		"make build",
		"make validate",
		"make check",
	} {
		if !strings.Contains(content, needle) {
			t.Fatalf("workflow missing %q", needle)
		}
	}

	for _, stepName := range []string{
		"Verify ai-sdlc pin",
		"Checkout ai-sdlc at pin",
		"Build validate binary",
		"Product check gate",
		"Product validate gate",
	} {
		if !strings.Contains(content, stepName) {
			t.Fatalf("workflow missing step name %q", stepName)
		}
	}
}
