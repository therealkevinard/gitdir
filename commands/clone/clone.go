package clone

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/dirtools"
)

const (
	name     = "clone"
	synopsis = "clone a remote repo url"
	usage    = "usage here"
)

type Command struct {
	collectionRoot string
	repoURL        string
	localDir       string
}

func (c *Command) Name() string     { return name }
func (c *Command) Synopsis() string { return synopsis }
func (c *Command) Usage() string    { return usage }
func (c *Command) SetFlags(set *flag.FlagSet) {
	set.StringVar(
		&c.collectionRoot,
		"root",
		"$HOME/Workspaces",
		"path within home directory to root the clone tree under. supports environment expansion.",
	)
}

func (c *Command) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// check args
	if f.Arg(0) == "" {
		log.Println("repo url must be provided as only positional argument")
		return subcommands.ExitUsageError
	}

	// TODO: flag defaults are mishandled here. this belongs in main.go, and should be passed to the command constructor
	if c.collectionRoot == "" {
		c.collectionRoot = "$HOME/Workspaces"
	}
	c.collectionRoot = os.ExpandEnv(c.collectionRoot)

	c.repoURL = f.Arg(0)

	subPath, err := dirtools.NormalizeRepoURL(c.repoURL)
	if err != nil {
		fmt.Printf("!! failed normalizing repo: %v", err)
		return subcommands.ExitFailure
	}

	// create clone directory
	c.localDir = dirtools.CompileDirPath(c.collectionRoot, subPath)
	if _, err = os.Stat(c.localDir); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("!! directory exists. not re-creating %s\n", c.localDir)
		return subcommands.ExitFailure
	}

	fmt.Printf("> creating %s\n", c.localDir)
	//nolint:gomnd
	if err = os.MkdirAll(c.localDir, 0o750); err != nil {
		fmt.Printf("!! error creating base directory: %v", err)
		return subcommands.ExitFailure
	}

	// clone operation
	fmt.Printf("> cloning %s into %s\n", c.repoURL, c.localDir)
	out, err := c.cloneRepo()
	if err != nil {
		fmt.Printf("!! error cloning. leaving empty dir at %s\n", c.localDir)
		return subcommands.ExitFailure
	}

	// status output
	fmt.Printf("> finished. git says: \n%s\n-------\n", out)

	return subcommands.ExitSuccess
}

// cloneRepo runs git clone command in localDir.
// git's stdout and stderr are captured in the first return.
// hard errors and no-zero exits are returned by err.
func (c *Command) cloneRepo() ([]byte, error) {
	var out bytes.Buffer

	cmd := exec.Command("git", "clone", c.repoURL, c.localDir)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = c.localDir

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running git clone: %w", err)
	}

	return out.Bytes(), nil
}
