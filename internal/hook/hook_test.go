package hook

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWriteCommitMsgHookCreatesExecutableHook(t *testing.T) {
	hookPath := filepath.Join(t.TempDir(), ".git", "hooks", "commit-msg")

	if err := WriteCommitMsgHook(hookPath, InstallOptions{}); err != nil {
		t.Fatalf("WriteCommitMsgHook() error = %v", err)
	}

	got, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != CommitMsgHook {
		t.Fatalf("hook content mismatch:\n%s", got)
	}

	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Mode().Perm() != 0755 {
		t.Fatalf("hook mode = %v, want 0755", info.Mode().Perm())
	}
}

func TestWriteCommitMsgHookDoesNotOverwriteDifferentHook(t *testing.T) {
	hookPath := filepath.Join(t.TempDir(), ".git", "hooks", "commit-msg")
	if err := os.MkdirAll(filepath.Dir(hookPath), 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\necho custom\n"), 0755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := WriteCommitMsgHook(hookPath, InstallOptions{}); err == nil {
		t.Fatal("WriteCommitMsgHook() error = nil, want error")
	}

	got, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != "#!/bin/sh\necho custom\n" {
		t.Fatalf("hook was overwritten:\n%s", got)
	}
}

func TestWriteCommitMsgHookForceOverwritesDifferentHook(t *testing.T) {
	hookPath := filepath.Join(t.TempDir(), ".git", "hooks", "commit-msg")
	if err := os.MkdirAll(filepath.Dir(hookPath), 0755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(hookPath, []byte("#!/bin/sh\necho custom\n"), 0755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	if err := WriteCommitMsgHook(hookPath, InstallOptions{Force: true}); err != nil {
		t.Fatalf("WriteCommitMsgHook() error = %v", err)
	}

	got, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(got) != CommitMsgHook {
		t.Fatalf("hook content mismatch:\n%s", got)
	}
}
