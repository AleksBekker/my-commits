package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// https://stackoverflow.com/a/7786922

func GetCommits(params LogParams) ([]string, error) {
	cmd := exec.Command("git", "-C", params.Dir, "log",
		fmt.Sprintf(`--pretty=format:%s%%H`, params.CommitLink),
		fmt.Sprintf(`--author=%s`, params.Email),
		fmt.Sprintf(`--since=%s`, params.Since), "--all")

	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	links := strings.Split(string(stdout), "\n")
	for idx, link := range links {
		links[idx] = strings.TrimSpace(link)
	}

	return links, nil
}
