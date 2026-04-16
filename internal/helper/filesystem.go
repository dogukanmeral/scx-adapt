package helper

import (
	"fmt"
	"os"

	paths "github.com/dogukanmeral/scx-adapt/internal"
)

// Creates directory with permission '700' if it does not exist already.
func CreateDirIfNotExist(dirPath string) error {
	if !IsFileExist(dirPath) {
		err := os.Mkdir(dirPath, 0700)

		if err != nil {
			return fmt.Errorf("Error occured while creating directory '%s': %s", dirPath, err)
		}
	}

	return nil
}

func CopyFile(sourcePath string, destinationPath string, filePerm int) error {
	input, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}

	err = os.WriteFile(destinationPath, input, os.FileMode(filePerm))
	if err != nil {
		return err
	}

	return nil
}

// Writes data to file with permissions '0644'
func Write(path string, data string) {
	err := os.WriteFile(path, []byte(data), 0644)

	if err != nil {
		panic(err) // TODO: Convert to error returning function as other components of project.
	}
}

// Removes the lock file
func RemoveLock() error {
	if err := os.Remove(paths.LOCKFILEPATH); err != nil {
		return fmt.Errorf("Error: Removing lock file at '%s' failed: %s", paths.LOCKFILEPATH, err)
	}

	return nil
}

// Creates the lock file
func CreateLock() error {
	if _, err := os.Create(paths.LOCKFILEPATH); err != nil {
		return fmt.Errorf("Error: Creating lock file at '%s': %s", paths.LOCKFILEPATH, err)
	}

	return nil
}

// Returns if file exists or not
func IsFileExist(path string) bool {
	_, err := os.Open(path)

	return err == nil
}
