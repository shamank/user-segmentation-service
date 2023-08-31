package main

import (
	"flag"
	"github.com/shamank/user-segmentation-service/internal/app"
)

const (
	configsDirs = "./configs"

	localFile = "/local.yaml"
	prodFile  = "/prod.yaml"
)

func main() {

	prodFlag := flag.Bool("prod", false, "product launch")

	configFile := localFile

	if *prodFlag {
		configFile = prodFile
	}

	app.Run(configsDirs + configFile)
}
