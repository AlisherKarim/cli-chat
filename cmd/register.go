// cmd/register.go
package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/spf13/cobra"
)

var registerCmd = &cobra.Command{
    Use:   "register",
    Short: "Register a new user",
    Run: func(cmd *cobra.Command, args []string) {
        var username, password string

        // Ask for username and password
        survey.AskOne(&survey.Input{
            Message: "Choose a username:",
        }, &username)

        survey.AskOne(&survey.Password{
            Message: "Choose a password:",
        }, &password)

        // Make API request
        response, err := utils.Register(username, password)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        // Print response
        fmt.Println("Registration successful:", response)
    },
}

func init() {
    RootCmd.AddCommand(registerCmd)
}
