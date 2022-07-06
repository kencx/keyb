package list

import "fmt"

func (m *Model) View() string {

	view := fmt.Sprintf("%s\n"+
		"%s\n"+
		"%s\n"+
		"\n"+
		" keys: %d",
		m.Title,
		m.searchBar.View(),
		m.viewport.View(),
		m.table.LineCount)

	if m.debug {
		view = fmt.Sprintf("%s\tLine: %d YOffset: %d Height: %d",
			view, m.cursor, m.viewport.YOffset, m.viewport.Height)
	}
	// m.setStyle()
	return m.table.BodyStyle.Render(view)
}
