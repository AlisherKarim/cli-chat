package teamodels

// A simple example demonstrating the use of multiple text input components
// from the Bubbles component library.

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alisherkarim/cli-chat/env"
	"github.com/alisherkarim/cli-chat/types"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type LoginModel struct {
	env *env.Env
	prevPage tea.Model
	focusIndex int
	inputs     []textinput.Model
	errorMessage string
	isRequesting bool
	loadingSpinner  spinner.Model
}

func CreateLoginModel(env *env.Env, prevPage tea.Model) LoginModel {
	m := LoginModel{
		env: env,
		prevPage: prevPage,
		inputs: make([]textinput.Model, 2),
		loadingSpinner: spinner.New(),
	}

	m.loadingSpinner.Spinner = spinner.Points
	m.loadingSpinner.Style = utils.LoadingSpinner


	var t textinput.Model
	for i := range m.inputs {
		t = textinput.New()
		t.Cursor.Style = utils.CursorStyle
		t.CharLimit = 32

		switch i {
		case 0:
			t.Placeholder = "Username"
			t.Focus()
			t.PromptStyle = utils.FocusedStyle
			t.TextStyle = utils.FocusedStyle
		case 1:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m LoginModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadingSpinner.Tick,
		textinput.Blink,
	)
}

func (m LoginModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c":
					return m, tea.Quit
				case "esc":
					return m.prevPage, nil
				case "tab", "shift+tab", "enter", "up", "down":
					s := msg.String()

					if s == "enter" && m.focusIndex == len(m.inputs) {
						if m.CheckInputs() {
							m.isRequesting = true
							go m.RequestLogin()
						}
						return m, nil
					}

					// Cycle indexes
					if s == "up" || s == "shift+tab" {
						m.focusIndex--
					} else {
						m.focusIndex++
					}

					if m.focusIndex > len(m.inputs) {
						m.focusIndex = 0
					} else if m.focusIndex < 0 {
						m.focusIndex = len(m.inputs)
					}

					cmds := make([]tea.Cmd, len(m.inputs))
					for i := 0; i <= len(m.inputs)-1; i++ {
						if i == m.focusIndex {
							// Set focused state
							cmds[i] = m.inputs[i].Focus()
							m.inputs[i].PromptStyle = utils.FocusedStyle
							m.inputs[i].TextStyle = utils.FocusedStyle
							continue
						}
						// Remove focused state
						m.inputs[i].Blur()
						m.inputs[i].PromptStyle = utils.NoStyle
						m.inputs[i].TextStyle = utils.NoStyle
					}

					return m, tea.Batch(cmds...)
				}
		case spinner.TickMsg:
			var cmd tea.Cmd
			m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
			return m, cmd
		case types.ResponseMsg:
			m.isRequesting = false
			return CreateChatListModel(m.env), nil
		case types.ErrorMsg:
			m.errorMessage = msg.Err.Error()
			m.isRequesting = false
			return m, nil
	}

	// Handle character input and blinking
	cmd := m.updateInputs(msg)
	cmd = tea.Batch(cmd, m.loadingSpinner.Tick)

	return m, cmd
}

func (m *LoginModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m LoginModel) View() string {
	var b strings.Builder

	b.WriteString(utils.PageNameStyle.Render("\nLogin Page\n"))
	b.WriteRune('\n')
	b.WriteRune('\n')

	for i := range m.inputs {
		b.WriteString(m.inputs[i].View())
		if i < len(m.inputs)-1 {
			b.WriteRune('\n')
		}
	}

	button := &utils.BlurredButton
	if m.focusIndex == len(m.inputs) {
		button = &utils.FocusedButton
	}
	fmt.Fprintf(&b, "\n\n%s\n\n", *button)

	if m.isRequesting {
		b.WriteString(fmt.Sprintf("\n%s", m.loadingSpinner.View()))
	}

	if m.errorMessage != "" {
		b.WriteString(utils.ErrorStyle.Render(fmt.Sprintf("\n%s", m.errorMessage)))
	}

	b.WriteRune('\n')
	b.WriteString(utils.HelpStyle.Render("esc to go back"))

	return b.String()
}


// Helper methods

func (m *LoginModel) CheckInputs() bool {
	for _, v := range m.inputs {
		if strings.TrimSpace(v.Value()) == "" {
			m.errorMessage = "Please fill all the fields!"
			return false
		}
	}

	m.errorMessage = ""
	return true
}

func (m *LoginModel) RequestLogin() {
	res, err := utils.Login(m.inputs[0].Value(), m.inputs[1].Value())
	if err != nil {
		m.errorMessage = "Something went wrong while requesting. Please, try again"
		m.env.CurrentProgram.Send(types.ErrorMsg{Err: err})
		return
	}
	
	var data map[string]string
	err = json.Unmarshal([]byte(res), &data)
	if err != nil {
		m.errorMessage = "Something went wrong while requesting. Please, try again"
		m.env.CurrentProgram.Send(types.ErrorMsg{Err: err})
		return
	}

	m.env.SetUser(data["username"], data["email"])
	m.env.SetSession(data["access_token"])
	m.env.CurrentProgram.Send(types.ResponseMsg{Res: res})
}