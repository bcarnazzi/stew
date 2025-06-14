package cmd

import (
	"context"
	"os"
	"path/filepath"
	"stew/utils"

	"github.com/urfave/cli/v3"
)

func Adopt(home, repository string) *cli.Command {
	return &cli.Command{
		Name:    "adopt",
		Usage:   "Adopt unmanaged dotfiles",
		Aliases: []string{"ad"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "",
				Usage:   "Package name",
			},
		},
		ArgsUsage: "TARGET...",
		Action: func(_ context.Context, cmd *cli.Command) error {
			var packageName string
			args := cmd.Args().Slice()

			if len(args) == 0 {
				cli.ShowSubcommandHelpAndExit(cmd, 1)
			}

			name := cmd.String("name")
			if len(args) == 1 {
				if name == "" {
					packageName = filepath.Base(args[0])
				} else {
					packageName = name
				}
			} else {
				if name == "" {
					utils.LogFatal("Package name must be provided when adopting multiple files")
				} else {
					packageName = name
				}
			}
			// utils.LogInfo("package name: " + packageName)

			for _, path := range args {
				absPath, err := filepath.Abs(path) // /home/xxx/.config/package
				if err != nil {
					return err
				}

				_, err = os.Stat(absPath)
				if err != nil {
					return err
				}

				// utils.LogInfo("abs path: " + absPath)

				relName, err := filepath.Rel(home, absPath) // .config/package
				if err != nil {
					return err
				}
				// utils.LogInfo("rel name: " + relName)

				dirName := filepath.Dir(relName) // .config
				utils.LogInfo("dir name: " + dirName)
				repoName := filepath.Join(repository, packageName, dirName) // /home/xxx/.dotfiles/package/.config
				// utils.LogInfo("repo name: " + repoName)

				err = os.MkdirAll(repoName, 0750)
				if err != nil {
					return err
				}

				destName := filepath.Join(repoName, filepath.Base(absPath))
				err = os.Rename(absPath, destName)
				if err != nil {
					return err
				}

			}

			return nil
		},
	}
}
