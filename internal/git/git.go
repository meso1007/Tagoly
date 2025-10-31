package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Check if there are any staged changes
func HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return false
	}
	return out.Len() > 0
}

// Execute the actual `git commit` command
func Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Stdout = nil
	cmd.Stderr = nil
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to commit: %v", err)
	}
	return nil
}

// Get the list of staged files
func GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--cached")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %v", err)
	}
	files := strings.Split(strings.TrimSpace(string(out)), "\n")

	// Return empty slice if no files are staged
	if len(files) == 1 && files[0] == "" {
		return []string{}, nil
	}

	return files, nil
}
