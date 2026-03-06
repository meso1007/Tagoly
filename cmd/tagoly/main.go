package main

import (
	"fmt"
	"log"
	"tagoly/internal/config"
	"tagoly/internal/generator"
	"tagoly/internal/git"
	"tagoly/internal/prompt"
)

func main() {

	if !git.HasStagedChanges() {
		fmt.Println("No staged files to commit. Please run 'git add' first.")
		return
	}
	// 1. setting load
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
		scope = prompt.SelectScope(scopeList, scope)
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
		if err := git.Commit(finalMessage); err != nil {
			log.Fatal("Failed to execute git commit: ", err)
		}
		fmt.Println("✅ Commit completed successfully!")
	} else {
		fmt.Println("❌ Commit canceled.")
	}
}
