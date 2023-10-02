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
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commandtools"
)

const (
	name     = "cd"
	synopsis = "root-aware cd. move into a local gitdir directory"
	usage    = `
gitdir cd - 
cd to directory within your collection root.

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

func (c *Command) Execute(_ context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	var cdTo string
	if f.NArg() != 1 || f.Arg(0) != "-" {
		return subcommands.ExitUsageError
	}

	// read path from stdin
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		cdTo = scanner.Text()
	}

	// write bash using selection
	if fileErr := c.writeCDToSelection(path.Join(c.CollectionRoot, cdTo)); fileErr != nil {
		log.Printf("error creating cd script: %v", fileErr)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (c *Command) writeCDToSelection(cdTo string) error {
	if cdTo == "" {
		return commandtools.ErrInvalidDirectory
	}

	// prepare cd-path. prepend the CollectionRoot if needed
	if !strings.HasPrefix(cdTo, c.CollectionRoot) {
		cdTo = path.Clean(path.Join(c.CollectionRoot, cdTo))
	}

	// prepare write-path
	cacheDir, _ := os.UserCacheDir()
	scriptpath := path.Clean(path.Join(cacheDir, "gitdir", "gdnext.sh"))

	// create script
	//nolint:gomnd
	_ = os.MkdirAll(path.Dir(scriptpath), 0o750) // TODO: check error
	f, fileErr := os.Create(scriptpath)
	if fileErr != nil {
		return fileErr //nolint:wrapcheck
	}
	defer func() { _ = f.Close() }()

	// write file. no need to handle the no-selection case, as os.Create has truncated the file already
	_, _ = f.WriteString(fmt.Sprintf(`cd %s`, cdTo))

	return nil
}
