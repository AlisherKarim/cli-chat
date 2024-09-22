package teamodels

import (
	"fmt"

	"github.com/alisherkarim/cli-chat/env"
	"github.com/alisherkarim/cli-chat/utils"
	tea "github.com/charmbracelet/bubbletea"
)

type ChatModel struct {
	env *env.Env
	messages []string
}

func CreateChatModel(env *env.Env, prevPage tea.Model) ChatModel {
	return ChatModel{
		env: env,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return nil
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c", "q":
					return m, tea.Quit
				case "enter":
					return m, nil
			}
	}
	return m, nil
}

func (m ChatModel) View() string {
	var s string

	if(m.env.GetUser().Username != "") {
		s += utils.SuccessStyle.Render(fmt.Sprintf("\nLogged in as %s", m.env.GetUser().Username))
		s += utils.SuccessStyle.Render(fmt.Sprintf("\nEmail: %s", m.env.GetUser().Email))
		s += utils.SuccessStyle.Render(fmt.Sprintf("\nAccess Token: %s\n", m.env.GetSession().AccessToken))
	}

	return s
}