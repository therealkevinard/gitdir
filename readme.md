# gitdir

a tiny little tool for organizing local git repos - immensely inspired by the old $GOROOT. that's it 🤷‍♀️

establishes a root directory that will hold git repos. within this root directory, repos are placed
in a directory that mirrors the clone url. clone url is normalized to a stable directory path.

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

- `gitdir clone $REPO_URL`: clone into directory under collection root
- `gitdir ls`: list all repos under your collection root
- `gitdir cd -`: a root-aware cd. reads target from stdin

