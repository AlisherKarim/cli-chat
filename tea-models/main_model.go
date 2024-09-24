package teamodels

import (
	"fmt"
	"time"

	"github.com/alisherkarim/cli-chat/env"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// sessionState is used to track which model is focused
type sessionState uint

const (
	defaultTime              = time.Minute
	timerView   sessionState = iota
	spinnerView
)

type MainModel struct {
	env *env.Env
	currentFocusedModelName string
	chatListModel tea.Model
	chatModel tea.Model
	spinner spinner.Model
}

func CreateMainModel(env *env.Env) MainModel {
	m := MainModel{
		env: env,
		currentFocusedModelName: "chat",
		chatModel: CreateChatModel(env, nil),
	}
	m.spinner = spinner.New()
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
				case "ctrl+c", "q":
					return m, tea.Quit
				case "tab":
					m.switchTab()
					return m, nil
				default:
					m.chatModel, cmd = m.chatModel.Update(msg)
					cmds = append(cmds, cmd)
			}
		case spinner.TickMsg:
			m.spinner, cmd = m.spinner.Update(msg)
			cmds = append(cmds, cmd)
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

	s += lipgloss.JoinHorizontal(lipgloss.Top, utils.ChatsListTabStyle.Render(fmt.Sprintf("%4s", "chat list here")), m.chatModel.View())
	s += "\n"
	s += utils.HelpStyle.Render(fmt.Sprintf("tab: focus next • n: new %s • q: exit", m.currentFocusedModelName))

	return s
}

func (m *MainModel) switchTab() {
	if m.currentFocusedModelName == "Chat" {
		m.currentFocusedModelName = "ChatList"
	} else {
		m.currentFocusedModelName = "Chat"
	}
}