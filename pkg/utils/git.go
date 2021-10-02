package utils

import "strings"

func GetCommitHash(arg string) (string, error) {
	wd := &WorkDir{}
	stdout, _, err := wd.RunCommand("git", "rev-parse", "-q", "--short", arg)
	if err != nil {
		return "", err
	}
	return strings.Trim(stdout, "\n"), nil
}
