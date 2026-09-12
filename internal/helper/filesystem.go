// SPDX-License-Identifier: GPL-2.0-only
// Copyright © 2026 Doğukan Meral <dogukan.meral@protonmail.com>

package helper

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dogukanmeral/scx-adapt/internal/paths"
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

// Creates the lock file and writes the active profile name into it
func CreateLock(profileName string) error {
	if err := os.WriteFile(paths.LOCKFILEPATH, []byte(profileName), 0700); err != nil {
		return fmt.Errorf("Error: Creating lock file at '%s': %s", paths.LOCKFILEPATH, err)
	}

	return nil
}

// Returns the name of the active profile written in the lock file
func ReadLock() (string, error) {
	data, err := os.ReadFile(paths.LOCKFILEPATH)
	if err != nil {
		return "", fmt.Errorf("Error: Reading lock file at '%s': %s", paths.LOCKFILEPATH, err)
	}

	return strings.TrimSpace(string(data)), nil
}

// Returns if file exists or not
func IsFileExist(path string) bool {
	_, err := os.Open(path)

	return err == nil
}

// Create log file
func CreateLogFile(schedulerName string) (*os.File, error) {
	if err := CreateDirIfNotExist(paths.LOGFOLDER); err != nil {
		return nil, err
	}

	currentTime := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_%s.log", currentTime, schedulerName)
	filepath := path.Join(paths.LOGFOLDER, filename)

	logFile, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
	if err != nil {
		return nil, fmt.Errorf("Creating log file at %s for scheduler %s failed: %s", filepath, schedulerName, err)
	}

	return logFile, nil
}
