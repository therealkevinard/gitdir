package dirtools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_RepoList_Keys(t *testing.T) {
	root := "/home/user/repos"
	paths := []string{
		"/home/user/repos/ccc",
		"/home/user/repos/ddd",
		"/home/user/repos/aaa",
		"/home/user/repos/bbb",
	}

	rl := NewRepoList(root, paths)
	keys := rl.Keys()

	// keys should be sorted
	assert.Equal(t, "/home/user/repos/aaa", keys[0])
	assert.Equal(t, "/home/user/repos/ddd", keys[3])
	// all should be present
	assert.Equal(t, 4, len(keys))
}

func Test_LocalRepoPath(t *testing.T) {
	root := "/home/user/repos"
	path := "/home/user/repos/aaa"

	lrp := NewLocalRepoPath(root, path)

	// assert full/long path
	assert.Equal(t, path, lrp.Path(true))
	assert.Equal(t, path, lrp.Long())

	// assert short path can recompile to full path
	assert.Equal(t, path, root+"/"+lrp.Short())
}

func Fuzz_LocalRepoPath(f *testing.F) {
	root := "/home/user/repos"

	f.Fuzz(func(t *testing.T, path string) {
		// prepend root to fuzz path
		asAbs := root + "/" + path
		lrp := NewLocalRepoPath(root, asAbs)
		// lrp.Short should eq fuzz path
		if lrp.Short() != path {
			t.Errorf("fail: %s", path)
		}
	})
}
