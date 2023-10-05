package dirtools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NormalizeRepoURL(t *testing.T) {
	errInvalidURLText := "invalid url "

	tests := []struct {
		name        string
		input       string
		want        string
		wantErrText string
	}{
		{
			name:  "gitlab - ssh",
			input: "ssh://git@git.site.com:2223/group/project.git",
			want:  "git.site.com/group/project",
		},
		{
			name:  "gitlab - https",
			input: "https://git.site.com/user/empty.git",
			want:  "git.site.com/user/empty",
		},
		{
			name:  "github - https",
			input: "https://github.com/group/project.git",
			want:  "github.com/group/project",
		},
		{
			name:  "github - ssh",
			input: "ssh://git@github.com:group/project.git",
			want:  "github.com/group/project",
		},
		{
			name:  "github - ssh noproto",
			input: "git@github.com:group/project.git",
			want:  "github.com/group/project",
		},
		{
			name:        "invalid",
			input:       "qwertyuiop",
			want:        "",
			wantErrText: errInvalidURLText,
		},
		{
			name:        "empty",
			input:       "",
			want:        "",
			wantErrText: errInvalidURLText,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := NormalizeRepoURL(test.input)
			if test.wantErrText != "" {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), test.wantErrText)
			} else {
				assert.Nil(t, err)
			}

			assert.Equal(t, test.want, got)
		})
	}
}

func Test_CompileDirPath(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "gitlab - ssh",
			input: "git.site.com/group/project",
			want:  "/home/test/git.site.com/group/project",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := CompileDirPath("/home/test", test.input)
			assert.Equal(t, got, test.want)
		})
	}
}
