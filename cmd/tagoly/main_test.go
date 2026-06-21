package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
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

func TestSubjectCursorEditingSupportsInsertionAndBackspace(t *testing.T) {
	subject, cursor := insertRunesAtCursor("", 0, []rune("fix class"))
	subject, cursor = insertRunesAtCursor(subject, 3, []rune(" CSS"))
	if subject != "fix CSS class" {
		t.Fatalf("subject = %q, want %q", subject, "fix CSS class")
	}
	if cursor != 7 {
		t.Fatalf("cursor = %d, want 7", cursor)
	}

	subject, cursor = deleteRuneBeforeCursor(subject, cursor)
	if subject != "fix CS class" {
		t.Fatalf("subject after delete = %q, want %q", subject, "fix CS class")
	}
	if cursor != 6 {
		t.Fatalf("cursor after delete = %d, want 6", cursor)
	}
}

func TestSubjectCursorEditingHandlesMultibyteRunes(t *testing.T) {
	subject, cursor := insertRunesAtCursor("修正する", 2, []rune("CSS"))
	if subject != "修正CSSする" {
		t.Fatalf("subject = %q, want %q", subject, "修正CSSする")
	}
	if cursor != 5 {
		t.Fatalf("cursor = %d, want 5", cursor)
	}

	subject, cursor = deleteRuneBeforeCursor(subject, cursor)
	if subject != "修正CSする" {
		t.Fatalf("subject after delete = %q, want %q", subject, "修正CSする")
	}
	if cursor != 4 {
		t.Fatalf("cursor after delete = %d, want 4", cursor)
	}
}

func TestSubjectCursorBlinkTogglesOnlyDuringSubjectStep(t *testing.T) {
	m := model{step: stepSubject, cursorVisible: true}
	updated, cmd := m.Update(cursorBlinkMsg{})
	got := updated.(model)

	if got.cursorVisible {
		t.Fatal("cursorVisible = true, want false after blink")
	}
	if cmd == nil {
		t.Fatal("blink command = nil, want next blink while editing subject")
	}

	m = model{step: stepConfirm, cursorVisible: true}
	updated, cmd = m.Update(cursorBlinkMsg{})
	got = updated.(model)

	if !got.cursorVisible {
		t.Fatal("cursorVisible changed outside subject step")
	}
	if cmd != nil {
		t.Fatal("blink command returned outside subject step")
	}
}

func TestSubjectInputMakesCursorVisible(t *testing.T) {
	m := model{step: stepSubject, cursorVisible: false}
	updated, _ := m.Update(testKeyMsg("x"))
	got := updated.(model)

	if got.subject != "x" {
		t.Fatalf("subject = %q, want x", got.subject)
	}
	if !got.cursorVisible {
		t.Fatal("cursorVisible = false, want true after typing")
	}
}

func TestSubjectClearShortcutsEraseEntireLine(t *testing.T) {
	m := model{step: stepSubject, subject: "fix CSS class", subjectCursor: 7}

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	got := updated.(model)
	if got.subject != "" || got.subjectCursor != 0 {
		t.Fatalf("ctrl+u did not clear subject: %#v", got)
	}

	m = model{step: stepSubject, subject: "fix CSS class", subjectCursor: 7}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyBackspace, Alt: true})
	got = updated.(model)
	if got.subject != "" || got.subjectCursor != 0 {
		t.Fatalf("alt+backspace did not clear subject: %#v", got)
	}

	m = model{step: stepSubject, subject: "fix CSS class", subjectCursor: 7}
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDelete})
	got = updated.(model)
	if got.subject != "" || got.subjectCursor != 0 {
		t.Fatalf("delete did not clear subject: %#v", got)
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

func testKeyMsg(value string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(value)}
}
