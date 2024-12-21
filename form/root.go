package form

import (
	"github.com/charmbracelet/bubbles/key"
)

type WithKeys interface {
	Keys() []key.Binding
}

type WithSkip interface {
	Skip() bool
}

type WithHide interface {
	Hide() bool
}

type WithValidation interface {
	Validate() bool
	Error() error
}

type WithValue interface {
	Value() string
}
