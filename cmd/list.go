package cmd

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"stew/utils"
	"strings"

	"github.com/urfave/cli/v3"
)

func List(repository string) *cli.Command {
	return &cli.Command{
		Name:      "list",
		Aliases:   []string{"ls"},
		Usage:     "List managed dotfiles",
		ArgsUsage: "[PACKAGE...]",
		Action: func(_ context.Context, cmd *cli.Command) error {
			args := cmd.Args().Slice()

			entries, err := os.ReadDir(repository)
			if err != nil {
				utils.LogFatal(err)
			}

			for _, entry := range entries {
				name := entry.Name()
				if entry.IsDir() && !strings.HasPrefix(name, ".") {
					display := len(args) == 0
					if !display {
						if slices.Contains(args, name) {
							display = true
						}
					}
					if display {
						fmt.Printf("%s:\n", name)
					} else {
						continue
					}
				} else {
					continue
				}

				dotFilePath := filepath.Join(repository, name)
				err = filepath.Walk(dotFilePath, func(path string, info fs.FileInfo, err error) error {
					if err != nil {
						return err
					}

					if info.IsDir() {
						return nil
					}

					relPath, err := filepath.Rel(dotFilePath, path)
					if err != nil {
						return err
					}

					if relPath != "." {
						fmt.Printf("  %s\n", relPath)
					}

					return nil
				})
				if err != nil {
					return err
				}

			}
			return nil
		},
	}
}
