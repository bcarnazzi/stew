package main

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
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

	app := &cli.App{
		Name:  "stew",
		Usage: "A simple dotfiles manager",
		Commands: []*cli.Command{
			{
				Name:  "list",
				Usage: "List managed dotfiles",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()

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
				Name:  "adopt",
				Usage: "Adopt unmanaged dotfiles",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
					if len(args) < 2 {
						return fmt.Errorf("need at least two arguments")
					}

					sources := args[:len(args)-1]
					destination := filepath.Join(repository, args[len(args)-1])

					_, err := os.Stat(destination)
					if os.IsNotExist(err) {
						err = os.Mkdir(destination, 0755)
						if err != nil {
							return err
						}
					}

					for _, source := range sources {
						// Check if symlink to stew repository
						info, err := os.Lstat(source)
						if err != nil {
							log.Fatal(err)
						}

						if info.Mode()&os.ModeSymlink != 0 {
							linkDest, err := os.Readlink(source)
							if err != nil {
								return err
							}

							absLinkDest, err := filepath.Abs(filepath.Join(filepath.Dir(source), linkDest))
							if err != nil {
								return err
							}

							// logInfo("Checking symlink: " + absLinkDest + " -> " + stewRepository)
							if strings.HasPrefix(absLinkDest, stewRepository) {
								logWarn(source + " is already managed. Skipping.")
								continue
							}
						}
						// logInfo("Not a symlink: " + source)

						absSource, err := filepath.Abs(source)
						if err != nil {
							return err
						}

						relSource, err := filepath.Rel(home, absSource)
						if err != nil {
							return err
						}

						absDest := filepath.Join(destination, relSource)

						_ = os.MkdirAll(filepath.Dir(absDest), 0755)

						logInfo("Moving " + source + " to " + absDest)
						err = os.Rename(source, absDest)
						if err != nil {
							return err
						}

						relDest, err := filepath.Rel(filepath.Dir(absSource), absDest)
						if err != nil {
							return err
						}

						logInfo("Linking " + relDest + " to " + source)
						err = os.Symlink(relDest, source)
						if err != nil {
							return err
						}
					}

					return nil
				},
			},
			{
				Name:  "link",
				Usage: "Link managed dotfiles",
				Action: func(c *cli.Context) error {
					args := c.Args().Slice()
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
				Name:  "doctor",
				Usage: "Check stew configuration and dependencies",
				Action: func(c *cli.Context) error {
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
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
