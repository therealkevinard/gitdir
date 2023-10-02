package main

import (
	"context"
	"flag"
	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commands/cd"
	"github.com/therealkevinard/gitdir/commands/clone"
	initCmd "github.com/therealkevinard/gitdir/commands/init"
	"github.com/therealkevinard/gitdir/commands/ls"
	"github.com/therealkevinard/gitdir/commandtools"
	context_keys "github.com/therealkevinard/gitdir/context-keys"
	"os"
	"path"
)

func main() {
	ctx := prepareCommandContext()

	const supportGroup = "support"
	subcommands.Register(subcommands.HelpCommand(), supportGroup)
	subcommands.Register(subcommands.FlagsCommand(), supportGroup)
	subcommands.Register(subcommands.CommandsCommand(), supportGroup)
	subcommands.Register(&initCmd.Command{}, supportGroup)

	const mgtGroup = "repo management"
	subcommands.Register(&clone.Command{}, mgtGroup)

	const navGroup = "navigation"
	subcommands.Register(&cd.Command{}, navGroup)
	subcommands.Register(&ls.Command{}, navGroup)

	flag.Parse()

	code := subcommands.Execute(ctx)
	os.Exit(int(code))
}

// prepareCommandContext initializes a context.Context for commands to run under.
// the prepared context will hold global runtime keys.
func prepareCommandContext() context.Context {
	ctx := context.Background()

	// command name, as installed/called
	_, selfCmd := path.Split(os.Args[0])
	ctx = context.WithValue(ctx, context_keys.SelfNameCtx, selfCmd)

	// collection root from env var
	collRoot := os.Getenv("GITDIR_COLLECTION_ROOT")
	ctx = context.WithValue(ctx, context_keys.CollRootCtx, collRoot)

	//
	userEnv := commandtools.InitUserEnvironment()
	ctx = context.WithValue(ctx, context_keys.UserEnvCtx, userEnv)

	return ctx
}
