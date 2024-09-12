package emulator

import (
	"EmuWatcher/utils/config"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	LDPlayerPathConfig = "LDPlayerPath"
	BackupFilePath     = "temp/LdPlayerBackup.ldbk"
)

// CreateInstance creates multiple instances of the emulator with the same configuration and HardwareId
func CreateInstance(emuName string, quantity int) error {
	exist, err := isEmuNameExist(emuName)
	if err != nil {
		return fmt.Errorf("error checking if the emulator name exists: %w", err)
	}

	if !exist {
		return fmt.Errorf("emulator name does not exist, please create it first before cloning")
	}

	emuPath, err := config.GetConfig(LDPlayerPathConfig)
	if err != nil {
		return fmt.Errorf("error getting the emulator path: %w", err)
	}

	if _, err := os.Stat(BackupFilePath); os.IsNotExist(err) {
		if err = createBackup(emuName); err != nil {
			return fmt.Errorf("error creating backup: %w", err)
		}
	}

	for i := 0; i < quantity; i++ {
		fmt.Printf("Creating instance: %d\n", i)
		args := []string{"add", "--name", fmt.Sprintf("EmuWatcher_%d", i)}
		if err := runCommand(emuPath, args...); err != nil {
			return fmt.Errorf("error creating instance: %w", err)
		}
	}

	devices, err := getInstanceList()
	if err != nil {
		return fmt.Errorf("error getting instance list: %w", err)
	}

	for _, device := range devices {
		if strings.Contains(device, "EmuWatcher") && device != emuName {
			args := []string{"restore", "--name", device, "--file", BackupFilePath}
			if err := runCommand(emuPath, args...); err != nil {
				return fmt.Errorf("error restoring backup: %w", err)
			}
		}
	}

	if err := os.Remove(BackupFilePath); err != nil {
		fmt.Printf("Warning: Error deleting the backup file: %v\n", err)
	}

	return nil
}

func createBackup(emuName string) error {
	emuPath, err := config.GetConfig(LDPlayerPathConfig)
	if err != nil {
		return fmt.Errorf("error getting the emulator path: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(BackupFilePath), 0755); err != nil {
		return fmt.Errorf("error creating temp directory: %w", err)
	}

	args := []string{"backup", "--name", emuName, "--file", BackupFilePath}
	return runCommand(emuPath, args...)
}

func getInstanceList() ([]string, error) {
	emuPath, err := config.GetConfig(LDPlayerPathConfig)
	if err != nil {
		return nil, fmt.Errorf("error getting the emulator path: %w", err)
	}

	output, err := runCommandWithOutput(emuPath, "list")
	if err != nil {
		return nil, fmt.Errorf("error getting instance list: %w", err)
	}

	return strings.Split(strings.TrimSpace(output), "\n"), nil
}

func isEmuNameExist(emuName string) (bool, error) {
	instances, err := getInstanceList()
	if err != nil {
		return false, err
	}

	for _, instance := range instances {
		if instance == emuName {
			return true, nil
		}
	}

	return false, nil
}

func runCommand(emuPath string, args ...string) error {
	cmd := exec.Command(filepath.Join(emuPath, "dnconsole"), args...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, stderr.String())
	}
	return nil
}

func runCommandWithOutput(emuPath string, args ...string) (string, error) {
	cmd := exec.Command(filepath.Join(emuPath, "dnconsole"), args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%w: %s", err, string(output))
	}
	return string(output), nil
}
