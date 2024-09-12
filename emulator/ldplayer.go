package emulator

import (
	"EmuWatcher/utils/config"
	"fmt"
	"os/exec"
)

func ArrangeEmulators() {
	emuPath, err := config.GetConfig("LDPlayerPath")
	if err != nil {
		fmt.Println("Error getting config:", err)
	}
	if emuPath == "" {
		fmt.Println("Please set the LDPlayer path in the config")
	}
	emuPath = emuPath + "\\dnconsole.exe"
	args := []string{
		"sortWnd",
	}
	err = exec.Command(emuPath, args...).Run()
	if err != nil {
		fmt.Println("Error arranging emulators:", err)
	}
}
