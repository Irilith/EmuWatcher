package utils

import (
	"EmuWatcher/utils/ui"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func ClearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func Exit(message string) {
	ui.OutShowFatalErrorModal(message, func() {
		os.Exit(0)
	})
}

func CheckFolderExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func CheckFileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func EnsureFolderExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("failed to create folder: %v", err)
		}
	}
	return nil
}

func EnsureFileExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create file: %v", err)
		}
		defer file.Close()
	}
	return nil
}

func ProgressBar(percentage int) string {
	barLength := 50
	filledLength := (percentage * barLength) / 100
	bar := make([]rune, barLength)
	for i := 0; i < filledLength; i++ {
		bar[i] = '#'
	}
	for i := filledLength; i < barLength; i++ {
		bar[i] = '-'
	}
	return string(bar)
}

func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	if err != nil {
		return false
	}
	return true
}

func GetFileHashSha256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("failed to copy file: %v", err)
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func SplitLines(in string) []string {
	clean := strings.ReplaceAll(in, "\r", "")
	return strings.Split(clean, "\n")
}

func Expect(s string, substr string) bool {
	return strings.Contains(s, substr)
}

func RemoveEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" && str != "\n" && str != "\r" && str != "\r\n" {
			r = append(r, str)
		}
	}
	return r
}
