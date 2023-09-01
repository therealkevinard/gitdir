package main

import (
	"fmt"
	"github.com/therealkevinard/gitdir/ui/styles"
	"log"
	"os"

	"github.com/therealkevinard/gitdir/commands"
	"github.com/therealkevinard/gitdir/commands/cd"
	"github.com/therealkevinard/gitdir/commands/clone"
)

func main() {
	if len(os.Args) == 1 {
		// TODO: maybe default to `cd` command, which can work with no arguments?
		log.Fatal("not enough arguments")
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
	styles.Println(styles.OKLevel, "[%s] parsed flags", cmd.GetName())

	// run it
	if err := cmd.Run(); err != nil {
		styles.Println(styles.FatalLevel, "[%s] execution exited with error: %v", cmd.GetName(), err)
		fmt.Println("...")
	}

	styles.Println(styles.OKLevel, "[%s] finished", cmd.GetName())
}
