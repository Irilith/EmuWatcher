package login

// This come later, i dont intend to move all to pkg (LAZY AF)
// this project folder structure like a dish of spaghetti lol
import (
	"EmuWatcher/adb"
	"EmuWatcher/emulator"
	"EmuWatcher/roblox"
	"EmuWatcher/utils"
	"fmt"
	"sync"
	"time"
)

func Login() {
	cookies, err := roblox.GetCookies()
	cookiesCount := len(cookies)
	if cookiesCount == 0 {
		fmt.Println("No cookies found")
		return
	}
	devices, err := adb.GetAllDevices() // Not use cached this one because we need to get the latest devices list
	if err != nil {
		fmt.Println("Error getting devices:", err)
		return
	}
	if len(devices) == 0 {
		fmt.Println("No devices found")
		return
	}

	if len(devices) < cookiesCount {
		fmt.Println("Not enough devices to login")
		return
	}
	var wg sync.WaitGroup

	for _, device := range devices {
		wg.Add(1)
		go func(device string) {
			fmt.Println("Rooting in device:", device)
			defer wg.Done()
			adb.Root(device)
		}(device)
	}

	wg.Wait()
	time.Sleep(7 * time.Second)
	utils.ClearScreen()
	for i, cookie := range cookies {
		wg.Add(1)
		go func(device string, cookie string, index int) {
			defer wg.Done()
			status, err := emulator.Login(device, cookie)
			if err != nil {
				fmt.Println("Error logging in:", err)
			}
			go func(status bool, device string) {
				if status {
					time.Sleep(5 * time.Second)
					adb.ForceStartGame(device)
				}
			}(status, device)
		}(devices[i], cookie, i)
	}
	wg.Wait()
}
