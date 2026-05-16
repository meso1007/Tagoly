package main

import (
	"fmt"
	"strings"
	"tagoly/internal/config"
	"tagoly/internal/search"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- Styling for tagdict ---
var (
	accentType     = lipgloss.Color("#EC4899")
	accentScope    = lipgloss.Color("#06B6D4")
	styleDictTitle = lipgloss.NewStyle().
			Background(brandColor).
			Foreground(textColor).
			Padding(0, 3).
			MarginBottom(1).
			Bold(true)

	styleDictSelected = lipgloss.NewStyle().
				Foreground(brandColor).
				Bold(true)

	styleDictHelp = lipgloss.NewStyle().
			Foreground(textMuted).
			MarginTop(1)

	styleDictType = lipgloss.NewStyle().
			Foreground(accentType).
			Bold(true)

	styleDictScope = lipgloss.NewStyle().
			Foreground(accentScope).
			Bold(true)

	styleDictResult = lipgloss.NewStyle().
			Foreground(brandColor)
)

type dictStep int

const (
	dictStepSelectMode dictStep = iota
	dictStepSelectType
	dictStepSelectScope
	dictStepSearching
	dictStepResults
	dictStepDone
)

type tagDictModel struct {
	step           dictStep
	cursor         int
	allTypes       []string
	allScopes      []string
	selectedType   string
	selectedScope  string
	searchMode     string // "type", "scope", or "both"
	results        []search.SearchResult
	canceled       bool
	err            error
	resultsCursor  int
	lastSearchOpts search.SearchOptions
}

func newTagDictModel() *tagDictModel {
	// Load configuration to get all available types
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

	return &tagDictModel{
		step:     dictStepSelectMode,
		cursor:   0,
		allTypes: types,
		allScopes: []string{
			"root",
			"auth",
			"api",
			"ui",
			"db",
			"config",
			"test",
			"docs",
			"ci",
		},
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
			case "enter", "esc", "q":
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
	m.lastSearchOpts = opts

	results, err := search.Search(opts)
	if err != nil {
		m.err = err
		m.results = []search.SearchResult{}
	} else {
		m.results = results
	}
}

func (m *tagDictModel) View() string {
	var b strings.Builder

	b.WriteString(styleDictTitle.Render("TAGOLY SEARCH") + "\n\n")

	switch m.step {
	case dictStepSelectMode:
		b.WriteString("What would you like to search by?\n\n")
		options := []string{"Commit Type", "Scope"}
		for i, opt := range options {
			line := "  " + opt
			if m.cursor == i {
				line = "❯ " + opt
				b.WriteString(styleDictSelected.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
		b.WriteString(styleDictHelp.Render("↑/↓ or j/k: move  Enter: select  q: quit"))

	case dictStepSelectType:
		b.WriteString("Select commit type:\n\n")
		for i, t := range m.allTypes {
			line := "  " + t
			if m.cursor == i {
				line = "❯ " + t
				b.WriteString(styleDictSelected.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
		b.WriteString(styleDictHelp.Render("↑/↓ or j/k: move  Enter: search  ESC: back  q: quit"))

	case dictStepSelectScope:
		b.WriteString("Select scope:\n\n")
		for i, s := range m.allScopes {
			line := "  " + s
			if m.cursor == i {
				line = "❯ " + s
				b.WriteString(styleDictSelected.Render(line) + "\n")
			} else {
				b.WriteString(line + "\n")
			}
		}
		b.WriteString(styleDictHelp.Render("↑/↓ or j/k: move  Enter: search  ESC: back  q: quit"))

	case dictStepResults:
		if len(m.results) == 0 {
			b.WriteString("No commits found.\n\n")
		} else {
			if m.searchMode == "type" {
				b.WriteString(fmt.Sprintf("Commits with type: %s\n\n", styleDictType.Render(m.selectedType)))
			} else {
				b.WriteString(fmt.Sprintf("Commits with scope: %s\n\n", styleDictScope.Render(m.selectedScope)))
			}

			start := m.resultsCursor
			end := start + 5
			if end > len(m.results) {
				end = len(m.results)
			}

			for i := start; i < end; i++ {
				result := m.results[i]
				marker := "  "
				if i == m.resultsCursor {
					marker = "❯ "
				}

				typeStr := styleDictType.Render(result.Type)
				scopeStr := styleDictScope.Render(result.Scope)
				hashStr := result.Hash[:7]

				b.WriteString(fmt.Sprintf("%s[%s] %s(%s): %s\n",
					marker, hashStr, typeStr, scopeStr, result.Subject))
			}

			b.WriteString("\n")
			b.WriteString(styleDictResult.Render(
				fmt.Sprintf("Showing %d-%d of %d results", start+1, end, len(m.results)),
			))
			b.WriteString("\n")
		}
		b.WriteString(styleDictHelp.Render("↑/↓ or j/k: scroll  q or ESC: quit"))
	}

	return b.String()
}

// processTagDictCommand handles the tagdict subcommand
func processTagDictCommand(args []string) error {
	p := tea.NewProgram(newTagDictModel())
	finalModel, err := p.Run()
	if err != nil {
		return fmt.Errorf("UI error: %v", err)
	}

	m, ok := finalModel.(*tagDictModel)
	if !ok {
		return fmt.Errorf("failed to get final model")
	}

	if m.canceled {
		fmt.Println("\n" + styleDictHelp.Render("Search canceled."))
		return nil
	}

	if m.err != nil {
		return fmt.Errorf("search error: %v", m.err)
	}

	return nil
}
