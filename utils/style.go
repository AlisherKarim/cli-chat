// utils/style.go
package utils

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

var (
    UserInfoBoxStyle = lipgloss.NewStyle().
			Width(120).
			Height(4).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("205"))

    PageNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("204")).Background(lipgloss.Color("235"))
    HelpStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("241")).
        Width(120).
        Height(1).
        Align(lipgloss.Left, lipgloss.Top).
        BorderStyle(lipgloss.NormalBorder()).
        BorderForeground(lipgloss.Color("241"))
    SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
    ErrorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6347"))
    FocusedStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	BlurredStyle        = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	CursorStyle         = FocusedStyle
	NoStyle             = lipgloss.NewStyle()
	CursorModeHelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

	FocusedButton = FocusedStyle.Render("[ Submit ]")
	BlurredButton = fmt.Sprintf("[ %s ]", BlurredStyle.Render("Submit"))

    ChatsListTabStyle = lipgloss.NewStyle().
			Width(20).
			Height(5).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.HiddenBorder())
    ChatTabStyle = lipgloss.NewStyle().
			Width(98).
			Height(30).
			Align(lipgloss.Left, lipgloss.Top).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("69"))

    LoadingSpinner = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
)
