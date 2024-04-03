package main

import (
	logger "github.com/rs/zerolog/log"
	"strings"
	uttils "terraform-spike-type-detection/utils"
)

func main() {
	println("Hello, Terraform!")
	hiddenFilePaths, err := uttils.FindHiddenFiles()
	if err != nil {
		logger.Error().Msgf("Error finding hidden files: %v", err)
	}

	fileMap := uttils.GetMap(hiddenFilePaths)
	for key, val := range fileMap {
		logger.Info().Msgf("%s: [%s]", key, strings.Join(val, ", "))
	}

}
