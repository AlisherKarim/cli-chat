package main

import (
	"log"

	"github.com/alisherkarim/cli-chat/env"
	teaModels "github.com/alisherkarim/cli-chat/tea-models"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	env := env.Env{
		// CurrentProgram: p,
	}
	entryModel := teaModels.CreateEntryModel(&env)
	p := tea.NewProgram(entryModel, tea.WithAltScreen())
	env.CurrentProgram = p

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}