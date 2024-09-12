package initialize

import (
	"EmuWatcher/adb"
	"EmuWatcher/roblox"
	"EmuWatcher/utils"
	"EmuWatcher/utils/cache"
	"EmuWatcher/utils/menu"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	totalSteps = 10
	scriptHash = "73c5d7c47f399816175ffc072774895d02ee32ef07e914a9f06d4022e986d33e"
)

func Initialize() {
	currentStep := 0

	updateProgress := func(status string) {
		currentStep++
		percentage := (currentStep * 100) / totalSteps
		line := fmt.Sprintf("\r[%-50s] %d%% - %s", utils.ProgressBar(percentage), percentage, status)
		fmt.Printf("\r%-100s", line)
		os.Stdout.Sync()
	}

	fmt.Println("Initializing...")
	fmt.Println("[--------------------------------------------------]")

	checks := []struct {
		name string
		fn   func() error
	}{
		{"Checking Operating System", checkOS},
		{"Checking permission", checkPermission},
		{"Checking datasets folder", checkDatasets},
		{"Checking tools folder", checkTools},
		{"Checking dependencies", checkDependencies},
		{"Checking data folder", checkDataFolder},
		{"Setting Environment", setupEnvironment},
		{"Initializing Adb server", initAdbServer},
		{"Refreshing device cache", refreshDeviceCache},
	}

	for _, check := range checks {
		updateProgress(check.name)
		if err := check.fn(); err != nil {
			utils.Exit(fmt.Sprintf("%s: %v", check.name, err))
		}
	}

	updateProgress("Initialization Complete!")
	fmt.Println("\n[--------------------------------------------------]")
	time.Sleep(1 * time.Second)
}

func checkOS() error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("this program only supports Windows")
	}
	return nil
}

func checkPermission() error {
	if !utils.IsAdmin() {
		return fmt.Errorf("please run the program as an administrator")
	}
	return nil
}

func checkDatasets() error {
	if err := utils.EnsureFolderExists("./assets/datasets"); err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}

	files := []string{"Cookies", "appStorage.json"}
	for _, file := range files {
		if !utils.CheckFileExists(filepath.Join("./assets/datasets", file)) {
			return fmt.Errorf("datasets file (%s) is not available, please reinstall the program (zip file) to fix this issue", file)
		}
	}
	return nil
}

func checkTools() error {
	if err := utils.EnsureFolderExists("./tools"); err != nil {
		return fmt.Errorf("error creating folder: %v", err)
	}
	return nil
}

func checkDependencies() error {
	dependencies := map[string]string{
		"./tools/adb":               "folder",
		"./tools/adb/adb.exe":       "file",
		"./tools/ocr":               "folder",
		"./tools/ocr/tesseract.exe": "file",
	}

	for path, kind := range dependencies {
		var exists bool
		if kind == "folder" {
			exists = utils.CheckFolderExists(path)
		} else {
			exists = utils.CheckFileExists(path)
		}
		if !exists {
			return fmt.Errorf("'%s' is not available, please follow the instruction to install it", path)
		}
	}
	return nil
}

func checkDataFolder() error {
	folders := []string{"./data", "./data/autoexec"}
	for _, folder := range folders {
		if err := utils.EnsureFolderExists(folder); err != nil {
			return fmt.Errorf("error creating folder: %v", err)
		}
	}

	if err := utils.EnsureFileExists("./data/cookies.txt"); err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	if !utils.CheckFileExists("./data/autoexec/EmuWatcher.lua") {
		if err := moveScript(); err != nil {
			return err
		}
	} else {
		hash, err := utils.GetFileHashSha256("./data/autoexec/EmuWatcher.lua")
		if err != nil {
			return fmt.Errorf("error getting file hash: %v", err)
		}
		if hash != scriptHash {
			if err := moveScript(); err != nil {
				return err
			}
		}
	}
	return nil
}

func setupEnvironment() error {
	envVars := map[string]string{
		"ADB_LOCAL_TRANSPORT_MAX_PORT": "65535",
		"ocr":                          "./tools/ocr/tesseract.exe",
		"adb":                          "./tools/adb/adb.exe",
	}

	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("error setting environment variable %s: %v", key, err)
		}
	}

	roblox.SetRunMenuCallback(menu.RunMenu)
	return nil
}

func initAdbServer() error {
	adb.RestartAdb()
	return nil
}

func refreshDeviceCache() error {
	return cache.RefreshDeviceCache()
}

func moveScript() error {
	scriptFile, err := os.Open("./scripts/EmuWatcher.lua")
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer scriptFile.Close()

	emuWatcherScript, err := os.Create("./data/autoexec/EmuWatcher.lua")
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	defer emuWatcherScript.Close()

	_, err = io.Copy(emuWatcherScript, scriptFile)
	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}
	return nil
}
