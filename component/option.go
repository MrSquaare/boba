package component

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type OptionProps struct {
	Label string
}

type OptionStyle struct {
	TextBase     lipgloss.Style
	TextFocus    lipgloss.Style
	TextActive   lipgloss.Style
	CursorBase   lipgloss.Style
	CursorFocus  lipgloss.Style
	CursorActive lipgloss.Style
}

type OptionState struct {
	focus  bool
	active bool
}

type Option struct {
	label  string
	cursor string
	style  OptionStyle
	state  OptionState
}

var (
	OptionDefaultStyle = OptionStyle{
		TextBase: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			MarginRight(1),
		TextFocus: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Background(lipgloss.Color("33")).
			MarginRight(1),
		TextActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true).
			MarginRight(1),
		CursorBase: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			PaddingRight(1),
		CursorFocus: lipgloss.NewStyle().
			Foreground(lipgloss.Color("15")).
			Bold(true).
			Background(lipgloss.Color("33")).
			PaddingRight(1),
		CursorActive: lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Bold(true).
			PaddingRight(1),
	}
)

func NewOption(label string) *Option {
	m := &Option{
		label:  label,
		cursor: ">",
		style:  OptionDefaultStyle,
		state: OptionState{
			focus:  false,
			active: false,
		},
	}

	return m
}

func (m *Option) Init() tea.Cmd {
	return nil
}

func (m *Option) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typedMsg := msg.(type) {
	case FocusMsg:
		m.state.focus = typedMsg.Focus
	case ActiveMsg:
		m.state.active = typedMsg.Active
	}

	return m, nil
}

func (m *Option) View() string {
	if m.state.focus {
		return m.style.CursorFocus.Render(m.cursor) + m.style.TextFocus.Render(m.label)
	}

	if m.state.active {
		return m.style.CursorActive.Render(m.cursor) + m.style.TextActive.Render(m.label)
	}

	return m.style.CursorBase.Render(m.cursor) + m.style.TextBase.Render(m.label)
}

func (m *Option) SetLabel(label string) *Option {
	m.label = label

	return m
}

func (m *Option) SetCursor(cursor string) *Option {
	m.cursor = cursor

	return m
}

func (m *Option) SetTextBaseStyle(style lipgloss.Style) *Option {
	m.style.TextBase = style

	return m
}

func (m *Option) SetTextFocusStyle(style lipgloss.Style) *Option {
	m.style.TextFocus = style

	return m
}

func (m *Option) SetTextActiveStyle(style lipgloss.Style) *Option {
	m.style.TextActive = style

	return m
}

func (m *Option) SetCursorBaseStyle(style lipgloss.Style) *Option {
	m.style.CursorBase = style

	return m
}

func (m *Option) SetCursorFocusStyle(style lipgloss.Style) *Option {
	m.style.CursorFocus = style

	return m
}

func (m *Option) SetCursorActiveStyle(style lipgloss.Style) *Option {
	m.style.CursorActive = style

	return m
}
