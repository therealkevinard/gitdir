package main

import (
	"fmt"
	"log"

	"github.com/therealkevinard/gitdir/commands"
	"github.com/therealkevinard/gitdir/commands/clone"
)

func main() {
	cmd := clone.NewCommand()
	exec(cmd)
}

func exec(cmd commands.Command) {
	// parse flags
	cmd.Flags()

	// run it
	fmt.Printf(">> running %s command\n", cmd.GetName())
	if err := cmd.Run(); err != nil {
		log.Fatalf("!! command %s execution exited with: %v", cmd.GetName(), err)
	}

	// all done
	fmt.Println("done!")
}
