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
	ue := &UserEnvironment{}

	// home dir
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	ue.userHomeDir = home

	// cache dir
	cache, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}
	ue.userCacheDir = cache

	// cd script
	ue.gdNextFilepath = path.Join(cache, "gitdir", "gdnext.sh")

	return ue
}

func (ue *UserEnvironment) Home() string        { return ue.userHomeDir }
func (ue *UserEnvironment) Cache() string       { return ue.userCacheDir }
func (ue *UserEnvironment) CDShellPath() string { return ue.gdNextFilepath }
