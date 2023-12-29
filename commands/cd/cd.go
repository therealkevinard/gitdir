// Package cd holds the command for changing directories to a git repo.
// it reads a target path from stdin and writes a cd script for that dir that
// can be sourced to cd the outer shell to the chosen dir.
//
// since binaries run as a subprocess, this needs an alias like `xxx='gitdir cd && source script.sh'
// where, on selection, `cd /chosen/path` is written to script.sh by the binary
package cd

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/therealkevinard/gitdir/dirtools"
	"github.com/therealkevinard/gitdir/errors"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/subcommands"
)

const (
	name     = "cd"
	synopsis = "root-aware cd. move from anywhere to a local gitdir directory"
	usage    = `
gitdir cd -  
cd to directory within your collection root.
to support this command, add source <(gitdir init) to your bash profile/rc file. 

reads target directory from stdin, prefixes your root, and writes a cd script to ~/Caches/gitdir/gdnext.sh 
it's important to source this script afterward to exec the actual cd. 
`
)

type Command struct {
	cdTo          string
	paths         *dirtools.UserPaths
	writeExecFile func(string, string) error // abstracted for test mocks
}

func New(ctx context.Context) *Command {
	return &Command{
		paths:         dirtools.GetUserPaths(ctx),
		writeExecFile: dirtools.WriteExecFile,
	}
}

func (c *Command) Name() string     { return name }
func (c *Command) Synopsis() string { return synopsis }
func (c *Command) Usage() string    { return usage }
func (c *Command) SetFlags(f *flag.FlagSet) {
	f.StringVar(&c.cdTo, "to", "", "target path")
}

func (c *Command) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	// read path from stdin if we have the conventional - argument
	if f.Arg(0) == "-" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			c.cdTo = scanner.Text()
		}
	}
	if c.cdTo == "" {
		return subcommands.ExitUsageError
	}

	// write bash using selection
	if fileErr := c.createAndWriteScript(); fileErr != nil {
		log.Printf("error creating cd script: %v", fileErr)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

// createAndWriteScript passes the output of getScriptContent through writeExecFile
func (c *Command) createAndWriteScript() error {
	content, err := c.getScriptContent()
	if err != nil {
		return err
	}

	if err := c.writeExecFile(c.paths.CDScriptPath, content); err != nil {
		return fmt.Errorf("error writing script content: %w", err)
	}

	return nil
}

// getScriptContent creates the cd script based off of c.cdTo
func (c *Command) getScriptContent() (string, error) {
	if c.cdTo == "" {
		return "", errors.ErrInvalidDirectory
	}

	// make absolute by prepending collectionroot
	if !strings.HasPrefix(c.cdTo, c.paths.CollectionRoot) {
		c.cdTo = path.Join(c.paths.CollectionRoot, c.cdTo)
	}

	return fmt.Sprintf(`cd %s`, path.Clean(c.cdTo)), nil
}
