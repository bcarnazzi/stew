package cmd

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/urfave/cli/v3"
)

func Status(repository string) *cli.Command {
	return &cli.Command{
		Name:  "status",
		Usage: "get status of the git repository",
		Action: func(_ context.Context, _ *cli.Command) error {
			cmd := exec.Command("git", "status")
			cmd.Dir = repository
			output, err := cmd.CombinedOutput()
			if err != nil {
				return err
			}

			fmt.Printf("%s", output)
			return nil
		},
	}
}
