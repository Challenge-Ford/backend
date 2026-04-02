package tui

import (
	"strings"
)

func (m Model) View() string {
	switch m.screen {
	case screenMain:
		return m.viewMain()
	case screenSeed:
		return m.viewSeed()
	case screenConfirm:
		return m.viewConfirm()
	case screenRunning:
		return m.viewRunning()
	case screenList:
		return m.viewList()
	case screenListLoading:
		return m.viewListLoading()
	case screenListView:
		return m.viewListResult()
	}
	return ""
}

func (m Model) viewMain() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Torque CLI"))
	b.WriteString("\n\n")

	for i, item := range m.mainItems {
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("› " + item))
		} else {
			b.WriteString(normalStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("\n↑↓ navigate · enter select · q quit"))
	return b.String()
}

func (m Model) viewSeed() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("Seed"))
	b.WriteString("\n\n")

	for i, item := range m.seedItems {
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("› " + item))
		} else {
			b.WriteString(normalStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("\n↑↓ navigate · enter run · esc back · q quit"))
	return b.String()
}

func (m Model) viewConfirm() string {
	var b strings.Builder

	action := m.seedItems[m.pendingSeed]
	b.WriteString(titleStyle.Render("Confirm"))
	b.WriteString("\n\n")
	b.WriteString(normalStyle.Render("This will clear and reseed: "))
	b.WriteString(selectedStyle.Render(action))
	b.WriteString("\n\n")
	b.WriteString(errorStyle.Render("⚠  Existing data will be removed."))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("\ny confirm · n cancel · esc back"))
	return b.String()
}

func (m Model) viewList() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render("List"))
	b.WriteString("\n\n")

	for i, item := range m.listItems {
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("› " + item))
		} else {
			b.WriteString(normalStyle.Render("  " + item))
		}
		b.WriteString("\n")
	}

	b.WriteString(helpStyle.Render("\n↑↓ navigate · enter select · esc back · q quit"))
	return b.String()
}

func (m Model) viewListLoading() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render(m.spinner.View() + " Loading..."))
	b.WriteString(helpStyle.Render("\n\nq quit"))
	return b.String()
}

func (m Model) viewListResult() string {
	var b strings.Builder

	b.WriteString(titleStyle.Render(m.listResult.title))
	b.WriteString("\n\n")

	if len(m.listResult.entries) == 0 {
		b.WriteString(dimStyle.Render("  no entries"))
		b.WriteString("\n")
	}

	for _, entry := range m.listResult.entries {
		b.WriteString(sectionStyle.Render("▸ " + entry.title))
		b.WriteString("\n")
		for _, d := range entry.details {
			b.WriteString(normalStyle.Render("    " + d))
			b.WriteString("\n")
		}
	}

	b.WriteString(helpStyle.Render("\nenter back · esc back · q quit"))
	return b.String()
}

func (m Model) viewRunning() string {
	var b strings.Builder

	if m.running {
		b.WriteString(titleStyle.Render(m.spinner.View() + " Running..."))
	} else {
		b.WriteString(titleStyle.Render("✓ Done"))
	}
	b.WriteString("\n\n")

	var lastSection string
	for _, entry := range m.progress {
		if entry.section != lastSection {
			if lastSection != "" {
				b.WriteString("\n")
			}
			b.WriteString(sectionStyle.Render("▸ " + entry.section))
			b.WriteString("\n")
			lastSection = entry.section
		}

		if entry.item != "" {
			if entry.ok {
				b.WriteString(successStyle.Render("  ✓ " + entry.item))
			} else {
				b.WriteString(errorStyle.Render("  ✗ " + entry.item))
			}
			b.WriteString("\n")
		}
	}

	if !m.running {
		b.WriteString(helpStyle.Render("\nenter to go back"))
	}

	return b.String()
}
