package teamodels

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EntryModel struct {
	options []string 
	selectedOptionIndex int
}

var (
	selectedOptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204"))
	defaultOptionStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
)

func CreateEntryModel() EntryModel {
	return EntryModel{
		options: []string{"Login", "Register"},
		selectedOptionIndex: 0,
	}
}

func (state EntryModel) Init() tea.Cmd {
	return nil
}

func (state EntryModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
				case "ctrl+c", "q":
					return state, tea.Quit
				case "up":
					state.selectedOptionIndex = (state.selectedOptionIndex - 1) % len(state.options)
					return state, nil
				case "down":
					state.selectedOptionIndex = (state.selectedOptionIndex + 1) % len(state.options)
					return state, nil
				case "enter":
					return state.OpenSelectedModel(), nil
			}
	}
	return state, nil
}

func (state EntryModel) View() string {
	var s string

	for idx, v := range state.options {
		if idx == state.selectedOptionIndex {
			s += selectedOptionStyle.Render(fmt.Sprintf("\n> %s", v))
		} else {
			s += defaultOptionStyle.Render(fmt.Sprintf("\n  %s", v))
		}
	}

	s += "\n\n Select option\n"

	return s
}

// Helper methods

func (state EntryModel) OpenSelectedModel() tea.Model {
	if state.selectedOptionIndex == 0 {
		return CreateLoginModel(state)
	} else {
		return CreateRegisterModel(state)
	}
}