package emulator

import (
	"EmuWatcher/adb"
	"EmuWatcher/ocr"
	"EmuWatcher/pkg/webhook"
	"EmuWatcher/utils/cache"
	"EmuWatcher/utils/config"
	"EmuWatcher/utils/imghandler"
	"EmuWatcher/utils/log"
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func getDevices() []string {
	devices := cache.GetCachedDevices()
	return devices
}

func Watch(placeId int64, onExit func()) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan struct{})
	if status, err := config.GetConfig("Run_Sequential"); err == nil && status == "true" {
		runSequential(placeId, done)
	} else {
		runParallel(placeId, done, watchRoblox)
		runParallel(placeId, done, checkCrashes)
		runParallel(placeId, done, checkStuckLoading)
	}

	if status, err := config.GetConfig("Watch_Using_Script"); err == nil && status == "true" {
		runParallel(placeId, done, checkOnline)
	}
	if status, err := config.GetConfig("Auto_Open_Ld"); err == nil && status == "true" {
		go func() {
			ticker := time.NewTicker(20 * time.Second)
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					OpenEmulator()
				}
			}
		}()
	}
	if status, err := config.GetConfig("Auto_Send_Webhook"); err == nil && status == "true" {
		url, err := config.GetConfig("Webhook_URL")
		if err != nil {
			fmt.Println("Error getting config:", err)
			return
		}
		go func() {
			ticker := time.NewTicker(1 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					webhookClient := webhook.NewWebhook(url)
					err := webhookClient.SendWebhook()
					if err != nil {
						fmt.Println("Error sending webhook:", err)
					}
				}
			}
		}()
	}
	if status, err := config.GetConfig("Auto_Arrange"); err == nil && status == "true" {
		go func() {
			ticker := time.NewTicker(35 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					ArrangeEmulators()
				}
			}
		}()
	}

	go restartAdb(placeId, done)

	<-sigs

	close(done)
	if onExit != nil {
		onExit()
	}
}

func runSequential(placeId int64, done chan struct{}) {
	go func() {
		config, err := config.GetConfig("Sequential_Interval")
		if err != nil {
			fmt.Println("Error getting config:", err)
			return
		}
		interval, err := strconv.Atoi(config)
		if err != nil {
			fmt.Println("Error parsing interval:", err)
			return
		}
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				watchRoblox(placeId)
				time.Sleep(4 * time.Second)
				checkCrashes(placeId)
				time.Sleep(4 * time.Second)
				checkStuckLoading(placeId)
			}
		}
	}()
}

func runParallel(placeId int64, done chan struct{}, task func(int64)) {
	go func() {
		config, err := config.GetConfig("Watch_Interval")
		if err != nil {
			fmt.Println("Error getting config:", err)
			return
		}
		interval, err := strconv.Atoi(config)
		if err != nil {
			fmt.Println("Error parsing interval:", err)
			return
		}
		ticker := time.NewTicker(time.Duration(interval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				task(placeId)
			}
		}
	}()
	return
}

func startRealTimeWatch(packageName string, placeId int64) {
	deviceWatchs := make(map[string]chan struct{})
	devices := getDevices()
	for _, device := range devices {
		device = strings.TrimSpace(device)
		if _, exists := deviceWatchs[device]; !exists {
			log.Greenf("Starting Real Time watcher for device: %s", device)
			stopCh := make(chan struct{})
			deviceWatchs[device] = stopCh
			go adb.RealTimeWatch(device, packageName, placeId, stopCh)
		}
	}
	for device, stopCh := range deviceWatchs {
		if !containsDevice(devices, device) {
			log.Yellowf("Stopping Real Time watcher for device: %s\n", device)
			close(stopCh)
			delete(deviceWatchs, device)
		}
	}
}

