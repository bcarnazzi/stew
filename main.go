package main

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v3"
	"slices"
)

var (
	yellow = color.New(color.FgYellow).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
)

func logInfo(s string) {
	log.Printf("["+green("INFO")+"] %s\n", s)
}

func logWarn(s string) {
	log.Printf("["+yellow("WARN")+"] %s\n", s)
}

func logOk(s string) {
	log.Printf("["+green("OK")+"] %s\n", s)
}

func main() {
	home := os.Getenv("HOME")

	stewRepository := os.Getenv("STEW_REPOSITORY")
	if stewRepository == "" {
		stewRepository = ".dotfiles"
	}

	repository := filepath.Join(home, stewRepository)

	cmd := &cli.Command{
		Name:  "stew",
		Usage: "A simple dotfiles manager",
		Commands: []*cli.Command{
			{
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
				Action: func(_ context.Context, cmd *cli.Command) error {
					var packageName string
					args := cmd.Args().Slice()

					name := cmd.String("name")
					if len(args) == 1 {
						if name == "" {
							packageName = filepath.Base(args[0])
						} else {
							packageName = name
						}
					} else {
						if name == "" {
							log.Fatal("Package name must be provided when adopting multiple files")
						} else {
							packageName = name
						}
					}
					// logInfo("package name: " + packageName)

					for _, path := range args {
						absPath, err := filepath.Abs(path) // /home/xxx/.config/package
						if err != nil {
							return err
						}

						_, err = os.Stat(absPath)
						if err != nil {
							return err
						}

						// logInfo("abs path: " + absPath)

						relName, err := filepath.Rel(home, absPath) // .config/package
						if err != nil {
							return err
						}
						// logInfo("rel name: " + relName)

						dirName := filepath.Dir(relName) // .config
						logInfo("dir name: " + dirName)
						repoName := filepath.Join(repository, packageName, dirName) // /home/xxx/.dotfiles/package/.config
						// logInfo("repo name: " + repoName)

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
			},
			{
				Name:  "doctor",
				Usage: "Check stew configuration and dependencies",
				Action: func(_ context.Context, _ *cli.Command) error {
					var errCode error

					if home == "" {
						logWarn("Undefined HOME")
						errCode = errors.New("undefined home")
					} else {
						logOk("home directory is " + home)
					}

					_, err := os.Stat(repository)
					if err != nil {
						logWarn("Cannot find repository at " + err.Error())
						errCode = err
					} else {
						logOk("repository is " + repository)
					}

					path, err := exec.LookPath("git")
					if err != nil {
						logWarn("git command not found")
						errCode = err
					}
					logOk("git command found at " + path)

					path, err = exec.LookPath("stow")
					if err != nil {
						logWarn("stow command not found")
						errCode = err
					}
					logOk("stow command found at " + path)

					return errCode
				},
			},
			{
				Name:    "link",
				Usage:   "Link managed dotfiles",
				Aliases: []string{"ln"},
				Action: func(_ context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()
					var errCode error
					if len(args) == 0 {
						return fmt.Errorf("need at least one arguments")
					}

					for _, p := range args {
						cmd := exec.Command("stow", "-d", repository, p)
						if err := cmd.Run(); err != nil {
							logWarn("Cannot link " + p)
							errCode = err
						} else {
							logInfo(p + " linked")
						}
					}
					return errCode
				},
			},
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List managed dotfiles",
				Action: func(_ context.Context, cmd *cli.Command) error {
					args := cmd.Args().Slice()

					entries, err := os.ReadDir(repository)
					if err != nil {
						log.Fatal(err)
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
			},
			{
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
			},
			{
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
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
