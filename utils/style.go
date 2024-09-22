// utils/style.go
package utils

import "github.com/charmbracelet/lipgloss"

var (
    HelpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
    SuccessStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#04B575"))
    ErrorStyle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF6347"))
)
