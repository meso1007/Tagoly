package git

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// CommitRecord represents a git commit with hash and message
type CommitRecord struct {
	Hash    string
	Message string
}

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

// GetCommitHistory retrieves git commit history with hash and message
// limit: maximum number of commits to retrieve (0 = all)
func GetCommitHistory(limit int) ([]CommitRecord, error) {
	args := []string{"log", "--pretty=format:%H%n%B%n---COMMIT_SEPARATOR---"}
	if limit > 0 {
		args = append([]string{"log", "-n", fmt.Sprintf("%d", limit)}, "--pretty=format:%H%n%B%n---COMMIT_SEPARATOR---")
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit history: %v", err)
	}

	var records []CommitRecord
	commitBlocks := strings.Split(string(out), "---COMMIT_SEPARATOR---")

	for _, block := range commitBlocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}

		lines := strings.SplitN(block, "\n", 2)
		if len(lines) < 2 {
			continue
		}

		hash := strings.TrimSpace(lines[0])
		message := strings.TrimSpace(lines[1])

		if hash != "" && message != "" {
			records = append(records, CommitRecord{
				Hash:    hash,
				Message: message,
			})
		}
	}

	return records, nil
}
