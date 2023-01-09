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
		Version:     "v0.0.11",
		Commands: []*cli.Command{
			{
				Name:    "init",
				Aliases: []string{"i"},
				Usage:   "Init configuration",

				Flags: []cli.Flag{
					&cli.PathFlag{
						Name:        "output",
						Aliases:     []string{"o"},
						Value:       "configuration.yaml",
						DefaultText: "configuration.yaml",
					},
				},
				Action: func(context *cli.Context) error {
					return os.WriteFile(context.Path("output"), []byte(templates.Configuration), 0666)
				},
			},
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
