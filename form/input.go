package form

import (
	"github.com/MrSquaare/boba/component"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type InputStyle struct {
	TextBase         lipgloss.Style
	TextFocus        lipgloss.Style
	Error            lipgloss.Style
	PromptBase       lipgloss.Style
	PromptFocus      lipgloss.Style
	PlaceholderBase  lipgloss.Style
	PlaceholderFocus lipgloss.Style
	Cursor           lipgloss.Style
}

type InputState struct {
	focus bool
}

type Input struct {
	validateFunc func(string) error
	err          error
	input        textinput.Model
	style        InputStyle
	state        InputState
}

var (
	InputDefaultStyle = InputStyle{
		TextBase:         lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		TextFocus:        lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		PromptBase:       lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		PromptFocus:      lipgloss.NewStyle().Foreground(lipgloss.Color("15")),
		PlaceholderBase:  lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		PlaceholderFocus: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
		Cursor:           lipgloss.NewStyle().Foreground(lipgloss.Color("33")),
		Error:            lipgloss.NewStyle().Foreground(lipgloss.Color("9")),
	}
)

func NewInput() *Input {
	m := &Input{
		validateFunc: func(string) error { return nil },
		err:          nil,
		input:        textinput.New(),
		style:        InputDefaultStyle,
		state: InputState{
			focus: false,
		},
	}

	m.input.TextStyle = m.style.TextBase

	return m
}

func (m *Input) Init() tea.Cmd {
	return textinput.Blink
}

func (m *Input) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	switch typedMsg := msg.(type) {
	case component.FocusMsg:
		m.state.focus = typedMsg.Focus
		m.updateStyle()

		if m.state.focus {
			cmds = append(cmds, m.input.Focus())
		} else {
			m.input.Blur()
		}
	}

	var cmd tea.Cmd

	m.input, cmd = m.input.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Input) View() string {
	var s string

	s += m.input.View()

	if m.err != nil {
		s += "\n" + m.style.Error.Render(m.err.Error())
	}

	return s
}

func (m *Input) Value() string {
	return m.input.Value()
}

func (m *Input) Validate() bool {
	m.err = m.validateFunc(m.input.Value())

	return m.err == nil
}

func (m *Input) Error() error {
	return m.err
}

func (m *Input) SetPlaceholder(placeholder string) *Input {
	m.input.Placeholder = placeholder

	return m
}

func (m *Input) SetValidateFunc(fn func(string) error) *Input {
	m.validateFunc = fn

	return m
}

func (m *Input) SetValue(value string) *Input {
	m.input.SetValue(value)

	return m
}

func (m *Input) SetTextBaseStyle(style lipgloss.Style) *Input {
	m.style.TextBase = style

	return m
}

func (m *Input) SetTextFocusStyle(style lipgloss.Style) *Input {
	m.style.TextFocus = style

	return m
}

func (m *Input) SetPromptBaseStyle(style lipgloss.Style) *Input {
	m.style.PromptBase = style

	return m
}

func (m *Input) SetPromptFocusStyle(style lipgloss.Style) *Input {
	m.style.PromptFocus = style

	return m
}

func (m *Input) SetPlaceholderBaseStyle(style lipgloss.Style) *Input {
	m.style.PlaceholderBase = style

	return m
}

func (m *Input) SetPlaceholderFocusStyle(style lipgloss.Style) *Input {
	m.style.PlaceholderFocus = style

	return m
}

func (m *Input) SetCursorStyle(style lipgloss.Style) *Input {
	m.style.Cursor = style

	return m
}

func (m *Input) SetErrorStyle(style lipgloss.Style) *Input {
	m.style.Error = style

	return m
}

func (m *Input) updateStyle() {
	if m.state.focus {
		m.input.TextStyle = m.style.TextFocus
		m.input.PromptStyle = m.style.PromptFocus
		m.input.PlaceholderStyle = m.style.PlaceholderFocus
	} else {
		m.input.TextStyle = m.style.TextBase
		m.input.PromptStyle = m.style.PromptBase
		m.input.PlaceholderStyle = m.style.PlaceholderBase
	}
	m.input.Cursor.Style = m.style.Cursor
}
