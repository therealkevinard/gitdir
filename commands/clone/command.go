package clone

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/therealkevinard/gitdir/dirtools"
	"github.com/therealkevinard/gitdir/ui/styles"
	"log"
	"os"
	"os/exec"
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

	// we have two cmds that lead here: `clone xxx` and `xxx`
	// for the `clone xxx` case, strip clone
	args := flag.Args()
	if args[0] == "clone" {
		args = args[1:]
	}
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
		styles.Println(styles.ErrorLevel, "failed normalizing repo: %v", err)
		return err
	}

	// create clone directory
	c.localDir = dirtools.CompileDirPath(c.collectionRoot, subPath)
	if _, err = os.Stat(c.localDir); !errors.Is(err, os.ErrNotExist) {
		err = os.ErrExist
		styles.Println(styles.ErrorLevel, "not-recreating (%s): %v", c.localDir, err)
		return err
	}

	styles.Println(styles.OKLevel, "creating %s", c.localDir)
	//nolint:gomnd
	if err = os.MkdirAll(c.localDir, 0o750); err != nil {
		styles.Println(styles.ErrorLevel, "error creating base directory: %v", err)
		return err
	}

	// clone operation
	styles.Println(styles.OKLevel, "cloning %s into %s", c.repoURL, c.localDir)
	out, err := c.cloneRepo()
	if err != nil {
		styles.Println(styles.ErrorLevel, "error cloning. leaving empty dir at %s", c.localDir)
		return err
	}

	// status output
	parts := []string{
		styles.OKTextf("finished. git says:"),
		styles.AltTextStyle.PaddingLeft(4).Render(string(out)),
		styles.OKTextf("done here."),
	}
	for _, v := range parts {
		fmt.Println(v)
	}

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
