package teamodels

import (
	"encoding/json"
	"fmt"

	"github.com/alisherkarim/cli-chat/env"
	"github.com/alisherkarim/cli-chat/types"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

// Room holds information about a single room
type Room struct {
	ID   string `json:"room_id"`
	Name string `json:"name"`
}

// ResponseData represents the structure of the JSON response
type ResponseData struct {
	Rooms []Room `json:"rooms"`
}

// ChatListModel represents the state of our Bubble Tea model
type ChatListModel struct {
	env 		  	*env.Env
	loading 		bool
	spinner 		spinner.Model
	data    		ResponseData
	selectedIndex 	int
	err     		error
}

// Init initializes the model, starts the spinner, and triggers data loading
func (m ChatListModel) Init() tea.Cmd {
	return nil
}

// loadData performs the backend request asynchronously
func (m* ChatListModel) loadData() tea.Msg {
	res, err := utils.RequestChatList(m.env.GetUser().Username)
	if err != nil {
		return types.ErrorMsg{Err: err}
	}
	return types.ResponseMsg{Res: res}
}

// Update handles messages and updates model state
func (m ChatListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c", "esc":
			return m, tea.Quit

		case "up":
			if len(m.data.Rooms) == 0 {
				return m, nil
			}
			ln := len(m.data.Rooms)
			m.selectedIndex = ((m.selectedIndex - 1) % ln + ln) % ln
			return m, nil

		case "down":
			if len(m.data.Rooms) == 0 {
				return m, nil
			}

			m.selectedIndex = (m.selectedIndex + 1) % len(m.data.Rooms)
			return m, nil
		
		case "enter":
			return m, nil

		default:
			return m, nil
		}

	case types.ResponseMsg:
		m.loading = false
		var err error
		m.data, err = ParseResponse(msg.Res)

		if err != nil {
			return m.Update(types.ErrorMsg{Err: err})
		}

		return m, nil

	case types.ErrorMsg:
		m.loading = false
		m.err = msg.Err
		return m, nil

	default:
		return m, nil
	}
}

// View renders the UI based on the model's state
func (m ChatListModel) View() string {
	var s string
	if m.loading {
		s = m.spinner.View() + selectedOptionStyle.Render(" Loading...")
	} else {
		if m.err != nil {
			s += fmt.Sprintf("something went wrong: %s\n", m.err)
		} else if len(m.data.Rooms) != 0 {
			s += fmt.Sprintf("total rooms: %d\n", len(m.data.Rooms))
			for idx, v := range m.data.Rooms {
				if idx == m.selectedIndex {
					s += selectedOptionStyle.Render(fmt.Sprintf("\n> %s", v.Name))
				} else {
					s += defaultOptionStyle.Render(fmt.Sprintf("\n  %s", v.Name))
				}
			}
		}
	}
	return s + "\n"
}

func CreateChatListModel(env *env.Env) ChatListModel {
	s := spinner.New()
	s.Spinner = spinner.Points
	s.Style = utils.LoadingSpinner

	model := ChatListModel{
		env: env,
		loading: true,
		spinner: s,
		data: ResponseData{Rooms: []Room{}},
	}

	return model
}

func ParseResponse(res string) (ResponseData, error) {
	var result ResponseData
	if err := json.Unmarshal([]byte(res), &result); err != nil {   // Parse []byte to go struct pointer
		return ResponseData{}, err
	}
	return result, nil
}