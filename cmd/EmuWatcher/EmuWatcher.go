package EmuWatcher

import (
	"EmuWatcher/initialize"
	"EmuWatcher/pkg/update"
	"EmuWatcher/utils/menu"
)

func Start() {
	initialize.Initialize()
	update.CheckForUpdate()
	menu.RunMenu()
}
