package clone

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
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
	runningCmd     *exec.Cmd
	ui             *model
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

type statusTextUpdate string

// Run wraps a complete execution, cloning repoURL under root
// it has little internal logic - mostly just composing other work with terminal output.
func (c *Command) Run() error {
	c.ui = &model{}
	c.ui.spinner = spinner.New()
	c.ui.spinner.Spinner = spinner.MiniDot

	// run ui in a separate goroutine
	go func() {
		if _, err := tea.NewProgram(c.ui).Run(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}
	}()

	c.ui.Update(statusTextUpdate("normalizing repo url"))
	subPath, err := dirtools.NormalizeRepoURL(c.repoURL)
	if err != nil {
		styles.Println(styles.ErrorLevel, "failed normalizing repo: %v", err)
		return err
	}

	// create clone directory
	c.ui.Update(statusTextUpdate("checking directories"))
	c.localDir = dirtools.CompileDirPath(c.collectionRoot, subPath)
	if _, err = os.Stat(c.localDir); !errors.Is(err, os.ErrNotExist) {
		err = os.ErrExist
		styles.Println(styles.ErrorLevel, "not-recreating (%s): %v", c.localDir, err)
		return err
	}

	//nolint:gomnd
	if err = os.MkdirAll(c.localDir, 0o750); err != nil {
		styles.Println(styles.ErrorLevel, "error creating base directory: %v", err)
		return err
	}

	// clone operation
	c.ui.Update(statusTextUpdate(fmt.Sprintf("cloning %s", c.repoURL)))
	out, err := c.cloneRepo()
	if err != nil {
		styles.Println(styles.ErrorLevel, "error cloning. leaving empty dir at %s", c.localDir)
		return err
	}

	c.ui.Update(statusTextUpdate("finished"))
	// status output
	styles.Println(styles.OKLevel, "finished. cloned %s into %s", c.repoURL, c.localDir)
	styles.Println(styles.OKLevel, "git says:")
	fmt.Println(styles.AltTextStyle.PaddingLeft(4).Render(string(out)))
	styles.Println(styles.OKLevel, "done here.")

	return err
}

func (c *Command) Stop() {
	if c.runningCmd != nil {
		_ = c.runningCmd.Process.Kill()
	}
	if c.ui != nil {
		c.ui.Update(statusTextUpdate("finished")) // magic string to shutdown bubble app. make a proper msg one day.
	}
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

	c.runningCmd = cmd
	defer func() { c.runningCmd = nil }()
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running git clone: %w", err)
	}

	ret := out.Bytes()[:64]
	ret = append(ret, []byte("...")...)
	return ret, nil
}
