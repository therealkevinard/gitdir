package context_keys

type ContextKey string

const (
	SelfNameCtx ContextKey = "self-name"
	CollRootCtx ContextKey = "collection-root"
	UserEnvCtx  ContextKey = "user-environment"
)
