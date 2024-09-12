package cmd

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

const (
	startupFolderPath = "AppData/Roaming/Microsoft/Windows/Start Menu/Programs/Startup"
	batFileName       = "EmuWatcher.bat"
)

func CreateStartup() error {
	startupFolder, err := getStartupFolder()
	if err != nil {
		return fmt.Errorf("error getting startup folder: %w", err)
	}

	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error getting executable path: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	batFilePath := filepath.Join(startupFolder, batFileName)

	batContent := fmt.Sprintf(`@echo off
cd /d "%s"
start "" "%s" watch`, exeDir, exePath)

	if err := os.WriteFile(batFilePath, []byte(batContent), 0644); err != nil {
		return fmt.Errorf("error writing .bat file: %w", err)
	}
	return nil
}

func DeleteStartup() error {
	startupFolder, err := getStartupFolder()
	if err != nil {
		return fmt.Errorf("error getting startup folder: %w", err)
	}

	batFilePath := filepath.Join(startupFolder, batFileName)

	if err := os.Remove(batFilePath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file does not exist: %w", err)
		}
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}

func getStartupFolder() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("error getting current user: %w", err)
	}
	return filepath.Join(currentUser.HomeDir, startupFolderPath), nil
}
