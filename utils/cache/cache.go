package cache

// *
// * This file contains the cache system for the programs
// -----------------------------------------------
// * Use the cache system to store data that is used to
// * prevent running the same adb multiple times
// * (For example, the device list)
// -----------------------------------------------
// *
import (
	"EmuWatcher/adb"
	"sync"
)

var (
	deviceCache []string
	cacheMutex  sync.Mutex
)

func GetCachedDevices() []string {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	devices := make([]string, len(deviceCache))
	copy(devices, deviceCache)
	return devices
}

func RefreshDeviceCache() error {
	cacheMutex.Lock()
	defer cacheMutex.Unlock()

	devices, err := adb.GetAllDevices()
	if err != nil {
		return err
	}
	deviceCache = devices
	return nil
}
