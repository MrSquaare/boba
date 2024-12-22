package form

import (
	"github.com/MrSquaare/boba/component"
	tea "github.com/charmbracelet/bubbletea"
)

type Hide struct {
	child component.Component
	hide  func() bool
}

func NewHide(child component.Component) *Hide {
	m := &Hide{
		child: child,
		hide:  func() bool { return true },
	}

	return m
}

func (m *Hide) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Hide) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.child, cmd = m.child.Update(msg)

	return m, cmd
}

func (m *Hide) View() string {
	if m.hide() {
		return ""
	}

	return m.child.View()
}

func (m *Hide) Child() component.Component {
	return m.child
}

func (m *Hide) Hide() bool {
	return m.hide()
}

func (m *Hide) SetHide(hide func() bool) *Hide {
	m.hide = hide

	return m
}
