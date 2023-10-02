package commandtools

import "os"

func CheckRoot(r string) string {
	if r == "" {
		r = "$HOME/Workspaces"
	}

	return os.ExpandEnv(r)
}
