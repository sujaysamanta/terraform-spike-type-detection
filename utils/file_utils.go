package utils

import (
	"fmt"
	logger "github.com/rs/zerolog/log"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type FilePath string

type FilePathOps interface {
	string() string
	isGitFolderPath() bool
	isTerraformFolderPath() bool
	getGitFolderPath() string
	getTerraformFolderPath() string
}

func toFilePath(s string) FilePath {
	return FilePath(s)
}

func getCurrentFolder() (string, error) {
	currentWorkingDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %v", err)
	}
	return currentWorkingDir, nil
}

func (fp FilePath) string() string {
	return string(fp)
}

func (fp FilePath) isGitFolderPath() bool {
	return strings.Contains(string(fp), ".git")
}

func (fp FilePath) isTerraformFolderPath() bool {
	return strings.Contains(string(fp), ".terraform")
}

func (fp FilePath) getGitFolderPath() string {
	gitPath := strings.Replace(fp.string(), ".git", "", -1)
	logger.Info().Msgf("Git path: %s", gitPath)
	// If the path is empty, it means the current directory is the root of git folder
	if "" == gitPath {
		currentWorkingDir, err := getCurrentFolder()
		if err == nil {
			folderTree := strings.Split(currentWorkingDir, "/")
			return folderTree[len(folderTree)-1]
		}
	}
	return gitPath
}

func (fp FilePath) getTerraformFolderPath() string {
	terraformPath := strings.Replace(fp.string(), ".terraform", "", -1)
	logger.Info().Msgf("Terraform path: %s", terraformPath)
	// If the path is empty, it means the current directory is the root of terraform folder
	if "" == terraformPath {
		currentWorkingDir, err := getCurrentFolder()
		if err == nil {
			folderTree := strings.Split(currentWorkingDir, "/")
			return folderTree[len(folderTree)-1]
		}
	}
	return terraformPath
}

func FindHiddenFiles() ([]FilePath, error) {
	var hiddenDirs []FilePath

	// Get the current directory
	cwd, err := getCurrentFolder()
	if err != nil {
		return nil, err
	}

	logger.Info().Msgf("Scanning for hidden directories in: %s", cwd)
	logger.Info().Msg("Hidden directories found:")

	// Walk the current directory
	err = filepath.WalkDir(cwd, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		// Check if it is a directory and if the name starts with a dot "."
		if d.IsDir() && filepath.Base(path)[0] == '.' {
			// Ensure it's not the current directory itself
			if path != cwd {
				relativePath, err := filepath.Rel(cwd, path)
				if err != nil {
					fmt.Println("Error calculating relative path:", err)
					return err
				}
				logger.Info().Msg(relativePath)
				hiddenDirs = append(hiddenDirs, toFilePath(relativePath))
			}
		}

		return nil
	})

	if err != nil {
		logger.Error().Msgf("Error walking through the directory: %v", err)
	}

	return hiddenDirs, err
}

func GetMap(filePaths []FilePath) map[string][]string {
	fileMap := make(map[string][]string)
	for _, path := range filePaths {
		if path.isGitFolderPath() {
			gitPath := path.getGitFolderPath()
			updateMap(fileMap, "git", gitPath)
			continue
		}

		if path.isTerraformFolderPath() {
			terraformPath := path.getTerraformFolderPath()
			updateMap(fileMap, "terraform", terraformPath)
		}
	}
	return fileMap
}

func updateMap(fileMap map[string][]string, key string, value string) {
	if _, ok := fileMap[key]; ok {
		paths := fileMap[key]
		paths = append(paths, strings.TrimSuffix(value, "/"))
		fileMap[key] = paths
	} else {
		fileMap[key] = []string{strings.TrimSuffix(value, "/")}
	}
}

func WriteSpec(yamlFile []byte) error {
	workingDir, err := getCurrentFolder()
	if err != nil {
		logger.Error().Msgf("Error getting current working directory: %v", err)
		return err
	}
	f, err := os.Create(workingDir + "/" + "spec.yaml")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Error().Msgf("Error closing file: %v", err)
		}
	}(f)

	_, err = io.WriteString(f, string(yamlFile))
	if err != nil {
		logger.Error().Msgf("Error writing to file: %v", err)
		return err
	}

	return nil
}
