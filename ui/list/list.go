package list

import (
	"github.com/kencx/keyb/config"
	"github.com/kencx/keyb/ui/table"

	"github.com/charmbracelet/bubbles/cursor"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type filterState int

const (
	unfiltered filterState = iota
	filtering
)

type Model struct {
	keys     KeyMap
	viewport viewport.Model
	table    *table.Model

	searchBar         textinput.Model
	search            bool
	startInSearchMode bool

	filterState    filterState
	filteredTable  *table.Model
	currentHeading string

	title   string
	debug   bool
	cursor  int
	maxRows int // max number of rows regardless of filterState

	margin         int
	padding        int
	scrollOffset   int
	border         lipgloss.Style
	counterStyle   lipgloss.Style
	promptLocation string
}

func New(t *table.Model, c *config.Config) Model {
	keyMap := CreateKeyMap(c.Keys)

	m := Model{
		keys: keyMap,
		viewport: viewport.Model{
			YOffset:           0,
			MouseWheelDelta:   3,
			MouseWheelEnabled: c.Mouse,
		},
		table: t,

		searchBar: textinput.Model{
			Prompt:           c.Prompt,
			PromptStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color(c.PromptColor)),
			Placeholder:      c.Placeholder,
			PlaceholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
			EchoCharacter:    '*',
			CharLimit:        0,
			Cursor:           cursor.New(),
			KeyMap:           textinput.KeyMap(keyMap.TextInputKeyMap),
		},
		startInSearchMode: c.SearchMode,

		filteredTable: table.NewEmpty(t.LineCount),

		title:   c.Title,
		debug:   c.Debug,
		cursor:  0,
		maxRows: t.LineCount,

		margin:         c.Margin,
		padding:        c.Padding,
		scrollOffset:   5,
		counterStyle:   lipgloss.NewStyle().Faint(true).Margin(0, 1),
		promptLocation: c.PromptLocation,
	}

	m.table.SepWidth = c.SepWidth
	m.filteredTable.SepWidth = c.SepWidth
	m.scrollOffset += (m.margin * 2) + (m.padding * 2)
	m.style(c)

	if m.startInSearchMode {
		m.startSearch()
	}

	return m
}

func (m *Model) style(c *config.Config) {
	if c.PlaceholderFg != "" || c.PlaceholderBg != "" {
		m.searchBar.PlaceholderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(c.PlaceholderFg)).
			Background(lipgloss.Color(c.PlaceholderBg))
	}

	if c.CounterFg != "" || c.CounterBg != "" {
		m.counterStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(c.CounterFg)).Background(lipgloss.Color(c.CounterBg)).Margin(0, 1)
	}

	var b lipgloss.Border
	switch c.BorderStyle {
	case "normal":
		b = lipgloss.NormalBorder()
	case "rounded":
		b = lipgloss.RoundedBorder()
	case "double":
		b = lipgloss.DoubleBorder()
	case "thick":
		b = lipgloss.ThickBorder()
	default:
		b = lipgloss.HiddenBorder()
	}
	m.border = lipgloss.NewStyle().BorderStyle(b).BorderForeground(lipgloss.Color(c.BorderColor))

	// row specific config
	if !m.table.Empty() {
		cursor := lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color(c.CursorFg)).
			Background(lipgloss.Color(c.CursorBg))

		s := table.RowStyles{
			Normal:          lipgloss.NewStyle().Margin(0, 2),
			Heading:         lipgloss.NewStyle().Margin(0, 1).Bold(true),
			Selected:        cursor.Copy().Margin(0, 2),
			SelectedHeading: cursor.Copy().Margin(0, 1).Bold(true),
			Filtered: lipgloss.NewStyle().
				Foreground(lipgloss.Color(c.FilterFg)).
				Background(lipgloss.Color(c.FilterBg)),
		}

		for _, row := range m.table.Rows {
			row.PrefixSep = c.PrefixSep
			row.Reversed = c.Reverse
			row.Styles = s
		}
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Resize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height - m.scrollOffset
}

// Resets list to unfiltered state
func (m *Model) Reset() {
	m.filteredTable.Reset()
	m.filterState = unfiltered
	m.currentHeading = ""
	m.cursorToBeginning()
	m.visibleRows()
}

// Sets items to be shown. All items are shown unless filtered
func (m *Model) visibleRows() {
	if !m.filteredTable.Empty() {
		m.SyncContent(m.filteredTable)
	} else {
		m.SyncContent(m.table)
	}
}

// Sync content by updating cursor and visible rows
func (m *Model) SyncContent(table *table.Model) {
	if table.Empty() {
		m.viewport.SetContent("")
		return
	}

	if m.cursor > table.LineCount {
		m.cursor = table.LineCount - 1
	}

	for i, row := range table.Rows {
		if i == m.cursor {
			row.IsSelected = true
			m.currentHeading = row.Heading
		} else {
			row.IsSelected = false
		}
	}
	m.viewport.SetContent(table.Render())
	m.maxRows = table.LineCount
}

func (m *Model) UnstyledString() string {
	return m.table.GetAlignedRows()
}

func (m *Model) searchMode() bool {
	return m.search && m.searchBar.Focused()
}

func (m *Model) startSearch() tea.Cmd {
	m.search = true
	m.filterState = filtering
	m.searchBar.Focus()
	return textinput.Blink
}

func (m *Model) cursorToBeginning() {
	m.cursor = 0
}

func (m *Model) cursorToEnd() {
	m.cursor = m.maxRows - 1
}

func (m *Model) cursorToViewTop() {
	m.cursor = m.viewport.YOffset + 3
}

func (m *Model) cursorToViewMiddle() {
	m.cursor = (m.viewport.YOffset + m.viewport.Height) / 2
}

func (m *Model) cursorToViewBottom() {
	m.cursor = m.viewport.YOffset + m.viewport.Height - 3
}

func (m *Model) cursorPastViewTop() bool {
	return m.cursor < m.viewport.YOffset
}

func (m *Model) cursorPastViewBottom() bool {
	return m.cursor > m.viewport.YOffset+m.viewport.Height-1
}

func (m *Model) cursorPastBeginning() bool {
	return m.cursor < 0
}

func (m *Model) cursorPastEnd() bool {
	return m.cursor > m.maxRows-1
}
