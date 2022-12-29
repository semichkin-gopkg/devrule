package main

import (
	"flag"
	"github.com/hairyhenderson/gomplate/v3"
	"github.com/semichkin-gopkg/devrule/templates"
	"log"
)

func main() {
	configFilePath := flag.String("c", "", "configuration file path (required)")
	outputFilePath := flag.String("o", "", "output file path (required)")
	flag.Parse()

	if *configFilePath == "" || *outputFilePath == "" {
		log.Fatal("Please specify the configuration and output file paths using the -c and -o flags, respectively.")
	}

	err := gomplate.RunTemplates(&gomplate.Config{
		Input:       templates.Makefile,
		DataSources: []string{"configuration=" + *configFilePath},
		OutputFiles: []string{*outputFilePath},
	})
	if err != nil {
		log.Fatal(err)
	}
}
