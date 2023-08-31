package commands

type Command interface {
	// GetName should return a friendly name for the command
	GetName() string
	// Flags configures and parses command flags
	Flags()
	// Run is the actual command execution
	Run() error
}
