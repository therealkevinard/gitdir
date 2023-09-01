package main

import (
	"log"
	"os"

	"github.com/therealkevinard/gitdir/commands"
	"github.com/therealkevinard/gitdir/commands/cd"
	"github.com/therealkevinard/gitdir/commands/clone"
)

func main() {
	if len(os.Args) == 0 {
		return
	}

	var cmd commands.Command

	switch os.Args[1] {
	case "clone":
		cmd = clone.NewCommand()

	case "cd":
		cmd = cd.NewCommand()

	default:
		cmd = clone.NewCommand()
	}

	exec(cmd)
}

func exec(cmd commands.Command) {
	// parse flags
	cmd.Flags()

	// run it
	if err := cmd.Run(); err != nil {
		log.Fatalf("!! command %s execution exited with: %v", cmd.GetName(), err)
	}
}
