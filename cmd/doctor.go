package cmd

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"stew/utils"

	"github.com/urfave/cli/v3"
)

func Doctor(home, repository string) *cli.Command {
	return &cli.Command{
		Name:  "doctor",
		Usage: "Check stew configuration and dependencies",
		Action: func(_ context.Context, _ *cli.Command) error {
			var errCode error

			if home == "" {
				utils.LogWarn("Undefined HOME")
				errCode = errors.New("undefined home")
			} else {
				utils.LogOk("home directory is " + home)
			}

			_, err := os.Stat(repository)
			if err != nil {
				utils.LogWarn("Cannot find repository at " + err.Error())
				errCode = err
			} else {
				utils.LogOk("repository is " + repository)
			}

			path, err := exec.LookPath("git")
			if err != nil {
				utils.LogWarn("git command not found")
				errCode = err
			}
			utils.LogOk("git command found at " + path)

			path, err = exec.LookPath("stow")
			if err != nil {
				utils.LogWarn("stow command not found")
				errCode = err
			}
			utils.LogOk("stow command found at " + path)

			return errCode
		},
	}
}
