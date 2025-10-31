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
	// 1. 設定読み込み
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration: ", err)
	}

	// 2. ステージ済みファイル取得
	files, err := git.GetChangedFiles()
	if err != nil {
		log.Fatal("Failed to get changed files: ", err)
	}

	// 3. スコープ検出
	scope, scopeList := generator.DetectScopeWithListImproved(files)
	if len(scopeList) > 1 {
		fmt.Println("Multiple scopes detected:", scopeList)
		promptScope := &survey.Select{
			Message:  "Select the scope for the commit:",
			Options:  scopeList,
			PageSize: 10,
			Default:  scope,
		}
		survey.AskOne(promptScope, &scope)
	}

	if scope == "root" {
		scope = ""
	}

	// 4. タグ選択
	customCommitTypes := cfg.CustomTags
	tag := prompt.SelectCommitType(customCommitTypes)

	// 5. メッセージ入力
	message := prompt.InputCommitMessage()

	// 6. コミットメッセージ生成
	var finalMessage string
	if scope != "" {
		finalMessage = fmt.Sprintf("%s(%s): %s", tag, scope, message)
	} else {
		finalMessage = fmt.Sprintf("%s: %s", tag, message)
	}

	// 7. 確認して git commit
	fmt.Println("\nFinal commit message:")
	fmt.Println(finalMessage)

	if prompt.ConfirmCommit("Commit with this message?") {
		cmd := exec.Command("git", "commit", "-m", finalMessage)
		if err := cmd.Run(); err != nil {
			log.Fatal("Failed to execute git commit: ", err)
		}
		fmt.Println("✅ Commit completed successfully!")
	} else {
		fmt.Println("❌ Commit canceled.")
	}
}
