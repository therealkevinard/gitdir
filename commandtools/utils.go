package commandtools

import (
	"context"
	"log"
	"os"

	context_keys "github.com/therealkevinard/gitdir/context-keys"
)

func CheckRoot(ctx context.Context) string {
	r, ok := ctx.Value(context_keys.CollRootCtx).(string)
	if !ok {
		log.Fatal("invalid value stored in collection root context key")
	}
	if r == "" {
		r = "$HOME/Workspaces"
	}

	return os.ExpandEnv(r)
}
