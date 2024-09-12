package emulator

import (
	"EmuWatcher/adb"
	"EmuWatcher/utils/cache"
)

func SetupAutoExec() {
	devices := cache.GetCachedDevices()
	for _, device := range devices {
		args := []string{
			"-s",
			device,
			"shell",
			"rm",
			"-r",
			"/sdcard/Codex/Autoexec",
		}
		adb.ExecuteADBCommand(args)
		args = []string{
			"-s",
			device,
			"push",
			".\\data\\autoexec",
			"/sdcard/Codex/Autoexec",
		}
		adb.ExecuteADBCommand(args)
	}
}
