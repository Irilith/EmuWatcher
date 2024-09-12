// LICENSE: GNU GPL-3.0
// LICENSE URL: https://raw.githubusercontent.com/Irilith/EmuWatcher/main/LICENSE

// DISCLAIMER:
// This program may reference or interact with third-party executors.
// The creator of this program did not create or have any involvement in the development of these executors.
// The use of this program, including any associated executors, is entirely at the user's own risk and discretion.
// The creator assumes no responsibility for any misuse, illegal activities, or consequences resulting from the use of this program or its associated executors.
// It is the user's sole responsibility to ensure that their use complies with all relevant laws and regulations.

package main

import (
	rootCmd "EmuWatcher/cli/root"
	"EmuWatcher/cmd/EmuWatcher"
	"fmt"
	"os"
)

func hasArguments() bool {
	return len(os.Args) > 1
}

// TODO:
// Add more emulators
// Auto open emulator if closed
func main() {
	if hasArguments() {
		rootCmd := rootCmd.Root()
		if err := rootCmd.Execute(); err != nil {
			fmt.Println("Error executing command:", err)
			os.Exit(1)
		}
	} else {
		EmuWatcher.Start()
	}
}
