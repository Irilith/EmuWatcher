package verison

import (
	"EmuWatcher/utils/version"
	"fmt"

	"github.com/spf13/cobra"
)

func Version() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "version",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.GetVersion())
		},
	}

	return cmd
}

func Commit() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "commit",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(version.GetCommit())
		},
	}
	return cmd
}
