package menu

import (
	"EmuWatcher/emulator"
	"EmuWatcher/pkg/commands/login"
	"EmuWatcher/roblox"
	"EmuWatcher/utils"
	"EmuWatcher/utils/cmd"
	"EmuWatcher/utils/config"
	"EmuWatcher/utils/version"
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	menuOptionWatchEmulator  = 1
	menuOptionConfigurations = 2
	menuOptionAddCookies     = 3
	menuOptionLogin          = 4
	menuOptionSetupAutoexec  = 5
	menuOptionExit           = 9
	menuOptionSetupStartup   = 10
)

func RunMenu() {
	for {
		displayHeader()
		displayMenuOptions()

		choice := getUserChoice()
		switch choice {
		case menuOptionWatchEmulator:
			handleWatchEmulator()
		case menuOptionConfigurations:
			config.EditConfig()
		case menuOptionAddCookies:
			roblox.AddCookies()
		case menuOptionLogin:
			login.Login()
		case menuOptionSetupAutoexec:
			emulator.SetupAutoExec()
		case menuOptionExit:
			os.Exit(0)
		case menuOptionSetupStartup:
			handleSetupStartup()
		default:
			handleInvalidChoice()
		}
	}
}

func displayHeader() {
	utils.ClearScreen()
	fmt.Println("EmuWatcher |", version.GetVersion())
	fmt.Println("Commit:", version.GetCommit())
	fmt.Println("Repository: github.com/Irilith/EmuWatcher")
	fmt.Println("License: GNU GPL-3.0, LICENSE URL: https://raw.githubusercontent.com/Irilith/EmuWatcher/main/LICENSE")
	fmt.Println("Discord: discord.gg/QfpGHB87jK")
	fmt.Println("Please enable ADB debugging and Root in Your LdPlayer config")
	fmt.Println("-----------------------------")
}

func displayMenuOptions() {
	fmt.Println("Select an option:")
	fmt.Println("1. Watch Emulator")
	fmt.Println("2. Configurations")
	fmt.Println("3. Add Cookies")
	fmt.Println("4. Login (no need logout)")
	fmt.Println("5. Setup Autoexec (Codex only fn)")
	fmt.Println("9. Exit")
	fmt.Println("10. Setup Startup")
}

func getUserChoice() int {
	var choice int
	fmt.Scan(&choice)
	return choice
}

func handleWatchEmulator() {
	configData, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func(placeId int64) {
		defer wg.Done()
		if watchCallbackFunc != nil {
			watchCallbackFunc(placeId)
		}
	}(configData.PlaceId)
	wg.Wait()
}

func handleSetupStartup() {
	utils.ClearScreen()
	fmt.Println("1. Create shortcut")
	fmt.Println("2. Delete shortcut")
	fmt.Println("3. Go Back")
	choice := getUserChoice()
	switch choice {
	case 1:
		cmd.CreateStartup()
	case 2:
		cmd.DeleteStartup()
	case 3:
		return
	default:
		return
	}
}

func handleInvalidChoice() {
	fmt.Println("Invalid choice")
	time.Sleep(500 * time.Millisecond)
}

// Dont care this
type WatchCallback func(placeId int64)

var watchCallbackFunc WatchCallback
