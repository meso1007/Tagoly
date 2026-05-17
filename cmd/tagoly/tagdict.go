package main

import (
	"fmt"
	"strings"
	"tagoly/internal/config"
	"tagoly/internal/search"

	tea "github.com/charmbracelet/bubbletea"
)

type dictStep int

const (
	dictStepSelectMode dictStep = iota
	dictStepSelectType
	dictStepSelectScope
	dictStepResults
)

type tagDictModel struct {
	step          dictStep
	cursor        int
	allTypes      []string
	allScopes     []string
	selectedType  string
	selectedScope string
	searchMode    string // "type" or "scope"
	results       []search.SearchResult
	resultsCursor int
	canceled      bool
}

func newTagDictModel() *tagDictModel {
	cfg, _ := config.LoadConfig()

	types := []string{}
	for _, t := range defaultCommitTypes() {
		types = append(types, t.Key)
	}
	if cfg != nil {
		for _, ct := range cfg.CustomTags {
			types = append(types, ct.Key)
		}
	}

	scopes := []string{
		"root", "auth", "api", "ui", "db",
		"config", "test", "docs", "ci", "perf",
	}

	return &tagDictModel{
		step:      dictStepSelectMode,
		cursor:    0,
		allTypes:  types,
		allScopes: scopes,
	}
}

func (m *tagDictModel) Init() tea.Cmd {
	return nil
}

func (m *tagDictModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.canceled = true
			return m, tea.Quit
		}

		switch m.step {
		case dictStepSelectMode:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < 1 {
					m.cursor++
				}
			case "enter":
				if m.cursor == 0 {
					m.searchMode = "type"
					m.step = dictStepSelectType
					m.cursor = 0
				} else {
					m.searchMode = "scope"
					m.step = dictStepSelectScope
					m.cursor = 0
				}
			}

		case dictStepSelectType:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.allTypes)-1 {
					m.cursor++
				}
			case "enter":
				m.selectedType = m.allTypes[m.cursor]
				m.performSearch()
				m.step = dictStepResults
				m.resultsCursor = 0
			case "esc":
				m.step = dictStepSelectMode
				m.cursor = 0
			}

		case dictStepSelectScope:
			switch msg.String() {
			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}
			case "down", "j":
				if m.cursor < len(m.allScopes)-1 {
					m.cursor++
				}
			case "enter":
				m.selectedScope = m.allScopes[m.cursor]
				m.performSearch()
				m.step = dictStepResults
				m.resultsCursor = 0
			case "esc":
				m.step = dictStepSelectMode
				m.cursor = 0
			}

		case dictStepResults:
			switch msg.String() {
			case "up", "k":
				if m.resultsCursor > 0 {
					m.resultsCursor--
				}
			case "down", "j":
				if m.resultsCursor < len(m.results)-1 {
					m.resultsCursor++
				}
			case "esc":
				m.step = dictStepSelectMode
				m.cursor = 0
			case "q":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m *tagDictModel) performSearch() {
	opts := search.SearchOptions{
		Type:  m.selectedType,
		Scope: m.selectedScope,
		Limit: 0,
	}

	results, _ := search.Search(opts)
	m.results = results
	if m.results == nil {
		m.results = []search.SearchResult{}
	}
}

func (m *tagDictModel) View() string {
	var b strings.Builder
	b.WriteString(styleTitle.Render("TAGOLY SEARCH") + "\n\n")

	switch m.step {
	case dictStepSelectMode:
		b.WriteString("Search by:\n\n")
		options := []string{"Commit Type", "Scope"}
		for i, opt := range options {
			prefix := "  "
			if m.cursor == i {
				prefix = "❯ "
				b.WriteString(styleSelected.Render(prefix+opt) + "\n")
			} else {
				b.WriteString(prefix + opt + "\n")
			}
		}
		b.WriteString(styleHelp.Render("\n↑/↓ or j/k: move  Enter: select  q: quit"))

	case dictStepSelectType:
		b.WriteString("Select commit type:\n\n")
		for i, t := range m.allTypes {
			prefix := "  "
			if m.cursor == i {
				prefix = "❯ "
				b.WriteString(styleSelected.Render(prefix+t) + "\n")
			} else {
				b.WriteString(prefix + t + "\n")
			}
		}
		b.WriteString(styleHelp.Render("\n↑/↓ or j/k: move  Enter: search  ESC: back  q: quit"))

	case dictStepSelectScope:
		b.WriteString("Select scope:\n\n")
		for i, s := range m.allScopes {
			prefix := "  "
			if m.cursor == i {
				prefix = "❯ "
				b.WriteString(styleSelected.Render(prefix+s) + "\n")
			} else {
				b.WriteString(prefix + s + "\n")
			}
		}
		b.WriteString(styleHelp.Render("\n↑/↓ or j/k: move  Enter: search  ESC: back  q: quit"))

	case dictStepResults:
		if m.searchMode == "type" {
			b.WriteString(fmt.Sprintf("Type: %s\n\n", styleSelected.Render(m.selectedType)))
		} else {
			b.WriteString(fmt.Sprintf("Scope: %s\n\n", styleSelected.Render(m.selectedScope)))
		}

		if len(m.results) == 0 {
			b.WriteString("No commits found.\n")
		} else {
			start := m.resultsCursor
			if start > len(m.results)-1 {
				start = len(m.results) - 1
			}
			end := start + 5
			if end > len(m.results) {
				end = len(m.results)
			}

			for i := start; i < end; i++ {
				result := m.results[i]
				prefix := "  "
				if i == m.resultsCursor {
					prefix = "❯ "
				}

				hashStr := result.Hash[:7]
				commitLine := fmt.Sprintf("%s[%s] %s(%s): %s",
					prefix, hashStr, result.Type, result.Scope, result.Subject)

				if i == m.resultsCursor {
					b.WriteString(styleSelected.Render(commitLine) + "\n")
				} else {
					b.WriteString(commitLine + "\n")
				}
			}

			b.WriteString("\n")
			b.WriteString(styleHelp.Render(
				fmt.Sprintf("Showing %d-%d of %d  |  ↑/↓: scroll  ESC: back  q: quit",
					start+1, end, len(m.results)),
			))
		}
	}

	return b.String()
}

func processTagDictCommand(args []string) error {
	p := tea.NewProgram(newTagDictModel())
	_, err := p.Run()
	if err != nil {
		return fmt.Errorf("UI error: %v", err)
	}
	return nil
}
