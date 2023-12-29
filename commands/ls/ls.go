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
	synopsis = "list all repos under your collection root"
	usage    = `
gitdir ls 
provides a list of local repositories housed under your collection root  
works from anywhere in your filesystem

walks the directory tree starting at collection root, listing-out all git repos. 
plays well with fzf.
`
)

type Command struct {
	optionFullPath bool
	paths          *dirtools.UserPaths
}

func New(ctx context.Context) *Command {
	return &Command{paths: dirtools.GetUserPaths(ctx)}
}

func (c *Command) Name() string     { return name }
func (c *Command) Synopsis() string { return synopsis }
func (c *Command) Usage() string    { return usage }
func (c *Command) SetFlags(f *flag.FlagSet) {
	f.BoolVar(&c.optionFullPath, "dirs", false, `report absolute path, not repo name. this works well with scripted cd.`)
}

func (c *Command) Execute(_ context.Context, _ *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	dirs, err := dirtools.FindGitDirs(c.paths.CollectionRoot)
	if err != nil {
		log.Printf("error finding git dirItems: %v", err)
		return subcommands.ExitFailure
	}

	list := dirtools.NewRepoList(c.paths.CollectionRoot, dirs)

	for _, k := range list.Keys() {
		fmt.Println(list[k].Path(c.optionFullPath))
	}

	return subcommands.ExitSuccess
}
