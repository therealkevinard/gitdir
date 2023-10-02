package ls

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commandtools"
	"github.com/therealkevinard/gitdir/dirtools"
	"log"
)

const (
	name     = "ls"
	synopsis = "list local repositories"
	usage    = `
gitdir ls 
provides a list of local repositories housed under your collection root 

walks the directory tree starting at collection root, listing-out all git repos. 
plays well with fzf.
`
)

type Command struct {
	collectionRoot string
}

func (c *Command) Name() string     { return name }
func (c *Command) Synopsis() string { return synopsis }
func (c *Command) Usage() string    { return usage }
func (c *Command) SetFlags(set *flag.FlagSet) {
	set.StringVar(&c.collectionRoot, "root", "$HOME/Workspaces", "path within home directory to root the clone tree under. supports environment expansion.")
}

func (c *Command) Execute(ctx context.Context, f *flag.FlagSet, args ...interface{}) subcommands.ExitStatus {
	c.collectionRoot = commandtools.CheckRoot(c.collectionRoot)

	dirs, err := dirtools.FindGitDirs(c.collectionRoot)
	if err != nil {
		log.Printf("error finding git dirItems: %v", err)
		return subcommands.ExitFailure
	}

	list := dirtools.NewRepoList(c.collectionRoot, dirs)

	for _, k := range list.Keys() {
		fmt.Println(list[k].Short())
	}

	return subcommands.ExitSuccess
}
