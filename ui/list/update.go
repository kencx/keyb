package list

import (
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/kencx/keyb/ui/table"
	"github.com/sahilm/fuzzy"
)

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {

	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:

		// to play nice with borders and truncation,
		// <2 results in border exceeding max width
		m.viewport.Width = msg.Width - max(2, (m.padding*2+m.margin*2))
		m.viewport.Height = msg.Height - m.scrollOffset

		m.table.MaxWidth = m.viewport.Width - m.padding*2
		m.filteredTable.MaxWidth = m.viewport.Width - m.padding*2

		if m.cursorPastViewBottom() {
			m.cursor = m.viewport.YOffset + m.viewport.Height - 1
		}

	case tea.MouseMsg:
		if !m.viewport.MouseWheelEnabled {
			break
		}
		switch msg.Type {
		case tea.MouseWheelUp:
			m.cursor -= m.viewport.MouseWheelDelta
			if m.cursorPastViewTop() {
				m.viewport.LineUp(m.viewport.MouseWheelDelta)
			}
		case tea.MouseWheelDown:
			m.cursor += m.viewport.MouseWheelDelta
			if m.cursorPastViewBottom() {
				m.viewport.LineDown(m.viewport.MouseWheelDelta)
			}
		}
	}

	switch {
	case m.searchMode():
		cmds = append(cmds, m.handleSearch(msg))
	default:
		cmds = append(cmds, m.handleNormal(msg))
	}

	// cursor loop around
	if m.cursorPastBeginning() {
		m.cursorToEnd()
		m.viewport.GotoBottom()
	} else if m.cursorPastEnd() {
		m.cursorToBeginning()
		m.viewport.GotoTop()
	}

	m.visibleRows()
	return m, tea.Batch(cmds...)
}

func (m *Model) handleNormal(msg tea.Msg) tea.Cmd {

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return tea.Quit

		case key.Matches(msg, m.keys.Search):
			return m.startSearch()

		case key.Matches(msg, m.keys.ClearSearch):
			m.searchBar.Reset()
			return m.startSearch()

		case key.Matches(msg, m.keys.Up):
			m.cursor--
			if m.cursorPastViewTop() {
				m.viewport.LineUp(1)
			}
		case key.Matches(msg, m.keys.Down):
			m.cursor++
			if m.cursorPastViewBottom() {
				m.viewport.LineDown(1)
			}

		case key.Matches(msg, m.keys.UpFocus):
			m.cursor--
			if m.cursorPastViewTop() {
				m.viewport.LineUp(1)
			}
		case key.Matches(msg, m.keys.DownFocus):
			m.cursor++
			if m.cursorPastViewBottom() {
				m.viewport.LineDown(1)
			}

		case key.Matches(msg, m.keys.HalfUp):
			m.cursor -= m.viewport.Height / 2
			if m.cursorPastViewTop() {
				m.viewport.HalfViewUp()
			}

			// don't loop around
			if m.cursorPastBeginning() {
				m.cursorToBeginning()
				m.viewport.GotoTop()
			}
		case key.Matches(msg, m.keys.HalfDown):
			m.cursor += m.viewport.Height / 2
			if m.cursorPastViewBottom() {
				m.viewport.HalfViewDown()
			}

			// don't loop around
			if m.cursorPastEnd() {
				m.cursorToEnd()
				m.viewport.GotoBottom()
			}

		case key.Matches(msg, m.keys.FullUp):
			m.cursor -= m.viewport.Height
			if m.cursorPastViewTop() {
				m.viewport.ViewUp()
			}

			// don't loop around
			if m.cursorPastBeginning() {
				m.cursorToBeginning()
				m.viewport.GotoTop()
			}

		case key.Matches(msg, m.keys.FullDown):
			m.cursor += m.viewport.Height
			if m.cursorPastViewBottom() {
				m.viewport.ViewDown()
			}

			// don't loop around
			if m.cursorPastEnd() {
				m.cursorToEnd()
				m.viewport.GotoBottom()
			}

		case key.Matches(msg, m.keys.GoToFirstLine):
			m.cursorToBeginning()
			m.viewport.GotoTop()
		case key.Matches(msg, m.keys.GoToLastLine):
			m.cursorToEnd()
			m.viewport.GotoBottom()

		case key.Matches(msg, m.keys.GoToTop):
			m.cursorToViewTop()

		case key.Matches(msg, m.keys.GoToMiddle):
			m.cursorToViewMiddle()

		case key.Matches(msg, m.keys.GoToBottom):
			m.cursorToViewBottom()

			// case key.Matches(msg, m.keys.CenterCursor):
			// 	middle := m.viewport.Height / 2
			// 	diff := m.cursor - middle
			// 	if diff > 0 {
			// 		m.viewport.LineDown(diff)
			// 	} else {
			// 		m.viewport.LineUp(diff)
			// 	}
		}
	}
	return nil
}

