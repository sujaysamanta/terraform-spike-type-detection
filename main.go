package main

import (
	logger "github.com/rs/zerolog/log"
	"strings"
	tfcSpecs "terraform-spike-type-detection/tfc-spec"
	utils "terraform-spike-type-detection/utils"
)

func main() {
	println("Hello, Terraform!")
	hiddenFilePaths, err := utils.FindHiddenFiles()
	if err != nil {
		logger.Error().Msgf("Error finding hidden files: %v", err)
	}

	fileMap := utils.GetMap(hiddenFilePaths)
	for key, val := range fileMap {
		logger.Info().Msgf("%s: [%s]", key, strings.Join(val, ", "))
	}

	logger.Info().Msgf("Generating project specs")
	projectSpec, err := tfcSpecs.GenerateProjectSpecs(fileMap)
	if err != nil {
		logger.Error().Msgf("Error generating project specs: %v", err)
	}
	logger.Info().Msgf("Project specs: %s", projectSpec)

	logger.Info().Msgf("Writing project specs to specs.yaml")
	err = utils.WriteSpec(projectSpec)
	if err != nil {
		logger.Error().Msgf("Error writing project specs yaml: %v", err)
	}
	logger.Info().Msgf("Project specs written to specs.yaml")

}
