package dirtools

import (
	"slices"
	"strings"
)

// LocalRepoPath represents a directory path string with or without a common prefix stripped
type LocalRepoPath struct {
	full  string
	short string
}

func NewLocalRepoPath(root, full string) *LocalRepoPath {
	return &LocalRepoPath{
		full:  full,
		short: strings.TrimPrefix(full, root+"/"), // TODO: this is more rigid than i'd like
	}
}

// Short returns the dir path with common prefix trimmed
func (rp *LocalRepoPath) Short() string { return rp.short }

// Long returns the full, unmodified directory path
func (rp *LocalRepoPath) Long() string { return rp.full }

// RepoList is a string-indexed map of *LocalRepoPath
type RepoList map[string]*LocalRepoPath

func NewRepoList(root string, paths []string) RepoList {
	list := make(RepoList, len(paths))
	for i := range paths {
		list[paths[i]] = NewLocalRepoPath(root, paths[i])
	}

	return list
}

// Keys returns a sorted list of the map keys
func (rl RepoList) Keys() []string {
	idx := make([]string, 0, len(rl))
	for k := range rl {
		idx = append(idx, k)
	}
	slices.Sort[[]string](idx)

	return idx
}
