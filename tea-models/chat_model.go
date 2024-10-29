package teamodels

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/alisherkarim/cli-chat/env"
	"github.com/alisherkarim/cli-chat/types"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/websocket"
)

type ChatModel struct {
	env *env.Env
	loading 		bool
	messages    []string
	err         error
	// conn				*websocket.Conn

	// ui
	spinner 		spinner.Model
	viewport    viewport.Model
	textarea    textarea.Model
	senderStyle lipgloss.Style
	myStyle lipgloss.Style
}


type SocketMsg struct {
	Msg string
}

type ConnectionMsg struct {
	Conn *websocket.Conn
}

// Message represents the structure of a WebSocket message.
type Message struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
	Type    string `json:"type"`  // e.g., "text", "image", "notification"
}


var conn *websocket.Conn

func CreateChatModel(env *env.Env) ChatModel {
	ta := textarea.New()
	ta.Placeholder = "Send a message..."
	ta.Focus()

	ta.Prompt = "┃ "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(3)

	// Remove cursor line styling
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 10)
	vp.SetContent(`Welcome to the chat room!
Type a message and press Enter to send.`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = utils.LoadingSpinner

	return ChatModel{
		env: 					env,
		loading: 			true,
		spinner: 			s,
		textarea:    	ta,
		messages:    	[]string{},
		viewport:    	vp,
		senderStyle: 	lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		myStyle: 			lipgloss.NewStyle().Foreground(lipgloss.Color("205")),
		err:         	nil,
	}
}

func (m ChatModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {

	case SocketMsg:
		v, err := ProcessMessage([]byte(msg.Msg))

		if err != nil {
			return m.Update(types.ErrorMsg{Err: err})
		}
		
		m.err = nil
		var sender string
		if v.Sender == m.env.GetUser().Username {
			sender = m.senderStyle.Render(v.Sender + ": ")
		} else {
			sender = m.myStyle.Render(v.Sender + ": ")
		}
		m.messages = append(m.messages, sender + v.Content)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()
		return m, nil

	case Message:
		m.messages = append(m.messages, m.senderStyle.Render(msg.Sender + ": ") + msg.Content)
		m.viewport.SetContent(strings.Join(m.messages, "\n"))
		m.textarea.Reset()
		m.viewport.GotoBottom()
		return m, nil

	case tea.KeyMsg:

		switch msg.Type {

		case tea.KeyCtrlC:
			fmt.Println(m.textarea.Value())
			// if conn != nil {
			// 	conn.Close()
			// }
			return m, tea.Quit

		case tea.KeyEsc:
			// if conn != nil {
			// 	conn.Close()
			// }
			chatListModel := CreateChatListModel(m.env)
			return chatListModel, tea.Batch(chatListModel.loadData, chatListModel.spinner.Tick)

		case tea.KeyEnter:
			if !m.textarea.Focused() {
				return m, nil
			}

			// m.messages = append(m.messages, m.senderStyle.Render("You: ")+m.textarea.Value())
			// m.viewport.SetContent(strings.Join(m.messages, "\n"))
			// m.textarea.Reset()
			// m.viewport.GotoBottom()

			payload := map[string]string{
        "sender": m.env.GetUser().Username,
        "content": m.textarea.Value(),
				"type": "text",
   		}
			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				return m.Update(types.ErrorMsg{Err: err})
			}

			// message := Message{
			// 	Sender: m.env.GetUser().Username,
			// 	Content: m.textarea.Value(),
			// 	Type: "text",
			// }

			if err := conn.WriteMessage(websocket.TextMessage, payloadBytes); err != nil {
				return m.Update(types.ErrorMsg{Err: err})
			}

			return m, nil
		}

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case ConnectionMsg:
		// conn = msg.Conn
		m.loading = false
		go m.listenToWebSocket()

	case types.ErrorMsg:
		m.loading = false
		// if conn != nil {
		// 	conn.Close()
		// }
		m.err = msg.Err
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m ChatModel) View() string {
	var s string

	if m.loading {
		return m.spinner.View() + selectedOptionStyle.Render(" Loading...")
	}

	s += fmt.Sprintf(
		"%s\n\n%s",
		m.viewport.View(),
		m.textarea.View(),
	) + "\n\n"

	if m.err != nil {
		s += utils.ErrorStyle.Render(m.err.Error() + "\n")
	}

	return s
}

func (m *ChatModel) Blur() {
	m.textarea.Blur()
}

func (m *ChatModel) Focus() {
	m.textarea.Focus()
}

func (m *ChatModel) ConnectToChat(id string) tea.Cmd {
	url := fmt.Sprintf("ws://localhost:8080/api/v1/chats/%s/ws", id)
	dialer := websocket.Dialer{
    HandshakeTimeout: 10 * time.Second,
	}

	c, _, err := dialer.Dial(url, nil)
	if err != nil {
		return func() tea.Msg {return types.ErrorMsg{Err: err}}
	}

	conn = c

	return func() tea.Msg {
		return ConnectionMsg{Conn: c}
	}
}

func (m *ChatModel) listenToWebSocket() {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			m.env.CurrentProgram.Send(types.ErrorMsg{Err: err})
			return
		}
		m.env.CurrentProgram.Send(SocketMsg{Msg: string(message)})
	}
}

// ProcessMessage processes a raw message, applies validation or custom logic.
func ProcessMessage(rawMessage []byte) (*Message, error) {
	var message Message
	err := json.Unmarshal(rawMessage, &message)
	if err != nil {
		return nil, err
	}

	// Apply custom logic here, like filtering content or routing based on type
	return &message, nil
}
