package open

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commands"
	"os"
	"os/exec"
)

const (
	name     = "open"
	synopsis = "open repo in browser"
	usage    = `
open a project's web url
`
)

type Command struct {
	CollectionRoot string
}

func (c *Command) Name() string             { return name }
func (c *Command) Synopsis() string         { return synopsis }
func (c *Command) Usage() string            { return usage }
func (c *Command) SetFlags(_ *flag.FlagSet) {}

func (c *Command) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var openpath string

	if p := f.Arg(0); p == "-" {
		// read path from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			openpath = scanner.Text()
		}
	} else {
		openpath = p
	}

	if openpath == "" {
		return subcommands.ExitUsageError
	}

	cmd := exec.Command("open", "https://"+openpath)
	if runErr := cmd.Run(); runErr != nil {
		commands.Notify(commands.NotifyError, fmt.Sprintf("error opening web url: %v", runErr))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
