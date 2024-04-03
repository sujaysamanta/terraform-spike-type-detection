package tfcSpecs

import (
	logger "github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"strings"
)

type Workspace string
type ProjectSpec struct {
	Project    string
	Workspaces []Workspace
}

type workspaceOps interface {
	String() string
}

func (w Workspace) String() string {
	return string(w)
}

func ToWorkspace(s string) Workspace {
	return Workspace(s)
}

func buildProjectSpecs(fileMap map[string][]string) ProjectSpec {
	var projectSpec ProjectSpec
	for key, val := range fileMap {
		if strings.Compare(key, "git") == 0 {
			projectSpec.Project = val[0]
		}
		if strings.Compare(key, "terraform") == 0 {
			for _, workspace := range val {
				projectSpec.Workspaces = append(projectSpec.Workspaces, ToWorkspace(workspace))
			}
		}
	}
	return projectSpec
}

func GenerateProjectSpecs(fileMap map[string][]string) ([]byte, error) {
	projectSpecs := buildProjectSpecs(fileMap)

	yamlSpecs, err := yaml.Marshal(projectSpecs)
	if err != nil {
		logger.Error().Msgf("Error generating project specs yaml: %v", err)
		return nil, err
	}

	return yamlSpecs, nil
}
