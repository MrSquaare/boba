package form

import (
	"github.com/MrSquaare/boba/component"
	tea "github.com/charmbracelet/bubbletea"
)

type Skip struct {
	child component.Component
	skip  func() bool
}

func NewSkip(child component.Component) *Skip {
	m := &Skip{
		child: child,
		skip:  func() bool { return true },
	}

	return m
}

func (m *Skip) Init() tea.Cmd {
	return m.child.Init()
}

func (m *Skip) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	m.child, cmd = m.child.Update(msg)

	return m, cmd
}

func (m *Skip) View() string {
	return m.child.View()
}

func (m *Skip) Child() component.Component {
	return m.child
}

func (m *Skip) Skip() bool {
	return m.skip()
}

func (m *Skip) SetSkip(skip func() bool) *Skip {
	m.skip = skip

	return m
}
