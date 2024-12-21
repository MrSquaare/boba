package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ButtonStyle struct {
	TextBase   lipgloss.Style
	TextFocus  lipgloss.Style
	TextActive lipgloss.Style
}

type ButtonState struct {
	focus  bool
	active bool
}

type Button struct {
	label string
	style ButtonStyle
	state ButtonState
}

var (
	ButtonDefaultStyle = ButtonStyle{
		TextBase: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Padding(0, 1),
		TextFocus: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Background(lipgloss.Color("33")).
			Padding(0, 1),
		TextActive: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("33")).
			Padding(0, 1),
	}
)

func NewButton(label string) *Button {
	m := &Button{
		label: label,
		style: ButtonDefaultStyle,
		state: ButtonState{
			focus:  false,
			active: false,
		},
	}

	return m
}

func (m *Button) Init() tea.Cmd {
	return nil
}

func (m *Button) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case FocusMsg:
		m.state.focus = typedMsg.Focus
	case ActiveMsg:
		m.state.active = typedMsg.Active
	}

	return m, nil
}

func (m *Button) View() string {
	if m.state.focus {
		return m.style.TextFocus.Render(m.label)
	}

	if m.state.active {
		return m.style.TextActive.Render(m.label)
	}

	return m.style.TextBase.Render(m.label)
}

func (m *Button) SetLabel(label string) *Button {
	m.label = label

	return m
}

func (m *Button) SetTextBaseStyle(style lipgloss.Style) *Button {
	m.style.TextBase = style

	return m
}

func (m *Button) SetTextFocusStyle(style lipgloss.Style) *Button {
	m.style.TextFocus = style

	return m
}

func (m *Button) SetTextActiveStyle(style lipgloss.Style) *Button {
	m.style.TextActive = style

	return m
}
