package form

import (
	"github.com/MrSquaare/boba/component"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
)

type FormItem struct {
	Name      string
	Component component.Component
}

type FormState struct {
	selectedIndex int
	step          int
	completed     bool
}

type Form struct {
	items []FormItem
	state FormState
}

type FormKeyMap struct {
	Prev key.Binding
	Next key.Binding
}

var (
	formKeyMap = FormKeyMap{
		Prev: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "Previous"),
		),
		Next: key.NewBinding(
			key.WithKeys("enter", "tab"),
			key.WithHelp("enter/tab", "Next"),
		),
	}
)

func NewForm(items []FormItem) *Form {
	m := &Form{
		items: items,
		state: FormState{
			selectedIndex: 0,
			step:          0,
			completed:     false,
		},
	}

	return m
}

func (m *Form) Init() tea.Cmd {
	return tea.Batch(m.initItems(), m.updateItems())
}

func (m *Form) Update(msg tea.Msg) (*Form, tea.Cmd) {
	cmds := []tea.Cmd{}

	selectedIndex := m.state.selectedIndex

	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(typedMsg, formKeyMap.Prev):
			msg = nil

			m.state.completed = false

			for i := m.state.selectedIndex - 1; i >= 0; i-- {
				if !isSkip(m.items[i].Component) {
					m.state.selectedIndex = i

					break
				}
			}
		case key.Matches(typedMsg, formKeyMap.Next):
			msg = nil

			if withValidation, ok := withValidation(m.items[m.state.selectedIndex].Component); ok {
				if !withValidation.Validate() {
					break
				}
			}

			for i := m.state.selectedIndex + 1; i <= len(m.items); i++ {
				if i == len(m.items) {
					m.state.completed = true

					break
				} else if !isSkip(m.items[i].Component) {
					m.state.selectedIndex = i

					break
				}
			}

			if m.state.step < m.state.selectedIndex {
				m.state.step = m.state.selectedIndex
			}
		}
	}

	if m.state.selectedIndex != selectedIndex {
		cmds = append(cmds, m.updateItems())
	}

	var cmd tea.Cmd

	m.items[m.state.selectedIndex].Component, cmd = m.items[m.state.selectedIndex].Component.Update(msg)

	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m *Form) View() string {
	var s string

	for i, item := range m.items {
		if i > m.state.step {
			continue
		}

		if withHide, ok := withHide(item.Component); ok {
			if withHide.Hide() {
				continue
			}
		}

		if i > 0 {
			s += "\n\n"
		}

		s += item.Component.View()
	}

	return s
}

func (m *Form) Keys() []key.Binding {
	keys := []key.Binding{}
	selectedItem := m.items[m.state.selectedIndex]

	if withKeys, ok := withKeys(selectedItem.Component); ok {
		keys = withKeys.Keys()
	}

	if m.state.selectedIndex > 0 {
		keys = append(keys, formKeyMap.Prev)
	}

	keys = append(keys, formKeyMap.Next)

	return keys
}

func (m *Form) Error(name string) error {
	for _, item := range m.items {
		if item.Name == name {
			if withValidation, ok := withValidation(item.Component); ok {
				return withValidation.Error()
			} else {
				return nil
			}
		}
	}

	return nil
}

func (m *Form) Errors() map[string]error {
	errors := make(map[string]error, len(m.items))

	for _, item := range m.items {
		errors[item.Name] = m.Error(item.Name)
	}

	return errors
}

func (m *Form) Value(name string) any {
	for _, item := range m.items {
		if item.Name == name {
			if withValue, ok := withValue(item.Component); ok {
				return withValue.Value()
			} else {
				return nil
			}
		}
	}

	return nil
}

func (m *Form) Values() map[string]any {
	values := make(map[string]any, len(m.items))

	for _, item := range m.items {
		values[item.Name] = m.Value(item.Name)
	}

	return values
}

func (m *Form) SetSelectedIndex(index int) *Form {
	m.state.selectedIndex = index

	if m.state.step < m.state.selectedIndex {
		m.state.step = m.state.selectedIndex
	}

	return m
}

func (m *Form) SetStep(step int) *Form {
	m.state.step = step

	return m
}

func (m *Form) Completed() bool {
	return m.state.completed
}

func (m *Form) initItems() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.items))

	for i, item := range m.items {
		cmds[i] = item.Component.Init()
	}

	return tea.Batch(cmds...)
}

func (m *Form) updateItems() tea.Cmd {
	cmds := make([]tea.Cmd, len(m.items))

	for i, item := range m.items {
		var cmd tea.Cmd

		item.Component, cmd = item.Component.Update(component.FocusMsg{Focus: i == m.state.selectedIndex})

		cmds[i] = cmd
	}

	return tea.Batch(cmds...)
}

func withKeys(m component.Component) (WithKeys, bool) {
	if m, ok := m.(WithKeys); ok {
		return m, ok
	}

	if m, ok := m.(component.WithChild); ok {
		return withKeys(m.Child())
	}

	return nil, false
}

func withSkip(m component.Component) (WithSkip, bool) {
	if m, ok := m.(WithSkip); ok {
		return m, ok
	}

	if m, ok := m.(component.WithChild); ok {
		return withSkip(m.Child())
	}

	return nil, false
}

func withHide(m component.Component) (WithHide, bool) {
	if m, ok := m.(WithHide); ok {
		return m, ok
	}

	if m, ok := m.(component.WithChild); ok {
		return withHide(m.Child())
	}

	return nil, false
}

func isSkip(m component.Component) bool {
	if withSkip, ok := withSkip(m); ok {
		return withSkip.Skip()
	}

	if withHide, ok := withHide(m); ok {
		return withHide.Hide()
	}

	return false
}

func withValidation(m component.Component) (WithValidation, bool) {
	if m, ok := m.(WithValidation); ok {
		return m, ok
	}

	if m, ok := m.(component.WithChild); ok {
		return withValidation(m.Child())
	}

	return nil, false
}

func withValue(m component.Component) (WithValue, bool) {
	if m, ok := m.(WithValue); ok {
		return m, ok
	}

	if m, ok := m.(component.WithChild); ok {
		return withValue(m.Child())
	}

	return nil, false
}
