package dirtools

import (
	"fmt"
	"github.com/therealkevinard/gitdir/errors"
	"io/fs"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

// regexes.
var (
	// ensure protocol is part of the url.
	rexNeedsProto = regexp.MustCompile(`(?mi)^https://|ssh://`)
	// strip port - replace with "".
	rexPortURL = regexp.MustCompile(`(?mi):([0-9]+)/`)
	// swap slash - urls like `git@github.com:group/project.git` become `git@github.com/group/project.git`.
	rexPortlessURL = regexp.MustCompile(`(?mi):([a-z]+)/`)
)

// CompileDirPath prefixes a directory tree with constant roots.
// it expects input to be normalized with NormalizeRepoURL.
func CompileDirPath(root, repoDirTree string) string {
	return path.Clean(path.Join(root, repoDirTree))
}

// NormalizeRepoURL normalizes repoURL to a partial directory tree.
// repoURL is a git clone url, and resulting string is safe for os.MkdirAll
// it handles removing/replacing ports, protocols, git users, .git suffix, and other conversions.
func NormalizeRepoURL(repoURL string) (string, error) {
	cleanRepoURL := repoURL

	// urls without protocol fail hard in url.Parse. stub a fake protocol to allow through.
	if !rexNeedsProto.MatchString(cleanRepoURL) {
		cleanRepoURL = "xxx://" + cleanRepoURL
	}

	// normalize colons
	cleanRepoURL = rexPortURL.ReplaceAllString(cleanRepoURL, "/")
	cleanRepoURL = rexPortlessURL.ReplaceAllString(cleanRepoURL, "/$1/")

	// parse normalized url
	parsed, err := url.Parse(cleanRepoURL)
	if err != nil {
		return "", fmt.Errorf("error parsing url: %w", err)
	}
	if parsed.Host == "" || parsed.Path == "" {
		return "", fmt.Errorf("invalid url %s (sanitized to %s): %w", repoURL, cleanRepoURL, errors.ErrInvalidURL)
	}

	// build final directory
	dir := path.Join(parsed.Host, strings.TrimSuffix(parsed.Path, ".git"))

	return dir, nil
}

// FindGitDirs uses filepath.Walk to recursively identify git repos, returning the slice of git paths
// repos are identified as `parent directory of a .git directory`.
func FindGitDirs(root string) ([]string, error) {
	items := make([]string, 0)
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() && d.Name() == ".git" {
			items = append(items, strings.TrimSuffix(path, "/.git"))
			return filepath.SkipDir
		}

		return nil
	})

	//nolint: wrapcheck
	return items, err
}

// permissions used for created bashes. should be executable.
const scriptPerms = 0o750

// WriteExecFile is used to write temporary shell scripts with +x . It should not be used for general file-writing.
func WriteExecFile(scriptPath, content string) error {
	_ = os.MkdirAll(path.Dir(scriptPath), scriptPerms) // TODO: check error
	f, fileErr := os.Create(scriptPath)
	if fileErr != nil {
		return fmt.Errorf("error creating script file: %w", fileErr)
	}
	defer func() { _ = f.Close() }()

	_, _ = f.WriteString(content)

	return nil
}