func restartAdb(placeId int64, done chan struct{}) {
	ticker := time.NewTicker(120 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			fmt.Println("Refreshing adb server...")
			adb.RestartAdb()
			err := cache.RefreshDeviceCache()
			if err != nil {
				fmt.Println("Error refreshing device cache:", err)
			}
			status, err := config.GetConfig("Real_Time_Watch")
			if err != nil {
				fmt.Println("Error getting config:", err)
				return
			}
			if status == "true" {
				go startRealTimeWatch("com.roblox.client", placeId)
			}
		}
	}
}

func checkOnline(placeId int64) {
	for _, device := range getDevices() {
		fmt.Println("Checking online status for device:", device)
		device = strings.TrimSpace(device)
		args := []string{
			"-s",
			device,
			"shell",
			"cat",
			"sdcard/Codex/Workspace/Emu.Watcher",
		}
		out, err := adb.ExecuteADBCommand(args)
		if err != nil {
			fmt.Println("Error checking raw status:", err)
			// If error, restart the game, because it's probably the game had not yet create a file (which i will fix later on)
			adb.StopGame(device)
			adb.ForceStartGame(device)
			time.Sleep(5 * time.Second)
			adb.JoinInstance(device, placeId)
			return
		}
		outString := string(out)
		outSplit := strings.Split(outString, "|")

		if len(outSplit) < 2 {
			fmt.Println("Error, not enough output")
			adb.StopGame(device)
			adb.ForceStartGame(device)
			time.Sleep(5 * time.Second)
			adb.JoinInstance(device, placeId)
			return
		}

		currentUnix := time.Now().Unix()
		emuUnix, err := strconv.ParseInt(outSplit[1], 10, 64)
		if err != nil {
			fmt.Println("Error parsing unix time:", err)
			return
		}

		if currentUnix >= emuUnix+60 {
			fmt.Printf("Device %s, User: %s is offline\n", device, outSplit[0])
			adb.StopGame(device)
			adb.ForceStartGame(device)
			time.Sleep(5 * time.Second)
			adb.JoinInstance(device, placeId)
		} else {
			fmt.Printf("Device %s, User: %s is online\n", device, outSplit[0])
		}
		time.Sleep(3 * time.Second)
	}
	return
}

func watchRoblox(placeId int64) {
	for _, device := range getDevices() {
		device = strings.TrimSpace(device)
		fmt.Println("Checking Roblox status for device:", device)
		if !adb.IsRobloxRunning(device) {
			log.Yellow("Roblox is not running")
			adb.StartGame(device)
			time.Sleep(4 * time.Second)
			adb.JoinInstance(device, placeId)
		}
		time.Sleep(3 * time.Second)
	}
	return
}

func checkCrashes(placeId int64) {
	for _, device := range getDevices() {
		fmt.Println("Checking crash for device:", device)
		device = strings.TrimSpace(device)
		if status, err := ocr.CheckCrash(device); err != nil {
			fmt.Println("Error checking crash status:", err)
			return
		} else if status {
			log.Yellow("Roblox is crashing")
			adb.StopGame(device)
			adb.ForceStartGame(device)
			time.Sleep(4 * time.Second)
			adb.JoinInstance(device, placeId)
		}
		time.Sleep(3 * time.Second)
	}
	return
}

func checkStuckLoading(placeId int64) {
	for _, device := range getDevices() {
		device = strings.TrimSpace(device)
		img, err := adb.CaptureScreen(device)
		if err != nil {
			fmt.Println("Error capturing screen:", err)
			return
		}
		bytesReader := bytes.NewBuffer(img)
		if status, err := imghandler.DetectColorRange(bytesReader); err == nil {
			if status >= 98 {
				log.Yellow("Roblox is stuck loading")
				adb.StopGame(device)
				adb.StartGame(device)
				time.Sleep(4 * time.Second)
				adb.JoinInstance(device, placeId)
			}
		} else {
			fmt.Println("Error detecting image:", err)
			return
		}
		time.Sleep(20 * time.Second)
	}
	return
}

func containsDevice(devices []string, device string) bool {
	for _, d := range devices {
		if d == device {
			return true
		}
	}
	return false
}
