package git

import (
	"os/exec"
	"regexp"
	"strings"
)

type LogParams struct {
	Dir        string
	Since      string
	Email      string
	CommitLink string
}

func GetLogParams(dir, targetRemote, since string) (LogParams, error) {
	params := LogParams{Dir: dir, Since: since}
	var err error

	if params.Email, err = GetEmail(dir); err != nil {
		return params, err
	}

	remote, err := GetRemote(dir, targetRemote)
	params.CommitLink = remoteToLink(remote)
	return params, err
}

// https://stackoverflow.com/a/7786922

func GetEmail(dir string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "config", "user.email")
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

func GetRemote(dir, targetRemote string) (string, error) {
	cmd := exec.Command("git", "-C", dir, "remote", "get-url", targetRemote)
	stdout, err := cmd.Output()
	return strings.TrimSpace(string(stdout)), err
}

func remoteToLink(remoteLink string) (commitLink string) {
	substitutions := []struct {
		pattern     string
		replacement string
	}{
		{`^(git@|https?://)github.com(:|/)`, `https://github.com/`},
		{`(/|\.git)?$`, `/commit/`},
	}

	commitLink = strings.TrimSpace(remoteLink)
	for _, sub := range substitutions {
		regex := regexp.MustCompile(sub.pattern)
		commitLink = regex.ReplaceAllString(commitLink, sub.replacement)
	}

	return
}

type strerr struct {
	str string
	err error
}
