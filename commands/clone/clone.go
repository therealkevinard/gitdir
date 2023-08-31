package clone

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/therealkevinard/gitdir/dirtools"
)

type Command struct {
	collectionRoot string
	repoURL        string
	localDir       string
}

func NewCommand() *Command {
	//nolint: exhaustruct
	return &Command{}
}

func (c *Command) GetName() string { return "clone" }

func (c *Command) Flags() {
	flag.StringVar(&c.collectionRoot, "root", "$HOME/Workspaces", "path within home directory to root the clone tree under. supports environment expansion.")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 || args[0] == "" {
		log.Fatal("repo url must be provided as only positional argument")
	}
	c.repoURL = args[0]
	c.collectionRoot = os.ExpandEnv(c.collectionRoot)
}

// Run wraps a complete execution, cloning repoURL under root
// it has little internal logic - mostly just composing other work with terminal output.
func (c *Command) Run() error {
	subPath, err := dirtools.NormalizeRepoURL(c.repoURL)
	if err != nil {
		log.Fatalf("!! failed normalizing repo: %v", err)
	}

	// create clone directory
	c.localDir = dirtools.CompileDirPath(c.collectionRoot, subPath)
	if _, err = os.Stat(c.localDir); !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("!! directory exists. not re-creating %s\n", c.localDir)
		log.Fatal("exiting")
	}

	fmt.Printf("> creating %s\n", c.localDir)
	//nolint:gomnd
	if err = os.MkdirAll(c.localDir, 0o750); err != nil {
		log.Fatalf("!! error creating base directory: %v", err)
	}

	// clone operation
	fmt.Printf("> cloning %s into %s\n", c.repoURL, c.localDir)
	out, err := c.cloneRepo()
	if err != nil {
		fmt.Printf("!! error cloning. leaving empty dir at %s\n", c.localDir)
		return err
	}

	// status output
	fmt.Printf("> finished. git says: \n%s\n-------\n", out)

	return err
}

// cloneRepo runs git clone command in localDir.
// git's stdout and stderr are captured in the first return.
// hard errors and no-zero exits are returned by err.
func (c *Command) cloneRepo() ([]byte, error) {
	var out bytes.Buffer

	cmd := exec.Command("git", "clone", c.repoURL, c.localDir)
	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = c.localDir

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running git clone: %w", err)
	}

	return out.Bytes(), nil
}
