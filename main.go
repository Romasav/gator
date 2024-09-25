package main

import (
	"fmt"
	"log"

	"github.com/Romasav/gator/internal/config"
)

func main() {
	con, err := config.Read()
	if err != nil {
		log.Fatalf("could not get config")
	}

	con.SetUpUser("kavuunnn")

	con, err = config.Read()
	if err != nil {
		log.Fatalf("could not get config")
	}

	fmt.Print(con)
}
