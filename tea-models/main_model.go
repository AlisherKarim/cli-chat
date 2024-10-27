package teamodels

import (
	"fmt"

	"github.com/alisherkarim/cli-chat/env"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)


type MainModel struct {
	env *env.Env
	currentFocusedModelName string
	chatListModel tea.Model
	chatModel ChatModel
	spinner spinner.Model

}

func CreateMainModel(env *env.Env) MainModel {
	m := MainModel{
		env: env,
		currentFocusedModelName: "chat",
		chatModel: CreateChatModel(env, nil),
	}
	m.spinner = spinner.New()
	m.currentFocusedModelName = "Chat"

	return m
}

func (m MainModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick)
}

func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c":
					return m, tea.Quit
				default:
					f, cmd := m.chatModel.Update(msg)
					m.chatModel = f.(ChatModel)
					cmds = append(cmds, cmd)
			}
		case spinner.TickMsg:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
		case tea.WindowSizeMsg:
			utils.ChatTabStyle = utils.ChatTabStyle.Height(msg.Height - 15).Width(msg.Width - 25)
			utils.UserInfoBoxStyle = utils.UserInfoBoxStyle.Width(msg.Width - 3)
	}
	return m, tea.Batch(cmds...)
}

func (m MainModel) View() string {
	var s string
	var info string

	if(m.env.GetUser().Username != "") {
		info += utils.SuccessStyle.Render(fmt.Sprintf("\nLogged in as %s", m.env.GetUser().Username))
		info += utils.SuccessStyle.Render(fmt.Sprintf("\nEmail: %s", m.env.GetUser().Email))
	}

	s += utils.UserInfoBoxStyle.Render(info)
	s += "\n\n"

	chatList := utils.ChatsListTabStyle.Render(fmt.Sprintf("%4s", "chat list here"))

	if m.currentFocusedModelName == "ChatList"{
		chatList = utils.ChatsListTabStyle.BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("204")).Render(fmt.Sprintf("%4s", "chat list here"))
	}

	s += lipgloss.JoinHorizontal(lipgloss.Top, chatList, m.chatModel.View())
	s += "\n"
	s += utils.HelpStyle.Render(fmt.Sprintf("tab: focus next • n: new %s • q: exit", m.currentFocusedModelName))

	return s
}

func (m *MainModel) switchTab() {
	if m.currentFocusedModelName == "Chat" {
		m.currentFocusedModelName = "ChatList"
		m.chatModel.Blur()
	} else {
		m.currentFocusedModelName = "Chat"
		m.chatModel.Focus()
	}
}