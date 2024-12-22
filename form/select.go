package form

import (
	"github.com/MrSquaare/boba/component"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type SelectItemProps struct {
	Component component.Component
	Value     string
}

type SelectState struct {
	focus         bool
	selectedIndex int
}

type Select struct {
	items  []SelectItemProps
	inline bool
	state  SelectState
}

type SelectKeyMap struct {
	Prev key.Binding
	Next key.Binding
}

var (
	selectKeyMap = SelectKeyMap{
		Prev: key.NewBinding(
			key.WithKeys("left", "up"),
			key.WithHelp("left/up", "Previous selection"),
		),
		Next: key.NewBinding(
			key.WithKeys("right", "down"),
			key.WithHelp("right/down", "Next selection"),
		),
	}
)

func NewSelect(items []SelectItemProps) *Select {
	m := &Select{
		items:  items,
		inline: false,
		state: SelectState{
			focus:         false,
			selectedIndex: 0,
		},
	}

	return m
}

func (m *Select) Init() tea.Cmd {
	return m.updateItems()
}

func (m *Select) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	focus := m.state.focus
	selectedIndex := m.state.selectedIndex

	switch typedMsg := msg.(type) {
	case component.FocusMsg:
		m.state.focus = typedMsg.Focus
	case tea.KeyMsg:
		msg = nil

		switch {
		case key.Matches(typedMsg, selectKeyMap.Prev):
			m.state.selectedIndex--

			if m.state.selectedIndex < 0 {
				m.state.selectedIndex = len(m.items) - 1
			}
		case key.Matches(typedMsg, selectKeyMap.Next):
			msg = nil

			m.state.selectedIndex++

			if m.state.selectedIndex >= len(m.items) {
				m.state.selectedIndex = 0
			}
		}
	}

	if m.state.focus != focus || m.state.selectedIndex != selectedIndex {
		cmds = append(cmds, m.updateItems())
	}

	var cmd tea.Cmd

	m.items[m.state.selectedIndex].Component, cmd = m.items[m.state.selectedIndex].Component.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Select) View() string {
	var s string

	for i, item := range m.items {
		s += item.Component.View()

		if i < len(m.items)-1 {
			if m.inline {
				s += " "
			} else {
				s += "\n"
			}
		}
	}

	return s
}

func (m *Select) Keys() []key.Binding {
	return []key.Binding{selectKeyMap.Prev, selectKeyMap.Next}
}

func (m *Select) Value() string {
	return m.items[m.state.selectedIndex].Value
}

func (m *Select) SetInline(inline bool) *Select {
	m.inline = inline

	return m
}

func (m *Select) SetSelectedIndex(index int) *Select {
	if index < 0 {
		index = 0
	}

	if index >= len(m.items) {
		index = len(m.items) - 1
	}

	m.state.selectedIndex = index

	return m
}

func (m *Select) updateItems() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.items)*2)

	for i, item := range m.items {
		var focusCmd tea.Cmd
		var activeCmd tea.Cmd

		item.Component, focusCmd = item.Component.Update(component.FocusMsg{Focus: i == m.state.selectedIndex && m.state.focus})
		item.Component, activeCmd = item.Component.Update(component.ActiveMsg{Active: i == m.state.selectedIndex})

		cmds[i*2] = focusCmd
		cmds[i*2+1] = activeCmd
	}

	return tea.Batch(cmds...)
}
