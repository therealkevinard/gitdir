package dirtools

import (
	"fmt"
	"net/url"
	"path"
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
		return "", fmt.Errorf("invalid url: %s (sanitized to %s)", repoURL, cleanRepoURL)
	}

	// build final directory
	dir := path.Join(parsed.Host, strings.TrimSuffix(parsed.Path, ".git"))

	return dir, nil
}
