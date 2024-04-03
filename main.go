package main

import (
	"fmt"
	"github.com/BentleyOph/monke/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monke programming language!\n", user.Username)
	fmt.Printf("Type in any commands to see generated tokens\n")
	repl.Start(os.Stdin, os.Stdout)
}
