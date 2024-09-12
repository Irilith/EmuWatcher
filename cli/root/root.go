package root

import (
	loginCmd "EmuWatcher/cli/login"
	versionCmd "EmuWatcher/cli/version"
	watchCmd "EmuWatcher/cli/watch"
	"EmuWatcher/utils/version"

	"github.com/spf13/cobra"
)

func Root() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "EmuWatcher",
		Short:   "EmuWatcher is a tool to watch emulators",
		Version: version.GetVersion(),
	}
	cmd.AddCommand(versionCmd.Version())
	cmd.AddCommand(versionCmd.Commit())
	cmd.AddCommand(watchCmd.Watch())
	cmd.AddCommand(loginCmd.Login())
	return cmd
}
