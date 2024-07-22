package git

import (
	"os/exec"
	"strings"
)

func IsGitRepo(dir string) (bool, error) {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--is-inside-work-tree")
	stdout, err := cmd.Output()
	return strings.TrimSpace(string(stdout)) == "true", err
}
