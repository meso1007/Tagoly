package main

import (
	"fmt"
	"os"
	"strings"
	"tagoly/internal/config"
	"tagoly/internal/generator"
	"tagoly/internal/git"
	"tagoly/internal/prompt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Warm Minimal Palette ---
var (
	brandColor = lipgloss.Color("#A3E635") // 鮮やかでモダンなライムグリーン
	textColor  = lipgloss.Color("#0F172A") // ライムにパキッと映える濃いネイビー/ブラック
	textMuted  = lipgloss.Color("#6B7280") // 操作説明用のグレー

	// タイトルを「大きく」見せるためのスタイル設定
	// 余白を広めに取り、塊としての存在感を強調
	styleTitle = lipgloss.NewStyle().
			Background(brandColor).
			Foreground(textColor).
			Padding(0, 3).   // 左右の余白を増やしてワイドに
			MarginBottom(1). // 下に少し空間を空ける
			Bold(true)

	styleSelected = lipgloss.NewStyle().
			Foreground(brandColor).
			Bold(true)

	styleHelp = lipgloss.NewStyle().
			Foreground(textMuted).
			MarginTop(1) // ヘルプテキストを少し離す
)

type step int

const (
	stepType step = iota
	stepScope
	stepSubject
	stepConfirm
	stepDone
)

type model struct {
	step          step
	cursor        int
	typeOptions   []prompt.CommitType
	scopeOptions  []string
	selectedType  prompt.CommitType
	selectedScope string
	subject       string
	committed     bool
	canceled      bool
	err           error
}

func defaultCommitTypes() []prompt.CommitType {
	return []prompt.CommitType{
		{Key: "feat", Label: "New feature"},
		{Key: "fix", Label: "Bug fix"},
		{Key: "docs", Label: "Documentation"},
		{Key: "refactor", Label: "Code improvement"},
		{Key: "style", Label: "UI / Style fixes"},
		{Key: "chore", Label: "Maintenance"},
	}
}

func newModel(customTags []prompt.CommitType, scopes []string, defaultScope string) model {
	types := append(defaultCommitTypes(), customTags...)
	if len(scopes) == 0 {
		scopes = []string{"root"}
	}

	scopeCursor := 0
	for i, scope := range scopes {
		if scope == defaultScope {
			scopeCursor = i
			break
		}
	}

	return model{
		step:         stepType,
		typeOptions:  types,
		scopeOptions: scopes,
		cursor:       0,
		selectedScope: func() string {
			if len(scopes) == 0 {
				return "root"
			}
			return scopes[scopeCursor]
		}(),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.canceled = true
			return m, tea.Quit
		}

		switch m.step {
		case stepType:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.typeOptions)-1 {
					m.cursor++
				}
			case "enter":
				m.selectedType = m.typeOptions[m.cursor]
				m.step = stepScope
				for i, scope := range m.scopeOptions {
					if scope == m.selectedScope {
						m.cursor = i
						return m, nil
					}
				}
				m.cursor = 0
			}
		case stepScope:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.scopeOptions)-1 {
					m.cursor++
				}
			case "enter":
				m.selectedScope = m.scopeOptions[m.cursor]
				m.step = stepSubject
				m.cursor = 0
			}
		case stepSubject:
			switch msg.String() {
			case "enter":
				m.subject = strings.TrimSpace(m.subject)
				if m.subject != "" {
					m.step = stepConfirm
					m.cursor = 0
				}
			case "backspace":
				if len(m.subject) > 0 {
					m.subject = m.subject[:len(m.subject)-1]
				}
			default:
				if len(msg.Runes) > 0 {
					m.subject += string(msg.Runes)
				}
			}
		case stepConfirm:
			switch msg.String() {
			case "left", "h", "up", "k":
				m.cursor = 0 // yes
			case "right", "l", "down", "j":
				m.cursor = 1 // no
			case "enter":
				if m.cursor == 1 {
					m.canceled = true
					return m, tea.Quit
				}

				finalMessage := fmt.Sprintf("%s(%s): %s", m.selectedType.Key, m.selectedScope, m.subject)
				if err := git.Commit(finalMessage); err != nil {
					m.err = err
				} else {
					m.committed = true
				}
				m.step = stepDone
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) finalMessage() string {
	return fmt.Sprintf("%s(%s): %s", m.selectedType.Key, m.selectedScope, strings.TrimSpace(m.subject))
}

