package main

import (
	"fmt"
	"log"
	"os/exec"

	"tagoly/internal/config"
	"tagoly/internal/generator"
	"tagoly/internal/git"
	"tagoly/internal/prompt"

	"github.com/AlecAivazis/survey/v2"
)

func main() {
	// 1. Load config
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// 2. Get staged files
	files, err := git.GetChangedFiles()
	if err != nil {
		log.Fatal("Failed to get changed files: ", err)
	}

	// 3. Detect scope
	scope, scopeList := generator.DetectScopeWithListImproved(files)
	if len(scopeList) > 1 {
		fmt.Println("Multiple scopes detected:", scopeList)
		promptScope := &survey.Select{
			Message:  "Select target scope:",
			Options:  scopeList,
			PageSize: 10,
			Default:  scope,
		}
		survey.AskOne(promptScope, &scope)
	}

	if scope == "root" {
		scope = ""
	}

	// 4. Convert custom tags (from config) into CommitType structs
	var customCommitTypes []prompt.CommitType
	for _, tag := range cfg.CustomTags {
		customCommitTypes = append(customCommitTypes, prompt.CommitType{
			Key:   tag,
			Label: "Custom tag",
		})
	}

	// 5. Select commit type
	tag := prompt.SelectCommitType(customCommitTypes)

	// 6. Enter commit message
	message := prompt.InputCommitMessage()

	// 7. Generate final commit message
	var finalMessage string
	if scope != "" {
		finalMessage = fmt.Sprintf("%s(%s): %s", tag, scope, message)
	} else {
		finalMessage = fmt.Sprintf("%s: %s", tag, message)
	}

	// 8. Confirm before committing
	fmt.Println("\nFinal commit message:")
	fmt.Println(finalMessage)

	if prompt.ConfirmCommit(finalMessage) {
		cmd := exec.Command("git", "commit", "-m", finalMessage)
		if err := cmd.Run(); err != nil {
			log.Fatal("Failed to execute git commit: ", err)
		}
		fmt.Println("✅ Commit completed successfully!")
	} else {
		fmt.Println("❌ Commit canceled.")
	}
}
