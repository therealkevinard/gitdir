package main

import (
	"context"
	"flag"
	"github.com/therealkevinard/gitdir/commands/ls"
	"os"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commands/cd"
	"github.com/therealkevinard/gitdir/commands/clone"
)

func main() {
	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	const mgtGroup = "repo management"
	subcommands.Register(&clone.Command{}, mgtGroup)

	const navGroup = "navigation"
	subcommands.Register(&cd.Command{}, navGroup)
	subcommands.Register(&ls.Command{}, navGroup)

	flag.Parse()

	ctx := context.Background()
	code := subcommands.Execute(ctx)
	os.Exit(int(code))
}
