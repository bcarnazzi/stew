package main

import (
	"context"
	"os"
	"path/filepath"

	"stew/cmd"
	"stew/tools"

	"github.com/urfave/cli/v3"
)

func main() {
	home := os.Getenv("HOME")

	stewRepository := os.Getenv("STEW_REPOSITORY")
	if stewRepository == "" {
		stewRepository = ".dotfiles"
	}

	repository := filepath.Join(home, stewRepository)

	stew := &cli.Command{
		Name:  "stew",
		Usage: "A simple dotfiles manager",
		Commands: []*cli.Command{
			cmd.Adopt(home, repository),
			cmd.Doctor(home, repository),
			cmd.Link(repository),
			cmd.List(repository),
			cmd.Status(repository),
			cmd.Sync(repository),
		},
	}

	if err := stew.Run(context.Background(), os.Args); err != nil {
		tools.LogFatal(err)
	}
}
