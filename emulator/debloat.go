package emulator

// Debloat.go is used to remove the bloatware from the emulator

// var bloatList = "bloat/bloatlist.txt"

// TODO:
// 1. Read the bloatlist.txt file
// 2. Get the device list from the emulator (must be running) using Emulator/adb package (GetAllDevices []string)
// 3. Loop through the device list and remove the bloatware using adb uninstall
// * adb shell pm uninstall <package_name>
