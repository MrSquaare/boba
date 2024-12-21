package main

import (
	"fmt"
	"os"
	"time"

	"github.com/MrSquaare/boba/component"
	"github.com/MrSquaare/boba/form"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-playground/validator/v10"
)

type Status = int

const (
	StatusForm Status = iota
	StatusLoading
	StatusSuccess
	StatusError
)

type Model struct {
	form    *form.Form
	spinner spinner.Model
	help    help.Model
	status  Status
}

type KeyMap struct {
	Exit key.Binding
	Help key.Binding
}

type loadingMsg struct{}

var (
	headerStyle = lipgloss.NewStyle().
			Foreground((lipgloss.Color("15"))).
			Bold(true).
			Background(lipgloss.Color("33")).
			Padding(0, 2)
	keyMap = KeyMap{
		Exit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "Exit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "Help"),
		),
	}
)

func newModel() Model {
	var myForm *form.Form

	myForm = form.NewForm([]form.FormItem{
		{
			Name: "host",
			Component: form.NewField(
				"Enter the server host",
				form.NewInput().
					SetValidateFunc(func(s string) error {
						validate := validator.New()

						if err := validate.Var(s, "required"); err != nil {
							return fmt.Errorf("host is required")
						}

						if err := validate.Var(s, "ip"); err != nil {
							return fmt.Errorf("host must be an IP address")
						}

						return nil
					}),
			),
		},
		{
			Name: "port",
			Component: form.NewField(
				"Enter the server port",
				form.NewInput().
					SetValidateFunc(func(s string) error {
						validate := validator.New()

						if err := validate.Var(s, "required"); err != nil {
							return fmt.Errorf("port is required")
						}

						if err := validate.Var(s, "number"); err != nil {
							return fmt.Errorf("port must be a number")
						}

						return nil
					}),
			),
		},
		{
			Name: "user",
			Component: form.NewField(
				"Enter the auth user",
				form.NewInput().
					SetValidateFunc(func(s string) error {
						validate := validator.New()

						if err := validate.Var(s, "required"); err != nil {
							return fmt.Errorf("user is required")
						}

						return nil
					}),
			),
		},
		{
			Name: "auth",
			Component: form.NewField(
				"Select an auth method",
				form.NewSelect([]form.SelectItemProps{
					{
						Value:     "key",
						Component: component.NewButton("Key"),
					},
					{
						Value:     "password",
						Component: component.NewButton("Password"),
					},
				}).
					SetInline(true),
			),
		},
		{
			Name: "key",
			Component: form.NewHide(
				form.NewField(
					"Select the auth key",
					form.NewLoader(
						func() component.Component {
							values := getKeyPaths()

							options := make([]form.SelectItemProps, len(values))

							for _, v := range values {
								options = append(options, form.SelectItemProps{
									Value:     v,
									Component: component.NewOption(v),
								})
							}

							return form.NewSelect(options)
						},
					).
						SetBindings(func() any {
							return myForm.Value("password")
						}),
				),
			).
				SetHide(func() bool {
					return myForm.Value("auth") != "key"
				}),
		},
		{
			Name: "password-note",
			Component: form.NewHide(
				form.NewSkip(
					component.NewText("Note: Passwords are not recommended for use in production environments."),
				).
					SetSkip(func() bool {
						return true
					}),
			).
				SetHide(func() bool {
					return myForm.Value("auth") != "password"
				}),
		},
		{
			Name: "password",
			Component: form.NewHide(
				form.NewField(
					"Enter the auth password",
					form.NewInput().
						SetValidateFunc(func(s string) error {
							validate := validator.New()

							if err := validate.Var(s, "required"); err != nil {
								return fmt.Errorf("password is required")
							}

							return nil
						}),
				),
			).
				SetHide(func() bool {
					return myForm.Value("auth") != "password"
				}),
		},
	})

	m := Model{
		form:    myForm,
		spinner: spinner.New(),
		help:    help.New(),
		status:  StatusForm,
	}

	return m
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.form.Init(), m.spinner.Tick)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}

	if m.status == StatusSuccess || m.status == StatusError {
		return m, tea.Quit
	}

	switch typedMsg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(typedMsg, keyMap.Exit):
			msg = nil

			m.status = StatusError
		}
	}

	switch m.status {
	case StatusForm:
		switch typedMsg := msg.(type) {
		case tea.KeyMsg:
			switch {
			case key.Matches(typedMsg, keyMap.Help):
				msg = nil

				m.help.ShowAll = !m.help.ShowAll
			}
		}

		var cmd tea.Cmd

		m.form, cmd = m.form.Update(msg)

		cmds = append(cmds, cmd)

		if m.form.Completed() {
			m.status = StatusLoading

			cmds = append(cmds, func() tea.Msg {
				connectToServer()

				return loadingMsg{}
			})
			cmds = append(cmds, m.spinner.Tick)
		}
	case StatusLoading:
		switch msg.(type) {
		case loadingMsg:
			msg = nil

			m.status = StatusSuccess
		}

		var cmd tea.Cmd

		m.spinner, cmd = m.spinner.Update(msg)

		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	var s string

	s += headerStyle.Render("examples: simple")

	if m.status == StatusForm {
		s += fmt.Sprintf("\n\n%s", m.form.View())

		s += fmt.Sprintf("\n\n%s", m.help.View(m))
	}

	if m.status == StatusLoading {
		s += fmt.Sprintf("\n\n%s", m.spinner.View())
	}

	if m.status == StatusSuccess {
		s += "\n\nSuccessfully connected to the server!"
	}

	if m.status == StatusError {
		s += "\n\nOperation cancelled."
	}

	s += "\n"

	return s
}

func (m Model) ShortHelp() []key.Binding {
	return []key.Binding{keyMap.Help, keyMap.Exit}
}

func (m Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		m.form.Keys(),
		{keyMap.Help, keyMap.Exit},
	}
}

func getKeyPaths() []string {
	time.Sleep(1 * time.Second)

	return []string{
		"/path/to/key-1",
		"/path/to/key-2",
		"/path/to/key-3",
	}
}

func connectToServer() error {
	time.Sleep(1 * time.Second)

	return nil
}

func main() {
	if _, err := tea.NewProgram(newModel()).Run(); err != nil {
		fmt.Printf("could not start program: %s\n", err)
		os.Exit(1)
	}
}
