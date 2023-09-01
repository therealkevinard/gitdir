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
	"syscall"
)

type Command struct {
	// the tea app intercepts os signals. as a workaround, we allow that, but in
	// the tea handler, we send a signal back to the channel from main.go. this lets
	// tea act as a signal proxy so everything can close-out
	signalsChan chan os.Signal
	// the global root dir of all collected git repos
	collectionRoot string
	// remote url to clone
	repoURL string
	// local directory to clone into
	localDir string
	// when the git command is running, its ref is kept here to support os signals
	runningCmd *exec.Cmd
	// a simple tea app with spinners
	ui *model
}

// NewCommand creates a new clone command. signalsChan is used to stop the outer application
func NewCommand(signalsChan chan os.Signal) *Command {
	//nolint: exhaustruct
	return &Command{
		signalsChan: signalsChan,
	}
}

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

		c.signalsChan <- syscall.SIGINT
	}()

	c.notifyUser("normalizing repo url")

	subPath, err := dirtools.NormalizeRepoURL(c.repoURL)
	if err != nil {
		styles.Println(styles.ErrorLevel, "failed normalizing repo: %v", err)
		return err
	}

	// create clone directory
	c.notifyUser("checking directories")
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
	c.notifyUser(fmt.Sprintf("cloning %s", c.repoURL))
	out, err := c.cloneRepo()
	if err != nil {
		styles.Println(styles.ErrorLevel, "error cloning. leaving empty dir at %s", c.localDir)
		return err
	}

	// status output
	c.notifyUser("finished")
	styles.Println(styles.OKLevel, "finished. cloned %s into %s", c.repoURL, c.localDir)
	styles.Println(styles.OKLevel, "git says:")
	fmt.Println(styles.AltTextStyle.PaddingLeft(4).Render(string(out)))
	styles.Println(styles.OKLevel, "done here.")

	return err
}

func (c *Command) Stop() {
	if c.runningCmd != nil {
		err := c.runningCmd.Process.Kill()
		c.runningCmd = nil
		fmt.Println(err)
	}
	if c.ui != nil {
		c.ui.Update(statusTextUpdate("finished")) // magic string to shutdown bubble app. make a proper msg one day.
	}
}

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

func (c *Command) GetName() string { return "clone" }

// notifyUser sends a message to the user. it's a tiny helper, but notifyUser is called
// frequently. refactoring appreciates it.
func (c *Command) notifyUser(msg string) {
	c.ui.Update(statusTextUpdate(msg))
	styles.OKTextf(msg)
}

// cloneRepo runs git clone command in localDir.
// git's stdout and stderr are captured in the first return.
// hard errors and no-zero exits are returned by err.
func (c *Command) cloneRepo() ([]byte, error) {
	var out bytes.Buffer

	cmd := exec.Command("git", "clone", c.repoURL, c.localDir)
	c.runningCmd = cmd
	defer func() { c.runningCmd = nil }()

	cmd.Stdout = &out
	cmd.Stderr = &out
	cmd.Dir = c.localDir

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error running git clone: %w", err)
	}

	ret := out.Bytes()[:64]
	ret = append(ret, []byte("...")...)

	return ret, nil
}
