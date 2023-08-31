// Package cd holds the command for changing directories to a git repo.
// it indexes all repos under root, presents a bubbletea list ui, and writes a cd script for that dir that
// can be sourced to cd the outer shell to the chosen dir.
// TODO: need to formalize the alias that makes this painless. `gitdir cd && source /Users/kard/Library/Caches/gitdirnextdir` should be close, but will need refinements
package cd

import (
	"flag"
	"fmt"
	"os"
	"path"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/therealkevinard/gitdir/dirtools"
)

type Command struct {
	collectionRoot string
	selection      string
	items          []list.Item
}

func NewCommand() *Command {
	//nolint: exhaustruct
	return &Command{}
}

func (c *Command) GetName() string { return "index" }

func (c *Command) Flags() {
	flag.StringVar(&c.collectionRoot, "root", "$HOME/Workspaces", "path within home directory to root the clone tree under. supports environment expansion.")
	flag.Parse()

	c.collectionRoot = os.ExpandEnv(c.collectionRoot)
}

func (c *Command) Run() error {
	// rebuild model items from git dirs
	gitDirs, err := dirtools.FindGitDirs(c.collectionRoot)
	if err != nil {
		return fmt.Errorf("error finding git dirItems: %w", err)
	}

	c.items = make([]list.Item, len(gitDirs))
	for i, dir := range gitDirs {
		c.items[i] = stringListItem(dir)
	}

	// reset selection
	c.selection = ""

	// run the ui picker
	c.ui()

	// if user has a directory selection, write a cd script that can be sourced
	// to cd from the shell itself. binary can't directly cd since it's a subprocess, but
	// the alias
	// gitdir index && source /Users/kard/Library/Caches/gitdirnextdir
	// does the trick
	if dir := c.selection; dir != "" {
		script := fmt.Sprintf(`cd %s`, c.selection)
		cacheDir, _ := os.UserCacheDir()
		f, fileErr := os.Create(path.Join(cacheDir, "gitdirnextdir.sh"))
		if fileErr != nil {
			return fmt.Errorf("error creating cd script: %w", fileErr)
		}
		defer func() { _ = f.Close() }()

		_, _ = f.WriteString(script)
	}

	return nil
}

func (c *Command) ui() {
	const (
		defaultWidth = 20
		listHeight   = 14
	)

	l := list.New(c.items, itemDelegate{}, defaultWidth, listHeight)
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

	if c.selection != "" {
		_ = os.Chdir(c.selection)
	}
}
