package init

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commandtools"
	context_keys "github.com/therealkevinard/gitdir/context-keys"
)

const (
	name     = "init"
	synopsis = "initializes shell env"
	usage    = `
source this into your shell's .profile to prepare environment, collection root, and add cd support
`
)

type Command struct {
	shell string
}

func (c *Command) Name() string             { return name }
func (c *Command) Synopsis() string         { return synopsis }
func (c *Command) Usage() string            { return usage }
func (c *Command) SetFlags(_ *flag.FlagSet) {}

func (c *Command) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	selfCmd, ok := ctx.Value(context_keys.SelfNameCtx).(string)
	if !ok {
		err := fmt.Errorf("command name not found in command context: %w", commandtools.ErrContextMismatch)
		log.Fatal(err)
		return subcommands.ExitFailure
	}

	ue, ok := ctx.Value(context_keys.UserEnvCtx).(*commandtools.UserEnvironment)
	if !ok {
		err := fmt.Errorf("user environment not found in command context: %w", commandtools.ErrContextMismatch)
		log.Fatal(err)
		return subcommands.ExitFailure
	}

	gdcdAlias := fmt.Sprintf("%s ls | fzf | %s cd - && source %s", selfCmd, selfCmd, ue.CDShellPath())

	//nolint
	fmt.Printf(`
#!/usr/bin/env bash
# env var unset. setting
[ -z "$GITDIR_COLLECTION_ROOT" ] && GITDIR_COLLECTION_ROOT="$HOME/Workspaces"
# create root if not exist 
[ -d "$GITDIR_COLLECTION_ROOT" ]  || mkdir -p "$GITDIR_COLLECTION_ROOT"
# gitdir fzf alias 
alias gdcd="%s"
`, gdcdAlias)

	return subcommands.ExitSuccess
}
