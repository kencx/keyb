package list

import (
	"sort"

	"github.com/kencx/keyb/ui/table"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"
)

type filterState int

const (
	unfiltered filterState = iota
	filtering
)

type Model struct {
	Title  string
	keys   KeyMap
	styles Styles

	debug    bool
	viewport viewport.Model
	table    *table.Model

	searchBar textinput.Model
	search    bool

	filterState   filterState
	filteredTable *table.Model

	cursor  int
	padding int // vertical padding - necessary to stabilize scrolling
	maxRows int // max number of rows regardless of filterState
	Settings
}

type Settings struct {
	MouseEnabled bool
}

func New(title string, t *table.Model) Model {
	m := Model{
		Title:     title,
		table:     t,
		keys:      DefaultKeyMap(),
		styles:    DefaultStyle(),
		viewport:  viewport.Model{YOffset: 0, MouseWheelDelta: 3},
		searchBar: textinput.New(),
		padding:   4,
		Settings: Settings{
			MouseEnabled: true,
		},
	}
	m.debug = true

	if t.Empty() {
		m.table = table.NewEmpty(1)
		m.table.AppendRow(table.NewRow("No key bindings found", "", ""))
		m.table.Styles = m.styles.Table
		m.filteredTable = table.NewEmpty(m.table.LineCount)
		return m
	}

	m.maxRows = m.table.LineCount
	m.table.Styles = m.styles.Table
	m.filteredTable = table.NewEmpty(m.table.LineCount)
	m.filteredTable.Styles = m.styles.Table
	return m
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Resize(width, height int) {
	m.viewport.Width = width
	m.viewport.Height = height - m.padding
}

// Resets list to unfiltered state
func (m *Model) Reset() {
	m.filteredTable.Reset()
	m.filterState = unfiltered
	m.searchBar.Reset()
	m.cursorToBeginning()
	m.visibleRows()
}

// Sets items to be shown. All items are shown unless filtered
func (m *Model) visibleRows() {
	if !m.filteredTable.Empty() {
		m.SyncContent(m.filteredTable)
		m.maxRows = m.filteredTable.LineCount

	} else {
		// TODO for some reason, len(m.table.Output) != m.table.LineCount here
		m.SyncContent(m.table)
		m.maxRows = m.table.LineCount
	}
}

// Sync content by updating cursor and visible rows
func (m *Model) SyncContent(table *table.Model) {
	if table.Empty() {
		m.viewport.SetContent("")
		return
	}

	for i, row := range table.Rows {
		if i == m.cursor {
			row.IsSelected = true
		} else {
			row.IsSelected = false
		}
	}
	table.Render()
	m.viewport.SetContent(table.String())
}

func filter(term string, target []string) fuzzy.Matches {
	matches := fuzzy.Find(term, target)
	sort.Stable(matches)
	return matches
}

// TODO check this
func (m *Model) String() string {
	return m.table.String()
}

func (m *Model) UnstyledString() string {
	return m.table.Align()
}

func (m *Model) searchMode() bool {
	return m.search && m.searchBar.Focused()
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
