package search

import (
	"sort"
	"strings"
	"tagoly/internal/git"
	"tagoly/internal/parser"
)

// SearchOptions defines search filter criteria
type SearchOptions struct {
	Type    string // Filter by commit type (e.g., "feat", "fix")
	Scope   string // Filter by scope (e.g., "auth", "api")
	Subject string // Filter by subject text (substring match)
	Limit   int    // Maximum number of results (0 = all)
}

// SearchResult represents a single search result
type SearchResult struct {
	Hash    string // Git commit hash
	Type    string // Commit type
	Scope   string // Commit scope
	Subject string // Commit subject
	Raw     string // Full commit message
}

// Search queries git history and returns matching commits
func Search(opts SearchOptions) ([]SearchResult, error) {
	// Fetch commit history from git
	records, err := git.GetCommitHistory(0) // 0 = fetch all commits
	if err != nil {
		return nil, err
	}

	var results []SearchResult

	for _, record := range records {
		// Parse commit message to Tagoly format
		parsed := parser.ParseCommitMessage(record.Hash, record.Message)
		if parsed == nil {
			// Skip commits not in Tagoly format
			continue
		}

		// Apply filters
		if !matchesFilters(parsed, opts) {
			continue
		}

		results = append(results, SearchResult{
			Hash:    parsed.Hash,
			Type:    parsed.Type,
			Scope:   parsed.Scope,
			Subject: parsed.Subject,
			Raw:     parsed.Raw,
		})

		// Stop if limit is reached
		if opts.Limit > 0 && len(results) >= opts.Limit {
			break
		}
	}

	return results, nil
}

// SearchByType searches commits by type (tag)
func SearchByType(commitType string, limit int) ([]SearchResult, error) {
	return Search(SearchOptions{
		Type:  commitType,
		Limit: limit,
	})
}

// SearchByScope searches commits by scope
func SearchByScope(scope string, limit int) ([]SearchResult, error) {
	return Search(SearchOptions{
		Scope: scope,
		Limit: limit,
	})
}

// SearchByTypeAndScope searches by both type and scope
func SearchByTypeAndScope(commitType, scope string, limit int) ([]SearchResult, error) {
	return Search(SearchOptions{
		Type:  commitType,
		Scope: scope,
		Limit: limit,
	})
}

// AvailableScopes returns scopes found in Tagoly-formatted commit history.
func AvailableScopes() ([]string, error) {
	records, err := git.GetCommitHistory(0)
	if err != nil {
		return nil, err
	}

	seen := map[string]bool{}
	for _, record := range records {
		parsed := parser.ParseCommitMessage(record.Hash, record.Message)
		if parsed == nil {
			continue
		}
		seen[parsed.Scope] = true
	}

	scopes := make([]string, 0, len(seen))
	for scope := range seen {
		scopes = append(scopes, scope)
	}
	sort.Strings(scopes)
	return scopes, nil
}

// matchesFilters checks if a parsed commit matches all search filters
func matchesFilters(parsed *parser.ParsedCommit, opts SearchOptions) bool {
	// Filter by type (if specified)
	if opts.Type != "" && !strings.EqualFold(parsed.Type, opts.Type) {
		return false
	}

	// Filter by scope (if specified)
	if opts.Scope != "" && !strings.EqualFold(parsed.Scope, opts.Scope) {
		return false
	}

	// Filter by subject (case-insensitive substring match, if specified)
	if opts.Subject != "" {
		if !strings.Contains(strings.ToLower(parsed.Subject), strings.ToLower(opts.Subject)) {
			return false
		}
	}

	return true
}
