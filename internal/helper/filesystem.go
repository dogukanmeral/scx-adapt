package helper

import (
	"fmt"
	"internal/checks"
	"os"
)

// Creates directory with permission '700' if it does not exist already.
func CreateDirIfNotExist(dirPath string) error {
	if !checks.IsFileExist(dirPath) {
		err := os.Mkdir(dirPath, 0700)

		if err != nil {
			return fmt.Errorf("Error occured while creating directory '%s': %s", dirPath, err)
		}
	}

	return nil
}
