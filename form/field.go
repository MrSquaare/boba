package form

import (
	"github.com/MrSquaare/boba/component"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FieldStyle struct {
	TextBase lipgloss.Style
}

type Field struct {
	label string
	child component.Component
	style FieldStyle
}

var (
	fieldDefaultStyle = FieldStyle{
		TextBase: lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")),
	}
)

func NewField(label string, child component.Component) *Field {
	m := &Field{
		label: label,
		child: child,
		style: fieldDefaultStyle,
	}

	return m
}

func (m *Field) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Field) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.child, cmd = m.child.Update(msg)

	return m, cmd
}

func (m *Field) View() string {
	var s string

	s += m.style.TextBase.Render(m.label) + "\n"
	s += m.child.View()

	return s
}

func (m *Field) Child() component.Component {
	return m.child
}

func (m *Field) SetTextBaseStyle(style lipgloss.Style) *Field {
	m.style.TextBase = style

	return m
}
