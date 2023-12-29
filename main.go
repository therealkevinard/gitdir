package main

import (
	"context"
	"flag"
	"github.com/therealkevinard/gitdir/commands/open"
	"github.com/therealkevinard/gitdir/dirtools"
	"log"
	"os"
	"path"

	"github.com/google/subcommands"
	"github.com/therealkevinard/gitdir/commands/cd"
	"github.com/therealkevinard/gitdir/commands/clone"
	initCmd "github.com/therealkevinard/gitdir/commands/init"
	"github.com/therealkevinard/gitdir/commands/ls"
)

func main() {
	p := prepareUserPaths()                               // setup common paths
	ctx := dirtools.SetUserPaths(context.Background(), p) // add paths to context

	const supportGroup = "support"
	subcommands.Register(subcommands.HelpCommand(), supportGroup)
	subcommands.Register(subcommands.FlagsCommand(), supportGroup)
	subcommands.Register(subcommands.CommandsCommand(), supportGroup)
	subcommands.Register(initCmd.New(ctx), supportGroup)

	const mgtGroup = "repo management"
	subcommands.Register(clone.New(ctx), mgtGroup)

	const navGroup = "navigation"
	subcommands.Register(open.New(ctx), navGroup)
	subcommands.Register(cd.New(ctx), navGroup)
	subcommands.Register(ls.New(ctx), navGroup)

	flag.Parse()

	code := subcommands.Execute(ctx)
	os.Exit(int(code))
}

// prepareUserPaths builds a dirtools.UserPaths for cli environment
func prepareUserPaths() *dirtools.UserPaths {
	var err error

	paths := &dirtools.UserPaths{}

	if paths.UserHomeDir, err = os.UserHomeDir(); err != nil {
		log.Fatal(err)
	}
	if paths.UserCacheDir, err = os.UserCacheDir(); err != nil {
		log.Fatal(err)
	}
	if paths.CollectionRoot = os.ExpandEnv(os.Getenv("GITDIR_COLLECTION_ROOT")); paths.CollectionRoot == "" {
		log.Fatal("invalid collection root")
	}
	if paths.OwnBinaryName = os.Args[0]; paths.OwnBinaryName == "" {
		log.Fatal("i have no name")
	}

	paths.CDScriptPath = path.Clean(path.Join(paths.UserCacheDir, "gitdir", "gdnext.sh"))

	return paths
}
