package chatting

import (
	"fmt"

	"github.com/chzyer/readline"
)

func StartChat(username string) {
	rl, err := readline.New("> ")
	if err != nil {
					panic(err)
	}
	defer rl.Close()

	messages := []string{}

	for {
		line, err := rl.Readline()
		if err != nil {
						break
		}

		if line == "/quit" {
			break
		}

		messages = append(messages, line)

		// Move the cursor to the top of the terminal
		fmt.Print("\033[1;1H")
		// Clear the screen
		fmt.Print("\033[2J")

		for i := 0; i < len(messages); i++ {
			// fmt.Println(messages[i])
			fmt.Printf("\033[32m%s:\033[0m %s\n", username, messages[i])
		}

		// Print the user input prompt at the bottom
		fmt.Print("\033[99;1H")
		fmt.Print("> ")
		rl.SetPrompt("> ")
	}

	
}