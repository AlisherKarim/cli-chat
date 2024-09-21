// cmd/login.go
package cmd

import (
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/alisherkarim/cli-chat/utils"
	"github.com/spf13/cobra"
)

type LoginResponse struct {
	Username string `json:"username"`
}

var loginCmd = &cobra.Command{
    Use:   "login",
    Short: "Login to the system",
    Run: func(cmd *cobra.Command, args []string) {
        var username, password string

        // Ask for username and password
        survey.AskOne(&survey.Input{
            Message: "Enter your username:",
        }, &username)

        survey.AskOne(&survey.Password{
            Message: "Enter your password:",
        }, &password)

        // Make API request
        response, err := utils.Login(username, password)
        if err != nil {
            fmt.Println("Error:", err)
            return
        }

        // Print response
        fmt.Println("Login successful:", response)
    },
}

func init() {
    RootCmd.AddCommand(loginCmd)
}
