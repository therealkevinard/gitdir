# gitdir

a tiny little tool for organizing local git repos - immensely inspired by the old $GOROOT. that's it 🤷‍♀️

establishes a root directory that will hold git repos. within this root directory, repos are placed
in a directory that mirrors the clone url. clone url is normalized to a stable directory path. 


## install

- install the binary

```shell
go install github.com/therealkevinard/gitdir@latest 
```

- setup your env: `gitdir init` outputs some snippets that belong in your .profile file (an env var and an alias)

## subcommands

- `gitdir clone $REPO_URL`: clone into directory under collection root
- `gitdir ls`: list all repos under your collection root
- `gitdir cd -`: a root-aware cd. reads target from stdin