func (m model) View() string {
	var b strings.Builder

	// 存在感のあるタイトル
	b.WriteString(styleTitle.Render("TAGOLY") + "\n\n")

	switch m.step {
	case stepType:
		b.WriteString("Select commit type:\n\n")
		for i, opt := range m.typeOptions {
			line := fmt.Sprintf("  %s (%s)", opt.Key, opt.Label)
			if m.cursor == i {
				line = "❯ " + strings.TrimPrefix(line, "  ")
				b.WriteString(styleSelected.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
		b.WriteString(styleHelp.Render("↑/↓ or j/k: move  Enter: select  q: quit"))
	case stepScope:
		b.WriteString(fmt.Sprintf("Type: %s\n\n", styleSelected.Render(m.selectedType.Key)))
		b.WriteString("Select scope:\n\n")
		for i, scope := range m.scopeOptions {
			line := "  " + scope
			if m.cursor == i {
				line = "❯ " + scope
				b.WriteString(styleSelected.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
		b.WriteString(styleHelp.Render("↑/↓ or j/k: move  Enter: select  q: quit"))
	case stepSubject:
		b.WriteString(fmt.Sprintf("Type: %s  Scope: %s\n\n", styleSelected.Render(m.selectedType.Key), styleSelected.Render(m.selectedScope)))
		b.WriteString("Enter commit subject:\n")
		// 入力中のテキストの後に改行を「2つ」入れて、ヘルプテキストを確実に押し下げる
		b.WriteString(styleSelected.Render("❯ ") + m.subject + "\n\n")
		b.WriteString(styleHelp.Render("Type text  Backspace: delete  Enter: continue  q: quit"))
	case stepConfirm:
		b.WriteString("Confirm commit message:\n\n")
		b.WriteString("  " + styleSelected.Render(m.finalMessage()) + "\n\n")
		yes := "  Yes "
		no := "  No "
		if m.cursor == 0 {
			yes = styleSelected.Render("❯ Yes ")
		} else {
			no = styleSelected.Render("❯ No ")
		}
		b.WriteString(fmt.Sprintf("%s   %s\n", yes, no))
		b.WriteString(styleHelp.Render("←/→ or h/l: choose  Enter: confirm  q: quit"))
	}

	return b.String()
}

func main() {
	if !git.HasStagedChanges() {
		fmt.Println("No staged files to commit. Please run 'git add' first.")
		return
	}

	cfg, err := config.LoadConfig()
	if err != nil {
		// Suppress error text to keep UI clean, fallback to default
	}

	changedFiles, err := git.GetChangedFiles()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	defaultScope, scopes := generator.DetectScopeWithListImproved(changedFiles)
	if len(scopes) == 0 {
		scopes = []string{defaultScope}
	}

	customTags := []prompt.CommitType{}
	if cfg != nil {
		customTags = cfg.CustomTags
	}

	p := tea.NewProgram(newModel(customTags, scopes, defaultScope))
	finalModel, err := p.Run()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	m, ok := finalModel.(model)
	if !ok {
		fmt.Println("Error: failed to finalize UI state.")
		os.Exit(1)
	}

	if m.canceled {
		fmt.Println("\n" + styleHelp.Render("Commit canceled."))
		return
	}
	if m.err != nil {
		fmt.Printf("\nError: %v\n", m.err)
		return
	}
	if m.committed {
		fmt.Printf("\n✅ %s\n", styleSelected.Render("Committed: "+m.finalMessage()))
	}
}
