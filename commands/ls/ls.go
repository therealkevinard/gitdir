package ls

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/dirtools"
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
	CollectionRoot string
}

func (c *Command) Name() string             { return name }
func (c *Command) Synopsis() string         { return synopsis }
func (c *Command) Usage() string            { return usage }
func (c *Command) SetFlags(_ *flag.FlagSet) {}

func (c *Command) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	dirs, err := dirtools.FindGitDirs(c.CollectionRoot)
	if err != nil {
		log.Printf("error finding git dirItems: %v", err)
		return subcommands.ExitFailure
	}

	list := dirtools.NewRepoList(c.CollectionRoot, dirs)

	for _, k := range list.Keys() {
		fmt.Println(list[k].Short())
	}

	return subcommands.ExitSuccess
}
