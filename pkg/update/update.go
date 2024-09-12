package update

import (
	"EmuWatcher/utils"
	"EmuWatcher/utils/version"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type ReleaseInfo struct {
	TagName string `json:"tag_name"`
}

func CheckForUpdate() error {
	latestVersion, err := GetLatestVersion()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	currentVersion := version.GetVersion()
	if latestVersion != "" && currentVersion != "DEV" && latestVersion != currentVersion {
		fmt.Printf("New version available: %s\n", latestVersion)
		fmt.Println("Please run Updater.exe to update to the latest version.")
		time.Sleep(5 * time.Second)
		utils.ClearScreen()
	}
	return nil
}

func GetLatestVersion() (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/Irilith/EmuWatcher/releases/latest")
	if err != nil {
		return "", fmt.Errorf("failed to get latest version: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	var data ReleaseInfo
	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return data.TagName, nil
}
