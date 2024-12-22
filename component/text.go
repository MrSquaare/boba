package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type TextStyle struct {
	TextBase lipgloss.Style
}

type Text struct {
	content string
	style   TextStyle
}

var (
	TextDefaultStyle = TextStyle{
		TextBase: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")),
	}
)

func NewText(label string) *Text {
	m := &Text{
		content: label,
		style:   TextDefaultStyle,
	}

	return m
}

func (m *Text) Init() tea.Cmd {
	return nil
}

func (m *Text) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *Text) View() string {
	return m.style.TextBase.Render(m.content)
}

func (m *Text) SetContent(content string) *Text {
	m.content = content

	return m
}

func (m *Text) SetTextBaseStyle(style lipgloss.Style) *Text {
	m.style.TextBase = style

	return m
}
