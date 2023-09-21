// Package cd holds the command for changing directories to a git repo.
// it indexes all repos under root, presents a bubbletea list ui, and writes a cd script for that dir that
// can be sourced to cd the outer shell to the chosen dir.
//
// since binaries run as a subprocess, this needs an alias like `xxx='gitdir cd && source script.sh'
// where, on selection, `cd /chosen/path` is written to script.sh by the binary
package cd

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/dirtools"
)

const (
	name     = "cd"
	synopsis = "cd into a local gitdir directory"
	usage    = "usage here"
)

type Command struct {
	collectionRoot string
	selection      string
	items          []list.Item
}

func (c *Command) Name() string     { return name }
func (c *Command) Synopsis() string { return synopsis }
func (c *Command) Usage() string    { return usage }
func (c *Command) SetFlags(set *flag.FlagSet) {
	set.StringVar(&c.collectionRoot, "root", "$HOME/Workspaces", "path within home directory to root the clone tree under. supports environment expansion.")
}

func (c *Command) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	if c.collectionRoot == "" {
		c.collectionRoot = "$HOME/Workspaces"
	}
	c.collectionRoot = os.ExpandEnv(c.collectionRoot)

	// build model items from git dirs
	gitDirs, err := dirtools.FindGitDirs(c.collectionRoot)
	if err != nil {
		log.Printf("error finding git dirItems: %v", err)
		return subcommands.ExitFailure
	}

	c.items = make([]list.Item, len(gitDirs))
	for i, dir := range gitDirs {
		c.items[i] = stringListItem(dir)
	}

	// clear selection
	c.selection = ""

	// run the ui picker
	c.ui()

	// write bash using selection
	if fileErr := c.writeCDToSelection(); fileErr != nil {
		log.Printf("error creating cd script: %v", fileErr)
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (c *Command) ui() {
	const (
		defaultWidth = 20
		listHeight   = 14
	)

	l := list.New(c.items, itemDelegate{stripPrefix: c.collectionRoot}, defaultWidth, listHeight)
	l.Title = "found these repositories in your collection"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)

	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	m := &model{command: c, list: l, choice: "", quitting: false}

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (c *Command) writeCDToSelection() error {
	// prepare paths
	cacheDir, _ := os.UserCacheDir()
	scriptpath := path.Clean(path.Join(cacheDir, "gitdir", "gdnext.sh"))

	// create script
	_ = os.MkdirAll(path.Dir(scriptpath), 0o750) // TODO: check error
	f, fileErr := os.Create(scriptpath)
	if fileErr != nil {
		return fileErr //nolint:wrapcheck
	}
	defer func() { _ = f.Close() }()

	// write file. no need to handle the no-selection case, as os.Create has truncated the file already
	if c.selection != "" {
		_, _ = f.WriteString(fmt.Sprintf(`cd %s`, c.selection))
	}

	return nil
}
