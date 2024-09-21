package teamodels

import (
	"fmt"

	"github.com/alisherkarim/cli-chat/types"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppState struct {
	test string
	currentStateName string
	currentState tea.Model
	currentFocusedTab string
	currentUser types.User
	currentSession types.Session
}

var (
	helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

func (state AppState) Init() tea.Cmd {
	return nil
}

func (state AppState) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c", "q":
					return state, tea.Quit
				case "s":
					state.currentStateName = "login"
					return CreateLoginModel(state), nil
				case "f":
					state.test = "haha"
					return state, nil
		}
	}
	return state, nil
}

func (state AppState) View() string {
	var s string
	
	s += fmt.Sprintf("Logged in as %s\n", state.test)
	s += helpStyle.Render(fmt.Sprintf("\ntab: focus next • q: exit\n"))
	return s
}