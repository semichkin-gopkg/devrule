package main

import (
	"github.com/hairyhenderson/gomplate/v3"
	"github.com/semichkin-gopkg/devrule/templates"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:        "devrule",
		Usage:       "Development Makefile builder",
		Description: "A tool for generating rules for managing a large number of local microservices",
		Version:     "v0.0.5",
		Commands: []*cli.Command{
			{
				Name:    "build",
				Aliases: []string{"b"},
				Usage:   "Builds Makefile",
				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "configuration",
						Aliases:     []string{"c"},
						Value:       "configuration.yaml",
						DefaultText: "configuration.yaml",
					},
					&cli.PathFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Value:       "Makefile",
						DefaultText: "Makefile",
					},
				},
				Action: func(context *cli.Context) error {
					log.Println(context.Path("configuration"), context.Path("output"))

					return gomplate.RunTemplates(&gomplate.Config{
						Input:       templates.Makefile,
						DataSources: []string{"configuration=" + context.Path("configuration")},
						OutputFiles: []string{context.Path("output")},
					})
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
