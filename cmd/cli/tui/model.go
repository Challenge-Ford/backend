package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenMain screen = iota
	screenSeed
	screenConfirm
	screenRunning
	screenList
	screenListLoading
	screenListView
)

type progressMsg struct {
	section string
	item    string
	ok      bool
}

type runDoneMsg struct{}

type progressEntry struct {
	section string
	item    string
	ok      bool
}

type Config struct {
	KratosAdminURL string
	KetoReadURL    string
	KetoWriteURL   string
}

type Model struct {
	cfg         Config
	screen      screen
	cursor      int
	spinner     spinner.Model
	running     bool
	progress    []progressEntry
	msgCh       <-chan tea.Msg
	mainItems   []string
	seedItems   []string
	listItems   []string
	pendingSeed int
	pendingList int
	listResult  listResultMsg
}

func New(cfg Config) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = selectedStyle

	return Model{
		cfg:       cfg,
		screen:    screenMain,
		spinner:   s,
		mainItems: []string{"Seed", "List"},
		seedItems: []string{"Permissions", "Users"},
		listItems: []string{"Users", "Dealerships", "Roles"},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func waitForMsg(ch <-chan tea.Msg) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
