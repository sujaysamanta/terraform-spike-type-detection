package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func main() {
	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	fmt.Println("Scanning for hidden directories in:", cwd)
	fmt.Println("Hidden directories found:")

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
				fmt.Println(relativePath)
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error walking through the directory:", err)
	}
}
