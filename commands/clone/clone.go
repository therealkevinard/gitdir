package clone

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commands"
	"github.com/therealkevinard/gitdir/dirtools"
)

const (
	name     = "clone"
	synopsis = "clone a remote repo url"
	usage    = `
gitdir clone $REPO_URL 
clones $REPO_URL into a directory that mirrors the repo url. 

ssh urls, http auth, and many other nuances are normalized to a stable path within your collection root
`
)

type Command struct {
	repoURL  string
	localDir string
	paths    *dirtools.UserPaths
}

func New(ctx context.Context) *Command {
	return &Command{paths: dirtools.GetUserPaths(ctx)}
}

func (c *Command) Name() string             { return name }
func (c *Command) Synopsis() string         { return synopsis }
func (c *Command) Usage() string            { return usage }
func (c *Command) SetFlags(_ *flag.FlagSet) {}

func (c *Command) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// check args
	if f.Arg(0) == "" {
		commands.Notify(commands.NotifyError, "repo url must be provided as only positional argument")
		return subcommands.ExitUsageError
	}
	c.repoURL = f.Arg(0)

	subPath, err := dirtools.NormalizeRepoURL(c.repoURL)
	if err != nil {
		commands.Notify(commands.NotifyError, fmt.Sprintf("failed normalizing repo: %v", err))
		return subcommands.ExitFailure
	}

	// create clone directory
	c.localDir = dirtools.CompileDirPath(c.paths.CollectionRoot, subPath)
	if _, err = os.Stat(c.localDir); !errors.Is(err, os.ErrNotExist) {
		commands.Notify(commands.NotifyError, fmt.Sprintf("directory exists. not re-creating %s", c.localDir))
		return subcommands.ExitFailure
	}

	commands.Notify(commands.NotifyCreate, fmt.Sprintf("creating %s", c.localDir))

	//nolint:gomnd
	if err = os.MkdirAll(c.localDir, 0o750); err != nil {
		commands.Notify(commands.NotifyError, fmt.Sprintf("error creating base directory: %v", err))
		return subcommands.ExitFailure
	}

	// clone operation
	commands.Notify(commands.NotifyClone, fmt.Sprintf("cloning %s into %s", c.repoURL, c.localDir))
	out, err := c.cloneRepo()
	if err != nil {
		commands.Notify(commands.NotifyError, fmt.Sprintf("error cloning. leaving empty dir at %s", c.localDir))
		return subcommands.ExitFailure
	}

	// status output
	commands.Notify(commands.NotifyDone, fmt.Sprintf("finished. git says: \n%s", out))

	return subcommands.ExitSuccess
}

// cloneRepo runs git clone command in localDir.
// git's stdout and stderr are captured in the first return.
// hard errors and no-zero exits are returned by err.
func (c *Command) cloneRepo() ([]byte, error) {
	var out bytes.Buffer

	//nolint:gosec
	cmd := exec.Command("git", "clone", c.repoURL, c.localDir)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = c.localDir

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running git clone: %w", err)
	}

	return out.Bytes(), nil
}
