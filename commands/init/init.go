package init

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commandtools"
	context_keys "github.com/therealkevinard/gitdir/context-keys"
)

const (
	name     = "init"
	synopsis = "initialize global application config"
	usage    = `
TODO... 
`
)

type Command struct {
	shell string
}

func (c *Command) Name() string     { return name }
func (c *Command) Synopsis() string { return synopsis }
func (c *Command) Usage() string    { return usage }
func (c *Command) SetFlags(set *flag.FlagSet) {
	set.StringVar(&c.shell, "shell", "", "path within home directory to root the clone tree under. supports environment expansion.")
}

func (c *Command) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	selfCmd := ctx.Value(context_keys.SelfNameCtx)
	ue := ctx.Value(context_keys.UserEnvCtx).(*commandtools.UserEnvironment)

	gdcdAlias := fmt.Sprintf("%s ls | fzf | %s cd - && source %s", selfCmd, selfCmd, ue.CDShellPath())

	fmt.Printf(`
# add this alias to your .profile
alias gdcd="%s"
`, gdcdAlias)

	return subcommands.ExitSuccess
}
