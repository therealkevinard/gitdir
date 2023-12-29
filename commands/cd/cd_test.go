package cd

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/therealkevinard/gitdir/dirtools"
	"testing"
)

var testUserPaths = dirtools.UserPaths{
	CollectionRoot: "/home/test/gitdir/repos/",
}

func Test_Command_getScriptContent(t *testing.T) {
	c := &Command{
		cdTo:  "path/to/dir",
		paths: &testUserPaths,
	}

	type test struct {
		name string
		to   string
		want string
	}

	// covers various path joins and proper slash handling
	slashTests := []test{
		{
			name: "no lead or trailing slashes",
			to:   "path/to/dir",
			want: "cd /home/test/gitdir/repos/path/to/dir",
		},
		{
			name: "with leading slash",
			to:   "/path/to/dir",
			want: "cd /home/test/gitdir/repos/path/to/dir",
		},
		{
			name: "with leading slashes",
			to:   "//path/to/dir",
			want: "cd /home/test/gitdir/repos/path/to/dir",
		},
		{
			name: "with trailing slash",
			to:   "/path/to/dir/",
			want: "cd /home/test/gitdir/repos/path/to/dir",
		},
		{
			name: "with trailing slashes",
			to:   "//path/to/dir//",
			want: "cd /home/test/gitdir/repos/path/to/dir",
		},
	}

	// covers proper handling of absolute vs relative paths
	absrelTests := []test{
		{
			name: "absolute path doesn't prefix root",
			to:   "/home/test/gitdir/repos/path/to/dir",
			want: "cd /home/test/gitdir/repos/path/to/dir",
		},
	}

	tests := make([]test, 0)
	tests = append(tests, slashTests...)
	tests = append(tests, absrelTests...)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			content, err := c.getScriptContent()
			require.Nil(t, err)
			assert.Equal(t, test.want, content)
		})
	}
}
