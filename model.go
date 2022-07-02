package main

import (
	"fmt"
	"strings"
	"text/tabwriter"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
)

type Categories map[string]App

type App struct {
	prefix   string `yaml:",omitempty"`
	keybinds []KeyBind
}

type KeyBind struct {
	comment      string
	key          string
	ignorePrefix bool `yaml:"ignore_prefix,omitempty"`
}

type model struct {
	Styles
	Settings

	categories Categories // map of programs from keyb file
	headings   []string   // *ordered* slice of category names

	headingMap map[int]string // map of headings with line no.
	lineMap    map[int]string // map of all key bindings with line no.
	lineCount  int            // total number of lines

	// body is the final output split into lines. It is split to update the cursor
	// in updateBody()
	body []string

	height, width int
	maxWidth      int // for word wrapping
	padding       int // vertical padding - necessary to stabilize scrolling
	offset        int // for vertical scrolling
	cursor        int
}

type Styles struct {
	bodyStyle    lipgloss.Style
	headingStyle lipgloss.Style
	lineStyle    lipgloss.Style

	title       string
	curFg       string
	curBg       string
	border      string
	borderColor string
}

type Settings struct {
	mouseEnabled bool
	mouseDelta   int
}

func NewModel(cat Categories, config *Config) *model {

	st := Styles{
		bodyStyle:    lipgloss.NewStyle(),
		headingStyle: lipgloss.NewStyle().Bold(true).Margin(0, 1),
		lineStyle:    lipgloss.NewStyle().Margin(0, 2),
		title:        config.Title,
		curFg:        config.Cursor_fg,
		curBg:        config.Cursor_bg,
		border:       config.Border,
		borderColor:  config.Border_color,
	}

	s := Settings{
		mouseEnabled: true,
		mouseDelta:   3,
	}

	m := model{
		Styles:     st,
		Settings:   s,
		categories: cat,
		headings:   sortKeys(cat),

		height:   40,
		width:    60,
		maxWidth: 88,
		padding:  4,
	}
	if len(m.headings) > 0 {
		m.initBody()
	}
	return &m
}

func (m *model) initBody() {

	m.headingMap, m.lineMap = m.splitHeadingsAndKeys()
	keyTable := m.alignKeyLines()

	m.body = strings.Split(keyTable, "\n")
	m.insertHeadings()
}

// generate 2 maps from categories: untabbed headings and tabbed lines
func (m *model) splitHeadingsAndKeys() (map[int]string, map[int]string) {

	var headingIdx int
	var line string

	headingMap := make(map[int]string, len(m.headings))
	lineMap := make(map[int]string)

	for _, h := range m.headings {

		// headingIdx tracks the heading's line num and is offset
		// by the previous category's number of key lines
		headingMap[headingIdx] = m.headingStyle.Render(h)

		// this accounts for the wrapping of lines with width > max width
		var wrappedLineCount int

		app := m.categories[h]
		for i, key := range app.keybinds {

			// handle word wrapping for long lines
			if len(key.comment) >= m.maxWidth {
				wrappedLineCount = (len(key.comment) / m.maxWidth) + 1
				key.comment = wordwrap.String(key.comment, min(m.width, m.maxWidth))
			}

			if app.prefix != "" && !key.ignorePrefix {
				line = fmt.Sprintf("%s\t%s ; %s", key.comment, app.prefix, key.key)
			} else {
				line = fmt.Sprintf("%s\t%s", key.comment, key.key)
			}

			// each key's line num is offset's by its headingIdx
			lineMap[headingIdx+i+1] = m.lineStyle.Render(line)
		}

		// each category contributes to the total line count:
		// heading + num of keys + num of (extra) wrapped lines
		m.lineCount += len(app.keybinds) + 1 + wrappedLineCount

		// required offset
		headingIdx += len(app.keybinds) + max(1, wrappedLineCount)
	}
	return headingMap, lineMap
}

// generate aligned table
// TODO word wrapped lines do not align well
func (m *model) alignKeyLines() string {

	var sb strings.Builder
	tw := tabwriter.NewWriter(&sb, 20, 8, 10, ' ', 0)

	// order keys by line no.
	for i := 1; i < m.lineCount; i++ {
		if line, ok := m.lineMap[i]; ok {
			fmt.Fprintln(tw, line)
		}
	}
	tw.Flush()
	return sb.String()
}

// insert headings at respective line numbers
func (m *model) insertHeadings() {
	for i := 0; i < m.lineCount; i++ {
		if heading, ok := m.headingMap[i]; ok {
			m.body = insertAtIndex(i, heading, m.body)
		}
	}
}
