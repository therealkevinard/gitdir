package commandtools

import (
	"log"
	"os"
	"path"
)

// UserEnvironment stores key global info from the user's environment
type UserEnvironment struct {
	// user os' home dir
	userHomeDir string
	// user os' cache dir
	userCacheDir string
	// path to cd shell script. this sh is updated as part of the cd subcommand
	gdNextFilepath string
}

func InitUserEnvironment() *UserEnvironment {
	// home dir
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// cache dir
	cache, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	return &UserEnvironment{
		userHomeDir:    home,
		userCacheDir:   cache,
		gdNextFilepath: path.Clean(path.Join(cache, "gitdir", "gdnext.sh")),
	}
}

func (ue *UserEnvironment) Home() string        { return ue.userHomeDir }
func (ue *UserEnvironment) Cache() string       { return ue.userCacheDir }
func (ue *UserEnvironment) CDShellPath() string { return ue.gdNextFilepath }
