package main

import (
	"database/sql"
	"errors"
	"log"
	"os"

	"github.com/Romasav/gator/internal/config"
	"github.com/Romasav/gator/internal/database"

	_ "github.com/lib/pq"
)

func main() {
	con, err := config.Read()
	if err != nil {
		log.Fatalf("could not get config: %v", err.Error())
	}

	db, err := sql.Open("postgres", con.DbUrl)
	if err != nil {
		log.Fatalf("could not open a connection with db: %v", err.Error())
	}
	dbQueries := database.New(db)

	state := newState(dbQueries, con)

	commands := newCommands()
	commands.register("login", handlerLogin)
	commands.register("register", handlerRegister)
	commands.register("reset", handlerReset)
	commands.register("users", handlerUsers)
	commands.register("agg", handlerAggregator)
	commands.register("addfeed", handlerCreateFeed)
	commands.register("feeds", handlerFeeds)
	commands.register("follow", handlerFollow)
	commands.register("following", handlerFollowing)

	command, err := parseArgs()
	if err != nil {
		log.Fatalf("could not get the command: %v", err.Error())
	}

	err = commands.run(state, command)
	if err != nil {
		log.Fatalf("could not run the command: %v", err.Error())
	}
}

func parseArgs() (command, error) {
	if len(os.Args) < 2 {
		return command{}, errors.New("usage: ./app <command> [arguments...]")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]
	return *newCommand(cmdName, cmdArgs), nil
}
