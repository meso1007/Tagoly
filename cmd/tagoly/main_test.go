package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestProcessCommitCommandCreatesCommitWithFlags(t *testing.T) {
	t.Chdir(t.TempDir())
	runGit(t, "init", "-q")
	runGit(t, "config", "user.name", "Test User")
	runGit(t, "config", "user.email", "test@example.com")
	runGit(t, "config", "commit.gpgsign", "false")

	if err := os.WriteFile("file.txt", []byte("hello\n"), 0644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
	runGit(t, "add", "file.txt")

	if err := processCommitCommand([]string{
		"--type", "fix",
		"--scope", "frontend",
		"--subject", "fix CSS class mismatch",
	}); err != nil {
		t.Fatalf("processCommitCommand() error = %v", err)
	}

	got := runGit(t, "log", "-1", "--pretty=%s")
	want := "fix(frontend): fix CSS class mismatch"
	if strings.TrimSpace(got) != want {
		t.Fatalf("commit subject = %q, want %q", strings.TrimSpace(got), want)
	}
}

func TestProcessCommitCommandRequiresTypeForNonInteractiveCommit(t *testing.T) {
	err := processCommitCommand([]string{"--subject", "missing type"})
	if err == nil {
		t.Fatal("processCommitCommand() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "-type is required") {
		t.Fatalf("processCommitCommand() error = %q, want -type requirement", err)
	}
}

func runGit(t *testing.T, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s failed: %v\n%s", strings.Join(args, " "), err, out)
	}
	return string(out)
}
