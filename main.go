package main

import (
	"context"
	"flag"
	"os"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commands/clone"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(&clone.Command{}, "clone")

	flag.Parse()
	ctx := context.Background()
	code := subcommands.Execute(ctx)
	os.Exit(int(code))
}
