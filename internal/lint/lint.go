package lint

import (
	"fmt"
	"strings"
	"tagoly/internal/config"
	"tagoly/internal/parser"
)

type RuleSet struct {
	AllowedTypes map[string]bool
}

type Result struct {
	Message string
	Errors  []string
}

func DefaultRuleSet(cfg *config.Config) RuleSet {
	allowed := map[string]bool{
		"feat":     true,
		"fix":      true,
		"docs":     true,
		"refactor": true,
		"style":    true,
		"chore":    true,
	}

	if cfg != nil {
		for _, tag := range cfg.CustomTags {
			allowed[tag.Key] = true
		}
	}

	return RuleSet{AllowedTypes: allowed}
}

func ValidateMessage(message string, rules RuleSet) Result {
	result := Result{Message: strings.TrimSpace(message)}
	parsed := parser.ParseCommitMessage("", message)
	if parsed == nil {
		result.Errors = append(result.Errors, "message must match type(scope): subject or type: subject")
		return result
	}

	if !rules.AllowedTypes[parsed.Type] {
		result.Errors = append(result.Errors, fmt.Sprintf("type %q is not allowed", parsed.Type))
	}
	if strings.TrimSpace(parsed.Subject) == "" {
		result.Errors = append(result.Errors, "subject must not be empty")
	}

	return result
}

func (r Result) Valid() bool {
	return len(r.Errors) == 0
}
