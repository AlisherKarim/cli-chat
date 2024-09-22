package teamodels

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

type RegisterModel struct {
	env *env.Env
	prevPage tea.Model
	focusIndex int
	inputs     []textinput.Model
	errorMessage string
	isRequesting bool
	loadingSpinner  spinner.Model
}

func CreateRegisterModel(env *env.Env, prevPage tea.Model) RegisterModel {
	m := RegisterModel{
		env: env,
		prevPage: prevPage,
		inputs: make([]textinput.Model, 3),
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
			t.Placeholder = "Nickname"
			t.Focus()
			t.PromptStyle = utils.FocusedStyle
			t.TextStyle = utils.FocusedStyle
		case 1:
			t.Placeholder = "Email"
			t.CharLimit = 64
		case 2:
			t.Placeholder = "Password"
			t.EchoMode = textinput.EchoPassword
			t.EchoCharacter = '•'
		}

		m.inputs[i] = t
	}

	return m
}

func (m RegisterModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadingSpinner.Tick,
		textinput.Blink,
	)
}

func (m RegisterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "esc":
			return m.prevPage, nil

		// Set focus to next input
		case "tab", "shift+tab", "enter", "up", "down":
			s := msg.String()

			// Did the user press enter while the submit button was focused?
			// If so, exit.
			if s == "enter" && m.focusIndex == len(m.inputs) {
				if m.CheckInputs() {
					m.isRequesting = true
					go m.RequestRegister()
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
		return CreateLoginModel(m.env, m), nil
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

func (m *RegisterModel) updateInputs(msg tea.Msg) tea.Cmd {
	cmds := make([]tea.Cmd, len(m.inputs))

	// Only text inputs with Focus() set will respond, so it's safe to simply
	// update all of them here without any further logic.
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return tea.Batch(cmds...)
}

func (m RegisterModel) View() string {
	var b strings.Builder

	b.WriteString(utils.PageNameStyle.Render("\nRegister Page\n"))
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

	b.WriteString(utils.HelpStyle.Render("\n\nesc to go back"))

	return b.String()
}


// Helper methods

func (m *RegisterModel) CheckInputs() bool {
	for _, v := range m.inputs {
		if strings.TrimSpace(v.Value()) == "" {
			m.errorMessage = "Please fill all the fields!"
			return false
		}
	}

	m.errorMessage = ""
	return true
}

func (m *RegisterModel) RequestRegister() {
	res, err := utils.Register(m.inputs[0].Value(), m.inputs[1].Value(), m.inputs[2].Value())
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
	m.env.CurrentProgram.Send(types.ResponseMsg{Res: res})
}