func (m *Model) handleSearch(msg tea.Msg) tea.Cmd {

	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case msg.String() == "ctrl+c":
			return tea.Quit

		case key.Matches(msg, m.keys.ClearSearch):
			m.searchBar.Reset()
			return m.startSearch()

			// scrolling in search mode
		case key.Matches(msg, m.keys.UpFocus):
			m.cursor--
			if m.cursorPastViewTop() {
				m.viewport.LineUp(1)
			}
			return nil
		case key.Matches(msg, m.keys.DownFocus):
			m.cursor++
			if m.cursorPastViewBottom() {
				m.viewport.LineDown(1)
			}
			return nil

		case key.Matches(msg, m.keys.HalfUp):
			m.cursor -= m.viewport.Height / 2
			if m.cursorPastViewTop() {
				m.viewport.HalfViewUp()
			}

			// don't loop around
			if m.cursorPastBeginning() {
				m.cursorToBeginning()
				m.viewport.GotoTop()
			}
		case key.Matches(msg, m.keys.HalfDown):
			m.cursor += m.viewport.Height / 2
			if m.cursorPastViewBottom() {
				m.viewport.HalfViewDown()
			}

			// don't loop around
			if m.cursorPastEnd() {
				m.cursorToEnd()
				m.viewport.GotoBottom()
			}

		case key.Matches(msg, m.keys.FullUp):
			m.cursor -= m.viewport.Height
			if m.cursorPastViewTop() {
				m.viewport.ViewUp()
			}

			// don't loop around
			if m.cursorPastBeginning() {
				m.cursorToBeginning()
				m.viewport.GotoTop()
			}

		case key.Matches(msg, m.keys.FullDown):
			m.cursor += m.viewport.Height
			if m.cursorPastViewBottom() {
				m.viewport.ViewDown()
			}

			// don't loop around
			if m.cursorPastEnd() {
				m.cursorToEnd()
				m.viewport.GotoBottom()
			}

		case key.Matches(msg, m.keys.Normal):
			m.search = false
			m.searchBar.Blur()

			if m.filteredTable.Empty() {
				m.filterState = unfiltered
			}
			return nil
		}
	}

	// filter with search input
	m.searchBar, cmd = m.searchBar.Update(msg)
	cmds = append(cmds, cmd)

	prefix := "h:"
	if strings.HasPrefix(m.searchBar.Value(), prefix) {
		matchHeadings(m, prefix)
	} else {
		matchRows(m)
	}

	// reset if search input is empty regardless of filterState
	if m.searchBar.Value() == "" {
		m.Reset()

		// remain in filtering state
		// until user explicitly returns to Normal mode
		m.filterState = filtering
	}
	return tea.Batch(cmds...)
}

func filter(term string, target []string) fuzzy.Matches {
	matches := fuzzy.Find(term, target)
	sort.Stable(matches)
	return matches
}

func matchHeadings(m *Model, prefix string) {
	var matches fuzzy.Matches
	value := strings.TrimSpace(strings.TrimPrefix(m.searchBar.Value(), prefix))
	if !m.table.Empty() {
		matches = filter(value, m.table.GetPlainHeadings())
	}

	// present new filtered rows
	m.filteredTable.Reset()
	if len(matches) == 0 {
		m.filteredTable.AppendRow(table.EmptyRow())

	} else {
		var hlMatches []*table.Row
		// get non-pointers as filtering is ephemeral
		headings := m.table.GetCopyOfHeadings()

		for _, match := range matches {
			heading := headings[match.Index]
			heading.IsFiltered = true
			heading.MatchedIndex = match.MatchedIndexes

			hlMatches = append(hlMatches, &heading)
			hlMatches = append(hlMatches, m.table.GetAllRowsofHeading(heading.Text)...)
		}
		m.filteredTable.AppendRows(hlMatches...)
	}
}

func matchRows(m *Model) {
	var matches fuzzy.Matches
	if !m.table.Empty() {
		matches = filter(m.searchBar.Value(), m.table.GetPlainRowsWithoutHeadings())
	}

	// present new filtered rows
	m.filteredTable.Reset()
	if len(matches) == 0 {
		m.filteredTable.AppendRow(table.EmptyRow())

	} else {
		var hlMatches []*table.Row
		// get non-pointers as filtering is ephemeral
		rows := m.table.GetCopyOfRowsWithoutHeadings()

		for _, match := range matches {
			row := rows[match.Index]
			row.IsFiltered = true
			row.MatchedIndex = match.MatchedIndexes
			hlMatches = append(hlMatches, &row)
		}
		m.filteredTable.AppendRows(hlMatches...)
	}
}
