package main

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sync/errgroup"

	cloneorg "github.com/caarlos0/clone-org"
	"github.com/caarlos0/spin"
	"github.com/urfave/cli"
)

var version = "master"

func main() {
	app := cli.NewApp()
	app.Name = "clone-org"
	app.Usage = "Clone all repos of a github organization"
	app.Version = version
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "org, o",
		},
		cli.StringFlag{
			Name:   "token, t",
			EnvVar: "GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name: "destination, d",
		},
	}
	app.Action = func(c *cli.Context) error {
		var token = c.String("token")
		var org = c.String("org")
		if token == "" {
			return cli.NewExitError("missing github token", 1)
		}
		if org == "" {
			return cli.NewExitError("missing organization name", 1)
		}
		destination := c.String("destination")
		if destination == "" {
			destination = filepath.Join(os.TempDir(), org)
		}
		fmt.Printf("Destination: %v\n", destination)
		var s = spin.New("%s Finding repositories to clone...")
		s.Start()
		repos, err := cloneorg.AllOrgRepos(token, org)
		s.Stop()
		if err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		if err := os.Mkdir(destination, 0700); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		s = spin.New(fmt.Sprintf(
			"%v Cloning %v repositories...", "%v", len(repos),
		))
		s.Start()
		defer s.Stop()
		var g errgroup.Group
		for _, repo := range repos {
			repo := repo
			g.Go(func() error {
				return cloneorg.Clone(repo, destination)
			})
		}
		return g.Wait()
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
