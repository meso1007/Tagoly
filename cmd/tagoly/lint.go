package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"tagoly/internal/config"
	"tagoly/internal/git"
	tagolylint "tagoly/internal/lint"
)

func processLintCommand(args []string) error {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	message := fs.String("message", "", "Commit message to lint")
	messageFile := fs.String("message-file", "", "Path to a file containing a commit message")
	revRange := fs.String("range", "", "Git revision range to lint (e.g., main..HEAD)")
	limit := fs.Int("limit", 0, "Maximum number of commits to lint when using -range (0 = all)")

	if err := fs.Parse(args); err != nil {
		return fmt.Errorf("invalid lint arguments: %v", err)
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		return err
	}
	rules := tagolylint.DefaultRuleSet(cfg)

	switch {
	case *message != "":
		return printLintResults([]tagolylint.Result{
			tagolylint.ValidateMessage(*message, rules),
		})
	case *messageFile != "":
		data, err := os.ReadFile(*messageFile)
		if err != nil {
			return err
		}
		return printLintResults([]tagolylint.Result{
			tagolylint.ValidateMessage(string(data), rules),
		})
	case *revRange != "":
		records, err := git.GetCommitHistoryRange(*revRange, *limit)
		if err != nil {
			return err
		}
		results := make([]tagolylint.Result, 0, len(records))
		for _, record := range records {
			result := tagolylint.ValidateMessage(record.Message, rules)
			result.Message = fmt.Sprintf("%s %s", record.Hash[:7], firstMessageLine(record.Message))
			results = append(results, result)
		}
		return printLintResults(results)
	default:
		fmt.Println("Usage: tagoly lint -message <message>")
		fmt.Println("       tagoly lint -message-file <path>")
		fmt.Println("       tagoly lint -range <rev-range> [-limit <number>]")
		return nil
	}
}

func printLintResults(results []tagolylint.Result) error {
	failed := 0
	for _, result := range results {
		if result.Valid() {
			fmt.Printf("OK: %s\n", result.Message)
			continue
		}

		failed++
		fmt.Printf("FAIL: %s\n", result.Message)
		for _, lintErr := range result.Errors {
			fmt.Printf("  - %s\n", lintErr)
		}
	}

	if failed > 0 {
		return fmt.Errorf("%d commit message(s) failed lint", failed)
	}
	return nil
}

func firstMessageLine(message string) string {
	message = strings.TrimSpace(message)
	line, _, _ := strings.Cut(message, "\n")
	return line
}
