package init

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/dirtools"
)

const (
	name     = "init"
	synopsis = "initializes shell env"
	usage    = `
source this into your shell's .profile to prepare environment, collection root, and add cd support
`
)

type Command struct {
	paths *dirtools.UserPaths
}

func New(ctx context.Context) *Command {
	return &Command{paths: dirtools.GetUserPaths(ctx)}
}

func (c *Command) Name() string             { return name }
func (c *Command) Synopsis() string         { return synopsis }
func (c *Command) Usage() string            { return usage }
func (c *Command) SetFlags(_ *flag.FlagSet) {}

func (c *Command) Execute(ctx context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	gdcdAlias := fmt.Sprintf("%s ls | fzf | %s cd - && source %s", c.paths.OwnBinaryName, c.paths.OwnBinaryName, c.paths.CDScriptPath)
	gdopenAlias := fmt.Sprintf("%s ls | fzf | %s open -", c.paths.OwnBinaryName, c.paths.OwnBinaryName)

	//nolint
	fmt.Printf(`
#!/usr/bin/env bash
# env var unset. setting
[ -z "$GITDIR_COLLECTION_ROOT" ] && GITDIR_COLLECTION_ROOT="$HOME/Workspaces"
# create root if not exist 
[ -d "$GITDIR_COLLECTION_ROOT" ]  || mkdir -p "$GITDIR_COLLECTION_ROOT"
# gitdir fzf aliases 
alias gdcd="%s"
alias gdopen="%s"
`, gdcdAlias, gdopenAlias)

	return subcommands.ExitSuccess
}
