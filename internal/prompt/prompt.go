package prompt

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
)

// Commit Type Struct
type CommitType struct {
	Key   string
	Label string
}

// Standard Commit Types
var defaultCommitTypes = []CommitType{
	{"feat", "New feature"},
	{"fix", "Bug fix"},
	{"docs", "Documentation"},
	{"refactor", "Code improvement"},
	{"style", "UI / Style fixes"},
	{"chore", "Maintenance"},
}

// Select commit type, also supporting custom tags
func SelectCommitType(customTags []CommitType) string {
	options := []string{}

	// Standard tags
	for _, c := range defaultCommitTypes {
		options = append(options, fmt.Sprintf("%s (%s)", c.Key, c.Label))
	}

	// Custom tags
	for _, ct := range customTags {
		options = append(options, fmt.Sprintf("%s (%s, Custom)", ct.Key, ct.Label))
	}

	var selected string
	prompt := &survey.Select{
		Message:  "⚙️ Select commit type:",
		Options:  options,
		PageSize: 10,
	}
	survey.AskOne(prompt, &selected)

	// Return the key based on the selection result
	for _, c := range defaultCommitTypes {
		if selected == fmt.Sprintf("%s (%s)", c.Key, c.Label) {
			return c.Key
		}
	}
	for _, ct := range customTags {
		if selected == fmt.Sprintf("%s (%s, Custom)", ct.Key, ct.Label) {
			return ct.Key
		}
	}

	return "feat"
}

// Input commit message
func InputCommitMessage() string {
	var message string
	prompt := &survey.Input{
		Message: "📝 Enter commit message:",
	}
	survey.AskOne(prompt, &message)
	return message
}

// Confirmation dialog
func ConfirmCommit(final string) bool {
	var confirm bool
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("💬 Generated message:\n%s\n\n✅ Commit with this message?", final),
		Default: true,
	}
	survey.AskOne(prompt, &confirm)
	return confirm
}
