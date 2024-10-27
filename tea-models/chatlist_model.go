package teamodels

import (
	"fmt"
	"io"
	"strings"

	"github.com/alisherkarim/cli-chat/env"
	"github.com/alisherkarim/cli-chat/types"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)




const listHeight = 14
const listWidth = 14

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("> " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}




type ChatListModel struct {
	env *env.Env
	// chatList []types.ChatListItem
	chatList list.Model
	loadingSpinner spinner.Model
	requestError string
	response string
}

func CreateChatListModel(env *env.Env) ChatListModel {
	items := []list.Item{
		item("Chatroom 1"),
		item("Another chatroom"),
		item("Max"),
		item("Jane"),
		item("Curry Chat"),
		item("Testers"),
		item("Pasta"),
		item("Fillet Mignon"),
		item("Caviar"),
		item("Just Wine"),
	}
	m := ChatListModel{
		env: env,
		chatList: list.New(items, itemDelegate{}, listWidth, listHeight),
		loadingSpinner: spinner.New(),
	}
	m.loadingSpinner.Spinner = spinner.Points
	m.loadingSpinner.Style = utils.LoadingSpinner
	return m
}


// bubble tea methods

func (m ChatListModel) Init() tea.Cmd {
	return tea.Batch(m.loadingSpinner.Tick)
}

func (m ChatListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var chCmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.chatList.SetWidth(msg.Width)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.loadingSpinner, cmd = m.loadingSpinner.Update(msg)
		return m, cmd
	case types.ErrorMsg:
		m.requestError = msg.Err.Error()
		return m, nil
	case types.ResponseMsg:
		m.response = msg.Res
		return m, nil
	}

	m.chatList, chCmd = m.chatList.Update(msg)
	// m.spinner, spCmd = m.spinner.Update(msg)
	return m, tea.Batch(m.loadingSpinner.Tick, chCmd)
}

func (m ChatListModel) View() string {
	s := "\n" + m.chatList.View() + "\n\n" + m.loadingSpinner.View()
	if m.requestError != "" {
		s = s + "\n" + m.requestError
	} else {
		s = s + "\n" + m.response
	}
	return s
}

// Helper methods

func (m *ChatListModel) RequestChatListData() tea.Cmd {
	return func() tea.Msg {
		res, err := utils.RequestChatList(m.env.GetUser().Username)
		if err != nil {
			return types.ErrorMsg{Err: err}
		}
		return types.ResponseMsg{Res: res}
	}
}