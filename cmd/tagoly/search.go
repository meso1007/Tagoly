package main

import (
	"flag"
	"fmt"
	"tagoly/internal/search"

	"github.com/charmbracelet/lipgloss"
)

// --- Styling for search output ---
var (
	accentType  = lipgloss.Color("#EC4899")
	accentScope = lipgloss.Color("#06B6D4")

	styleHeader = lipgloss.NewStyle().
			Foreground(brandColor).
			Bold(true)

	styleType = lipgloss.NewStyle().
			Foreground(accentType).
			Bold(true)

	styleScope = lipgloss.NewStyle().
			Foreground(accentScope).
			Bold(true)

	styleHash = lipgloss.NewStyle().
			Foreground(textMuted)

	styleCount = lipgloss.NewStyle().
			Foreground(brandColor).
			Bold(true)
)

// Command represents a CLI command
type Command struct {
	Name        string
	Description string
	Flags       *flag.FlagSet
	Run         func(args []string) error
}

// handleSearch processes the search subcommand
func handleSearch(args []string) error {
	fs := flag.NewFlagSet("search", flag.ContinueOnError)
	searchType := fs.String("type", "", "Filter by commit type (e.g., feat, fix, docs)")
	searchScope := fs.String("scope", "", "Filter by scope (e.g., auth, api)")
	searchSubject := fs.String("subject", "", "Filter by subject text (substring match)")
	limit := fs.Int("limit", 0, "Maximum number of results (0 = all)")

	err := fs.Parse(args)
	if err != nil {
		return fmt.Errorf("invalid search arguments: %v", err)
	}

	// Validate: at least one filter should be specified
	if *searchType == "" && *searchScope == "" && *searchSubject == "" {
		fmt.Println("Error: Please specify at least one search filter")
		fmt.Println("\nUsage: tagoly search [options]")
		fmt.Println("\nOptions:")
		fs.PrintDefaults()
		return nil
	}

	// Perform search
	opts := search.SearchOptions{
		Type:    *searchType,
		Scope:   *searchScope,
		Subject: *searchSubject,
		Limit:   *limit,
	}

	results, err := search.Search(opts)
	if err != nil {
		return fmt.Errorf("search error: %v", err)
	}

	// Display results
	if len(results) == 0 {
		fmt.Println(styleHeader.Render("No commits found matching your search criteria"))
		return nil
	}

	// Print header
	fmt.Println(styleHeader.Render("Search Results"))
	fmt.Printf("Found %s:\n\n", styleCount.Render(fmt.Sprintf("%d commit(s)", len(results))))

	// Print each result
	for i, result := range results {
		fmt.Printf("[%d] %s\n", i+1, styleHash.Render(result.Hash[:7]))
		fmt.Printf("    Type:    %s\n", styleType.Render(result.Type))
		fmt.Printf("    Scope:   %s\n", styleScope.Render(result.Scope))
		fmt.Printf("    Subject: %s\n", result.Subject)
		fmt.Printf("    Message: %s\n", result.Raw)
		fmt.Println()
	}

	return nil
}

// searchCommand returns the search command configuration
func searchCommand() *Command {
	return &Command{
		Name:        "search",
		Description: "Search commits by type, scope, or subject",
		Run:         handleSearch,
	}
}

// processSearchCommand handles the search subcommand
func processSearchCommand(args []string) error {
	if len(args) < 1 {
		fmt.Println("Error: search command requires options")
		fmt.Println("\nUsage: tagoly search [options]")
		fmt.Println("\nOptions:")
		fs := flag.NewFlagSet("search", flag.ContinueOnError)
		fs.String("type", "", "Filter by commit type (e.g., feat, fix, docs)")
		fs.String("scope", "", "Filter by scope (e.g., auth, api)")
		fs.String("subject", "", "Filter by subject text (substring match)")
		fs.Int("limit", 0, "Maximum number of results (0 = all)")
		fs.PrintDefaults()
		return nil
	}

	return handleSearch(args)
}
