package login

import (
	"EmuWatcher/initialize"
	"EmuWatcher/pkg/commands/login"

	"github.com/spf13/cobra"
)

func Login() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Login to emulator using cookies",
		Long: `Login to emulator using cookieg 
Put your cookies to file data/cookies.txt`,
		Run: func(cmd *cobra.Command, args []string) {
			initialize.Initialize()
			login.Login()
		},
	}
	return cmd
}
