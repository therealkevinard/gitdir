package open

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commands"
	"github.com/therealkevinard/gitdir/dirtools"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

const (
	name     = "open"
	synopsis = "[alpha] open repo in browser"
	usage    = `
open a project's web url. local directory is parsed into the original browser url 
acceptable arguments are
- read from stdin. 
	echo "git@github.com:therealkevinard/gitdir.git | gitdir open -" 
. open project in $PWD.
	gitdir open .
<string> open 
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

func (c *Command) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var openpath string

	switch f.Arg(0) {
	case "-":
		// read path from stdin
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			openpath = scanner.Text()
		}

	case ".":
		// use $PWD
		dir, err := os.Getwd()
		if err != nil {
			commands.Notify(commands.NotifyError, fmt.Sprintf("error getting dot-path: %v", err))
			return subcommands.ExitFailure
		}
		openpath = dir

	default:
		// use scalar value
		openpath = f.Arg(0)
	}

	if openpath == "" {
		commands.Notify(commands.NotifyError, fmt.Sprintf("no valid path provided"))
		return subcommands.ExitUsageError
	}

	if strings.HasPrefix(openpath, c.paths.CollectionRoot) {
		openpath = strings.TrimPrefix(openpath, c.paths.CollectionRoot)
		openpath = strings.TrimPrefix(openpath, "/") // TODO: this is all so gross. make a better string cleaner one day.
	}

	dest, err := url.Parse("https://" + openpath)
	if err != nil {
		commands.Notify(commands.NotifyError, fmt.Sprintf("error parsing %s as url: %v", "https://"+openpath, err))
		return subcommands.ExitFailure
	}

	commands.Notify(commands.NotifyInfo, fmt.Sprintf("opening %s", dest.String()))

	cmd := exec.Command("open", dest.String())
	if runErr := cmd.Run(); runErr != nil {
		commands.Notify(commands.NotifyError, fmt.Sprintf("error opening web url: %v", runErr))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}
