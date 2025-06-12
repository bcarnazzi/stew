package cmd

import (
	"context"
	"os/exec"

	"github.com/urfave/cli/v3"
)

func Sync(repository string) *cli.Command {
	return &cli.Command{
		Name:  "sync",
		Usage: "Sync dotfiles to remote repository",
		Action: func(_ context.Context, _ *cli.Command) error {
			cmd := exec.Command("git", "pull")
			cmd.Dir = repository
			if err := cmd.Run(); err != nil {
				return err
			}

			cmd = exec.Command("git", "push")
			cmd.Dir = repository
			if err := cmd.Run(); err != nil {
				return err
			}

			return nil
		},
	}
}
