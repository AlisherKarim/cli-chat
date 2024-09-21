package main

import (
	"log"

	teaModels "github.com/alisherkarim/cli-chat/tea-models"
	tea "github.com/charmbracelet/bubbletea"
)

func createState() teaModels.AppState {
	return teaModels.AppState{}
}

func main() {
	p := tea.NewProgram(teaModels.CreateEntryModel(), tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}