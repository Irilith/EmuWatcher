package emulator

import (
	"EmuWatcher/utils"
	"EmuWatcher/utils/config"
	"fmt"
	"os/exec"
	"time"
)

func OpenEmulator() {
	emuList := utils.RemoveEmpty(GetEmulatorList())
	for _, emuName := range emuList {
		isRunning, err := IsRunning(emuName)
		if err != nil {
			fmt.Println("Error checking if emulator is running:", err)
		}
		if !isRunning {
			err := LaunchEmulator(emuName)
			if err != nil {
				fmt.Println("Error launching emulator:", err)
			}
			time.Sleep(7 * time.Second)
			continue
		}
	}
}

func LaunchEmulator(emuName string) error {
	emuPath, err := config.GetConfig("LDPlayerPath")
	if err != nil {
		return err
	}

	if emuPath == "" {
		return fmt.Errorf("Please set the LDPlayer path in the config")
	}
	emuPath = emuPath + "\\dnconsole.exe"
	args := []string{
		"launch",
		"--name",
		emuName,
	}

	cmd := exec.Command(emuPath, args...)

	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func GetEmulatorList() []string {
	emuPath, err := config.GetConfig("LDPlayerPath")
	if err != nil {
		fmt.Println("Error getting config:", err)
	}
	if emuPath == "" {
		fmt.Println("Please set the LDPlayer path in the config")
	}
	emuPath = emuPath + "\\dnconsole.exe"
	args := []string{
		"list", // only get the name
	}
	out, err := exec.Command(emuPath, args...).Output()
	if err != nil {
		fmt.Println("Error getting emulator list:", err)
	}

	return utils.SplitLines(string(out))
}

func IsRunning(emuName string) (bool, error) {
	emuPath, err := config.GetConfig("LDPlayerPath")
	if err != nil {
		fmt.Println("Error getting config:", err)
	}
	if emuPath == "" {
		return false, fmt.Errorf("Please set the LDPlayer path in the config")
	}
	emuPath = emuPath + "\\dnconsole.exe"
	args := []string{
		"isrunning",
		"--name",
		emuName,
	}
	out, err := exec.Command(emuPath, args...).Output()
	if err != nil {
		return false, err
	}

	return utils.Expect(string(out), "Running"), nil
}
