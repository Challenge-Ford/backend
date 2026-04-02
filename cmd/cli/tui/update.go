package tui

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		case "esc":
			switch m.screen {
			case screenSeed:
				m.screen = screenMain
				m.cursor = 0
			case screenConfirm:
				m.screen = screenSeed
			case screenList:
				m.screen = screenMain
				m.cursor = 0
			case screenListView:
				m.screen = screenList
				m.cursor = 0
			}

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			m.cursor++
			switch m.screen {
			case screenMain:
				if m.cursor >= len(m.mainItems) {
					m.cursor = len(m.mainItems) - 1
				}
			case screenSeed:
				if m.cursor >= len(m.seedItems) {
					m.cursor = len(m.seedItems) - 1
				}
			case screenList:
				if m.cursor >= len(m.listItems) {
					m.cursor = len(m.listItems) - 1
				}
			}

		case "y", "Y":
			if m.screen == screenConfirm {
				ch := make(chan tea.Msg, 256)
				m.screen = screenRunning
				m.running = true
				m.progress = nil
				m.msgCh = ch
				return m, tea.Batch(
					m.spinner.Tick,
					startSeed(m.pendingSeed, ch, m.cfg),
					waitForMsg(ch),
				)
			}

		case "enter":
			switch m.screen {
			case screenMain:
				switch m.cursor {
				case 0:
					m.screen = screenSeed
					m.cursor = 0
				case 1:
					m.screen = screenList
					m.cursor = 0
				}

			case screenSeed:
				m.pendingSeed = m.cursor
				m.screen = screenConfirm

			case screenRunning:
				if !m.running {
					m.screen = screenSeed
					m.cursor = 0
					m.msgCh = nil
				}

			case screenList:
				m.pendingList = m.cursor
				m.screen = screenListLoading
				m.listResult = listResultMsg{}
				return m, tea.Batch(m.spinner.Tick, fetchList(m.cursor, m.cfg))

			case screenListView:
				m.screen = screenList
				m.cursor = 0
			}

		case "n", "N":
			if m.screen == screenConfirm {
				m.screen = screenSeed
			}
		}

	case spinner.TickMsg:
		if m.running || m.screen == screenListLoading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case progressMsg:
		m.progress = append(m.progress, progressEntry{
			section: msg.section,
			item:    msg.item,
			ok:      msg.ok,
		})
		return m, waitForMsg(m.msgCh)

	case runDoneMsg:
		m.running = false

	case listResultMsg:
		m.listResult = msg
		m.screen = screenListView
		m.cursor = 0
	}

	return m, nil
}
