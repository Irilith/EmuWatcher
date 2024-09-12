package watch

import (
	"EmuWatcher/emulator"
	"EmuWatcher/initialize"
	"EmuWatcher/utils"
	"EmuWatcher/utils/config"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func Watch() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "watch",
		Short: "Start to watch Emulators",
		Run: func(cmd *cobra.Command, args []string) {
			initialize.Initialize()
			utils.ClearScreen()
			configData, err := config.LoadConfig()
			if err != nil {
				fmt.Println("Error loading config:", err)
				return
			}
			go emulator.OpenEmulator()
			emulator.Watch(configData.PlaceId, func() {
				os.Exit(0)
			})
		},
	}
	return cmd
}
