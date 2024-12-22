package component

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Component = tea.Model

type WithChild interface {
	Child() Component
}

type FocusMsg struct {
	Focus bool
}

type ActiveMsg struct {
	Active bool
}
