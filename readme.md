# gitdir

a tiny little tool for organizing local git repos - immensely inspired by the old $GOROOT. that's it 🤷‍♀️

establishes a root directory that will hold git repos. within this root directory, repos are placed
in directories that mirror the clone urls. clone url is normalized to a stable directory path and the cloned repo is placed in a deterministic location

## install 

- install the binary

```shell
go install github.com/therealkevinard/gitdir@latest 
```

- setup profile: add `source <(gitdir init)` to your .profile.  
  this script sets the `$GITDIR_COLLECTION_ROOT` to default $HOME/Workspaces if it's unset, and creates the directory if
  it doesn't exist.   
  it also creates the fzf alias that makes `gitdir cd` pleasant. 

## subcommands

Usage: gitdir <flags> <subcommand> <subcommand args>

Subcommands for navigation:
	cd               root-aware cd. move into a local gitdir directory
	ls               list local repositories
	open             open repo in browser

Subcommands for repo management:
	clone            clone a remote repo url

Subcommands for support:
	commands         list all command names
	flags            describe all known top-level flags
	help             describe subcommands and their syntax
	init             initializes shell env

