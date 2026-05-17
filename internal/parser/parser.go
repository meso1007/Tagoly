package parser

import (
	"regexp"
	"strings"
)

// ParsedCommit represents a parsed commit message in Tagoly format
type ParsedCommit struct {
	Hash    string // Git commit hash
	Type    string // e.g., "feat", "fix", "docs"
	Scope   string // e.g., "auth", "api", "root"
	Subject string // e.g., "add login feature"
	Raw     string // Original commit message
}

// commitsRegex matches Tagoly format: type(scope): subject
// Examples: "feat(auth): add login", "fix(api): resolve timeout", "docs: update readme"
var commitsRegex = regexp.MustCompile(`^([a-z]+)(?:\(([^)]+)\))?:\s+(.+)$`)

// ParseCommitMessage parses a commit message in Tagoly format
// Returns ParsedCommit if format is valid, otherwise returns nil
func ParseCommitMessage(hash, message string) *ParsedCommit {
	subjectLine := firstLine(message)
	matches := commitsRegex.FindStringSubmatch(subjectLine)
	if matches == nil {
		return nil // Not in Tagoly format
	}

	// matches[0] = full match
	// matches[1] = type (e.g., "feat")
	// matches[2] = scope (e.g., "auth") or empty
	// matches[3] = subject (e.g., "add login feature")

	scope := matches[2]
	if scope == "" {
		scope = "root"
	}

	return &ParsedCommit{
		Hash:    hash,
		Type:    matches[1],
		Scope:   scope,
		Subject: matches[3],
		Raw:     subjectLine,
	}
}

// IsValidCommitFormat checks if a message is in valid Tagoly format
func IsValidCommitFormat(message string) bool {
	return commitsRegex.MatchString(firstLine(message))
}

func firstLine(message string) string {
	message = strings.TrimSpace(message)
	if message == "" {
		return ""
	}
	line, _, _ := strings.Cut(message, "\n")
	return strings.TrimSpace(line)
}
