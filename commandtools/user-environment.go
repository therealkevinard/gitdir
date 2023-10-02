package commandtools

import (
	"log"
	"os"
	"path"
)

type UserEnvironment struct {
	userHomeDir    string
	userCacheDir   string
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
		gdNextFilepath: path.Join(cache, "gitdir", "gdnext.sh"),
	}
}

func (ue *UserEnvironment) Home() string        { return ue.userHomeDir }
func (ue *UserEnvironment) Cache() string       { return ue.userCacheDir }
func (ue *UserEnvironment) CDShellPath() string { return ue.gdNextFilepath }
