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
	"github.com/therealkevinard/gitdir/commands"
	context_keys "github.com/therealkevinard/gitdir/context-keys"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commandtools"
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
	CollectionRoot string
}

func (c *Command) Name() string             { return name }
func (c *Command) Synopsis() string         { return synopsis }
func (c *Command) Usage() string            { return usage }
func (c *Command) SetFlags(_ *flag.FlagSet) {}

func (c *Command) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	userEnvironment, ok := ctx.Value(context_keys.UserEnvCtx).(*commandtools.UserEnvironment)
	if !ok {
		commands.Notify(commands.NotifyError, "no UserEnvironment found in context")
		return subcommands.ExitFailure
	}

	var cdTo string
	if f.NArg() != 1 || f.Arg(0) != "-" {
		return subcommands.ExitUsageError
	}

	// read path from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cdTo = scanner.Text()
	}

	if cdTo == "" {
		return subcommands.ExitUsageError
	}

	// write bash using selection
	if fileErr := c.writeCDToSelection(userEnvironment, path.Join(c.CollectionRoot, cdTo)); fileErr != nil {
		log.Printf("error creating cd script: %v", fileErr)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

// writeCDToSelection writes a bash script to CDShellPath as defined in UserEnvironment that calls simply `cd {{.cdTo}}`
func (c *Command) writeCDToSelection(env *commandtools.UserEnvironment, cdTo string) error {
	if cdTo == "" {
		return commandtools.ErrInvalidDirectory
	}

	// prepare cd-path. prepend the CollectionRoot if needed
	if !strings.HasPrefix(cdTo, c.CollectionRoot) {
		cdTo = path.Clean(path.Join(c.CollectionRoot, cdTo))
	}

	// create script
	//nolint:gomnd
	_ = os.MkdirAll(path.Dir(env.CDShellPath()), 0o750) // TODO: check error
	f, fileErr := os.Create(env.CDShellPath())
	if fileErr != nil {
		return fileErr //nolint:wrapcheck
	}
	defer func() { _ = f.Close() }()

	// write file. no need to handle the no-selection case, as os.Create has truncated the file already
	_, _ = f.WriteString(fmt.Sprintf(`cd %s`, cdTo))

	return nil
}
