package cmd

import (
	"context"
	"os/exec"
	"stew/tools"

	"github.com/urfave/cli/v3"
)

func Link(repository string) *cli.Command {
	return &cli.Command{
		Name:      "link",
		Usage:     "Link managed dotfiles",
		Aliases:   []string{"ln"},
		ArgsUsage: "PACKAGE ...",
		Action: func(_ context.Context, cmd *cli.Command) error {
			args := cmd.Args().Slice()
			var errCode error
			if len(args) == 0 {
				cli.ShowSubcommandHelpAndExit(cmd, 1)
			}

			for _, p := range args {
				cmd := exec.Command("stow", "-d", repository, p)
				if err := cmd.Run(); err != nil {
					tools.LogWarn("Cannot link " + p)
					errCode = err
				} else {
					tools.LogInfo(p + " linked")
				}
			}
			return errCode
		},
	}
}
