package dirtools

import "context"

const userPathsContextKey = "user-paths"

func SetUserPaths(ctx context.Context, paths *UserPaths) context.Context {
	return context.WithValue(ctx, userPathsContextKey, paths)
}

func GetUserPaths(ctx context.Context) *UserPaths {
	// skipping the usual type-checking for now. just return the nil if it's !ok
	return ctx.Value(userPathsContextKey).(*UserPaths)
}

// UserPaths stores root-level paths used by the application
type UserPaths struct {
	// OwnBinaryName is the name of this binary, as installed/executed
	OwnBinaryName string
	// user os' home dir
	UserHomeDir string
	// user os' cache dir
	UserCacheDir string
	// path to cd shell script. this sh is updated as part of the cd subcommand
	CDScriptPath string
	// CollectionRoot is the root path for all git dirs
	CollectionRoot string
}
