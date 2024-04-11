package utils

import (
	"fmt"
	logger "github.com/rs/zerolog/log"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	cmd "terraform-spike-type-detection/tf-cmd"
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
	return strings.TrimSuffix(gitPath, "/")
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
	return strings.TrimSuffix(terraformPath, "/")
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

func GetMap(filePaths []FilePath) (map[string][]string, error) {
	fileMap := make(map[string][]string)
	for _, path := range filePaths {
		if path.isGitFolderPath() {
			gitPath := path.getGitFolderPath()
			updateMap(fileMap, "git", gitPath)
			continue
		}

		if path.isTerraformFolderPath() {
			terraformPath := path.getTerraformFolderPath()
			if err := getWorkspaces(terraformPath, fileMap["git"][0], fileMap); err != nil {
				logger.Error().Msgf("Error getting workspaces: %v", err)
				return nil, err
			}
		}
	}
	return fileMap, nil
}

// TODO: Refactor this function to reduce code duplication
func getWorkspaces(terraformPath string, gitPath string, fileMap map[string][]string) error {
	if strings.Compare(terraformPath, gitPath) == 0 {
		workspaces, err := cmd.GetWorkspaces()
		if err != nil {
			logger.Error().Msgf("Error getting workspaces: %v", err)
			return err
		}

		for _, workspace := range workspaces {
			if workspace != "" {
				workspace := terraformPath + "_" + workspace
				updateMap(fileMap, "terraform", workspace)
			}
		}

	} else {
		logger.Info().Msgf("Changing to directory: %s", terraformPath)
		err := os.Chdir(terraformPath)
		if err != nil {
			logger.Error().Msgf("Failed to change directory to %s: %v", terraformPath, err)
			return err
		}

		workspaces, err := cmd.GetWorkspaces()
		if err != nil {
			logger.Error().Msgf("Error getting workspaces: %v", err)
			return err
		}

		for _, workspace := range workspaces {
			if workspace != "" {
				workspace := gitPath + "_" + terraformPath + "_" + workspace
				sanitizedWorkspace := strings.ReplaceAll(workspace, "/", "_")
				updateMap(fileMap, "terraform", sanitizedWorkspace)
			}
		}
		parentPath := buildParentPath(terraformPath)
		logger.Info().Msgf("Changing to directory: %s", parentPath)
		err = os.Chdir(parentPath)

	}
	return nil
}

func updateMap(fileMap map[string][]string, key string, value string) {
	if _, ok := fileMap[key]; ok {
		paths := fileMap[key]
		paths = append(paths, value)
		fileMap[key] = paths
	} else {
		fileMap[key] = []string{value}
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

//func isChildPath(parentPath, childPath string) (bool, string, string, error) {
//	// Clean and resolve the absolute paths
//	absParentPath, err := filepath.Abs(filepath.Clean(parentPath))
//	if err != nil {
//		return false, "", "", err
//	}
//	absChildPath, err := filepath.Abs(filepath.Clean(childPath))
//	if err != nil {
//		return false, "", "", err
//	}
//
//	// Use Rel to find the relative path from parent to child
//	relPath, err := filepath.Rel(absParentPath, absChildPath)
//	if err != nil {
//		return false, "", "", err
//	}
//
//	// If the relative path starts with "..", the child is not within the parent
//	// Also, if the relative path is ".", it means both paths are the same
//	if !filepath.IsAbs(relPath) && !startsOrIsDotDot(relPath) {
//		return true, absParentPath, absChildPath, nil
//	}
//
//	return false, "", "", nil
//}
//
//// startsOrIsDotDot checks if the given path is ".." or starts with "../"
//func startsOrIsDotDot(path string) bool {
//	return path == ".." || filepath.HasPrefix(path, ".."+string(filepath.Separator))
//}

func buildParentPath(terraformPath string) string {
	// Split the path into parts
	parts := strings.Split(terraformPath, string(filepath.Separator))
	var walkBack = ""
	if len(parts) == 0 {
		walkBack = "../"
	} else {
		for i := 0; i < len(parts); i++ {
			walkBack += "../"
		}
	}
	return strings.Trim(walkBack, "")

}
