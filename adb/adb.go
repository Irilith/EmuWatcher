package adb

import (
	"EmuWatcher/utils/log"
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

const (
	RobloxPackage = "com.roblox.client"
	adb           = "./tools/adb/adb"
)

func ExecuteADBCommand(args []string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, adb, args...)
	return cmd.CombinedOutput()
}

func ClearLogBuffer(emuName string) error {
	_, err := ExecuteADBCommand([]string{"-s", emuName, "logcat", "-c"})
	return err
}

func RealTimeWatch(emuName, packageName string, placeId int64, done <-chan struct{}) {
	if err := ClearLogBuffer(emuName); err != nil {
		log.Redf("Error clearing log buffer for device %s: %v", emuName, err)
		return
	}

	cmd := exec.Command(adb, "-s", emuName, "logcat", "-s", "ActivityManager")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Redf("Error creating StdoutPipe for device %s: %v", emuName, err)
		return
	}

	if err := cmd.Start(); err != nil {
		log.Redf("Error starting logcat command for device %s: %v", emuName, err)
		return
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		select {
		case <-done:
			log.Yellowf("Stopping watcher for device %s", emuName)
			if err := cmd.Process.Kill(); err != nil {
				log.Redf("Error killing process for device %s: %v", emuName, err)
			}
			return
		default:
			line := scanner.Text()
			if (strings.Contains(line, "ANR") || strings.Contains(line, "Killing")) && strings.Contains(line, packageName) {
				log.Bluef("Detected unusual behaviour Roblox on device %s, Rejoin the game", emuName)
				StopGame(emuName)
				ForceStartGame(emuName)
				time.Sleep(5 * time.Second)
				JoinInstance(emuName, placeId)
				if err := ClearLogBuffer(emuName); err != nil {
					log.Redf("Error clearing log buffer for device %s: %v", emuName, err)
				}
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Redf("Error reading logcat output for device %s: %v", emuName, err)
	}
	log.Yellowf("Watcher ended for device %s", emuName)
}

func IsRobloxRunning(emuName string) bool {
	output, err := ExecuteADBCommand([]string{"-s", emuName, "shell", "pidof", RobloxPackage})
	return err == nil && len(output) > 1
}

func Root(emuName string) error {
	_, err := ExecuteADBCommand([]string{"-s", emuName, "root"})
	if err != nil {
		return fmt.Errorf("failed to root device: %w", err)
	}
	return nil
}

func RestartAdb() error {
	if err := exec.Command(adb, "kill-server").Run(); err != nil {
		return fmt.Errorf("failed to kill adb server: %w", err)
	}
	if err := exec.Command(adb, "start-server").Run(); err != nil {
		return fmt.Errorf("failed to start adb server: %w", err)
	}
	return nil
}

func ForceStartGame(emuName string) error {
	_, err := ExecuteADBCommand([]string{"-s", emuName, "shell", "am", "start", RobloxPackage})
	if err != nil {
		return fmt.Errorf("failed to start Roblox: %w", err)
	}
	return nil
}

func StopGame(emuName string) error {
	log.Greenf("Stopping Roblox In Emu %s...", emuName)
	if IsRobloxRunning(emuName) {
		_, err := ExecuteADBCommand([]string{"-s", emuName, "shell", "am", "force-stop", RobloxPackage})
		if err != nil {
			return fmt.Errorf("failed to stop Roblox: %w", err)
		}
	}
	return nil
}

func StartGame(emuName string) error {
	if !IsRobloxRunning(emuName) {
		log.Greenf("Starting Roblox In Emu %s...", emuName)
		return ForceStartGame(emuName)
	}
	return nil
}

func JoinInstance(emuName string, placeId int64) error {
	log.Greenf("Joining Place: %d, Devices: %s", placeId, emuName)
	_, err := ExecuteADBCommand([]string{
		"-s", emuName, "shell", "am", "start",
		"-a", "android.intent.action.VIEW",
		"-d", fmt.Sprintf("roblox://placeId=%d", placeId),
	})
	if err != nil {
		return fmt.Errorf("failed to join instance: %w", err)
	}
	return nil
}

func CaptureScreen(emuName string) ([]byte, error) {
	output, err := ExecuteADBCommand([]string{"-s", emuName, "exec-out", "screencap", "-p"})
	if err != nil {
		return nil, fmt.Errorf("failed to perform screencap: %w", err)
	}
	return output, nil
}

func GetAllDevices() ([]string, error) {
	output, err := ExecuteADBCommand([]string{"devices"})
	if err != nil {
		return nil, fmt.Errorf("failed to perform devices command: %w", err)
	}

	var devices []string
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, line := range lines[1:] { // Skip the first line (header)
		fields := strings.Fields(line)
		if len(fields) > 0 {
			devices = append(devices, fields[0])
		}
	}

	return devices, nil
}
