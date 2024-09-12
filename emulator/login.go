package emulator

import (
	"EmuWatcher/adb"
	"EmuWatcher/roblox"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	cookiesDBPath    = "./assets/datasets/Cookies"
	emuWatcherPath   = "./assets/datasets/Emu.Watcher"
	appStoragePath   = "./assets/datasets/appStorage.json"
	robloxClientPath = "/data/data/com.roblox.client"
)

func Login(emuName, cookie string) (bool, error) {
	valid, userInfo, err := roblox.ValidCookies(cookie)
	if err != nil {
		return false, fmt.Errorf("error validating cookies: %w", err)
	}
	if !valid {
		return false, fmt.Errorf("invalid cookies")
	}

	fmt.Printf("Logging in user %s...\n", userInfo.DisplayName)

	if err := updateCookiesDB(cookie); err != nil {
		return false, err
	}

	if err := adb.StopGame(emuName); err != nil {
		return false, fmt.Errorf("error stopping game: %w", err)
	}

	if err := updateEmuWatcherFile(userInfo.Name); err != nil {
		return false, err
	}

	if err := pushFiles(emuName); err != nil {
		return false, err
	}

	fmt.Printf("Successfully logged in %s\n", userInfo.DisplayName)
	return true, nil
}

func updateCookiesDB(cookie string) error {
	db, err := sql.Open("sqlite3", cookiesDBPath)
	if err != nil {
		return fmt.Errorf("error opening database: %w", err)
	}
	defer db.Close()

	_, err = db.Exec("UPDATE cookies SET value=? WHERE name='.ROBLOSECURITY'", cookie)
	if err != nil {
		return fmt.Errorf("error updating cookies: %w", err)
	}
	return nil
}

func updateEmuWatcherFile(userName string) error {
	file, err := os.OpenFile(emuWatcherPath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(fmt.Sprintf("%s|0", userName))
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}
	return nil
}

func pushFiles(emuName string) error {
	filesToPush := map[string]string{
		cookiesDBPath:  filepath.Join(robloxClientPath, "app_webview/Default/Cookies"),
		appStoragePath: filepath.Join(robloxClientPath, "files/appData/LocalStorage/appStorage.json"),
		emuWatcherPath: "/sdcard/Codex/Workspace/Emu.Watcher",
	}

	for src, dest := range filesToPush {
		args := []string{"-s", emuName, "push", src, dest}
		out, err := adb.ExecuteADBCommand(args)
		if err != nil {
			return fmt.Errorf("error pushing %s: %w, output: %s", filepath.Base(src), err, string(out))
		}
	}

	return nil
}
