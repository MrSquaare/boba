package form

import (
	md5 "crypto/md5"
	"encoding/json"

	"github.com/MrSquaare/boba/component"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type loaderMsg struct {
	id    int
	hash  [16]byte
	child component.Component
}

type Loader struct {
	id           int
	childFn      func() component.Component
	bindings     func() any
	bindingsHash [16]byte
	child        component.Component
	focus        bool
	loading      bool
	spinner      spinner.Model
}

func NewLoader(childFn func() component.Component) *Loader {
	m := &Loader{
		id:           0,
		childFn:      childFn,
		bindings:     func() any { return nil },
		bindingsHash: [16]byte{},
		child:        nil,
		focus:        false,
		loading:      true,
		spinner:      spinner.New(),
	}

	return m
}

func (m *Loader) Init() tea.Cmd {
	return nil
}

func (m *Loader) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	bindingsHash := m.bindingsHash
	child := m.child
	focus := m.focus

	switch typedMsg := msg.(type) {
	case loaderMsg:
		if typedMsg.id == m.id && typedMsg.hash == m.bindingsHash {
			m.child = typedMsg.child
			m.loading = false
		}
	case component.FocusMsg:
		m.focus = typedMsg.Focus
	}

	if m.focus != focus {
		if m.focus {
			m.bindingsHash = hash(m.bindings())
		}
	}

	if bindingsHash != m.bindingsHash {
		m.loading = true

		cmds = append(cmds, func() tea.Msg {
			child := m.childFn()

			return loaderMsg{
				id:    m.id,
				hash:  m.bindingsHash,
				child: child,
			}
		})
		cmds = append(cmds, m.spinner.Tick)
	}

	if m.loading {
		var cmd tea.Cmd

		m.spinner, cmd = m.spinner.Update(msg)

		cmds = append(cmds, cmd)
	} else if m.child != nil {
		if m.child != child {
			cmds = append(cmds, m.child.Init())
		}

		if m.child != child || m.focus != focus {
			var cmd tea.Cmd

			m.child, cmd = m.child.Update(component.FocusMsg{Focus: m.focus})

			cmds = append(cmds, cmd)
		}

		var cmd tea.Cmd

		m.child, cmd = m.child.Update(msg)

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m *Loader) View() string {
	if m.loading {
		return m.spinner.View()
	}

	if m.child != nil {
		return m.child.View()
	}

	return "No component"
}

func (m *Loader) Child() component.Component {
	return m.child
}

func (m *Loader) Bindings() any {
	return m.bindings()
}

func (m *Loader) SetBindings(bindings func() any) *Loader {
	m.bindings = bindings

	return m
}

func hash(value interface{}) [16]byte {
	data, _ := json.Marshal(value)

	return md5.Sum(data)
}
