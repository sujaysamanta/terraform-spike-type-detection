package tf_cmd

import (
	logger "github.com/rs/zerolog/log"
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
