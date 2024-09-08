package main

import (
	"flag"
	"fmt"
)

func main() {
	var (
		room   string
	)

	flag.StringVar(&room, "room", "c123456", "Room code")
	flag.Parse()

	fmt.Printf("Joining room %s...\n", room)
	fmt.Printf("You have joined room #%s successfully\n", room)
}