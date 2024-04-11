package tf_cmd

import (
	"fmt"
	logger "github.com/rs/zerolog/log"
	"os"
	"os/exec"
	"strings"
)

func GetWorkspaces() ([]string, error) {
	logger.Info().Msgf("Getting workspaces")
	cmd := exec.Command("terraform", "workspace", "list")
	out, err := cmd.Output()
	if err != nil {
		logger.Error().Msgf("Error getting workspaces: %v", err)
		return nil, err
	}
	output := strings.ReplaceAll(strings.ReplaceAll(string(out), "*", ""), " ", "")
	workspaces := strings.Split(output, "\n")
	return workspaces, nil
}

func GoToChildDirectory(directory string) error {
	// Name of the child directory
	childDirectory := "childDirectoryName"

	// Changing to the child directory
	err := os.Chdir(childDirectory)
	if err != nil {
		fmt.Printf("Failed to change directory to %s: %s\n", childDirectory, err)
		return err
	}

	fmt.Printf("Successfully changed to directory %s\n", childDirectory)
	return nil
}
