package main

import (
	"fmt"
	"os"

	"github.com/SantiagoZuluaga/fileserver/client"
	"github.com/SantiagoZuluaga/fileserver/server"
)

func main() {

	args := os.Args[1:]

	switch args[0] {
	case "client":
		client.RunTCPClient()
	case "server":
		server.RunTCPServer()
	default:
		fmt.Println("Invalid command \nCommands available:\nclient\nserver.")
	}

}
