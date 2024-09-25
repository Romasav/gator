package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Romasav/gator/internal/config"
)

func main() {
	con, err := config.Read()
	if err != nil {
		log.Fatalf("could not get config: %v", err.Error())
	}

	state := newState(&con)

	commands := newCommands()
	commands.register("login", handlerLogin)

	command, err := parseArgs()
	if err != nil {
		log.Fatalf("could not get the command: %v", err.Error())
	}

	err = commands.run(state, command)
	if err != nil {
		log.Fatalf("could not run the command: %v", err.Error())
	}

	fmt.Print(con)
}

func parseArgs() (command, error) {
	if len(os.Args) < 2 {
		return command{}, errors.New("usage: ./app <command> [arguments...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	return *newCommand(cmdName, cmdArgs), nil
}
