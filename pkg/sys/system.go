package sys

import (
	"fmt"
	"log"
	"runtime"
	"testing"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemInfo struct {
	CoreCount      int     `json:"core_count"`
	ThreadCount    int     `json:"thread_count"`
	CPUUtilization float64 `json:"cpu_utilization"`
	TotalRAM       uint64  `json:"total_ram"`
	AvailableRAM   uint64  `json:"available_ram"`
	TotalDisk      uint64  `json:"total_disk"`
	FreeDisk       uint64  `json:"free_disk"`
}

func bytesToGB(bytes uint64) float64 {
	return float64(bytes) / (1024 * 1024 * 1024)
}

func GetSystemInfo() (SystemInfo, error) {
	var info SystemInfo

	info.CoreCount = runtime.NumCPU()
	time.Sleep(60 * time.Second)
	percentages, err := cpu.Percent(0, false)
	if err != nil {
		return info, fmt.Errorf("error fetching CPU utilization: %w", err)
	}
	info.CPUUtilization = percentages[0]

	info.ThreadCount = runtime.NumGoroutine()

	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return info, fmt.Errorf("error fetching RAM info: %w", err)
	}
	info.TotalRAM = uint64(bytesToGB(memInfo.Total))
	info.AvailableRAM = uint64(bytesToGB(memInfo.Available))

	diskInfo, err := disk.Usage("/")
	if err != nil {
		return info, fmt.Errorf("error fetching disk info: %w", err)
	}
	info.TotalDisk = uint64(bytesToGB(diskInfo.Total))
	info.FreeDisk = uint64(bytesToGB(diskInfo.Free))

	return info, nil
}

func TestGetCPUInfo(t *testing.T) {
	systemInfo, err := GetSystemInfo()
	if err != nil {
		log.Fatalf("Failed to get system info: %v", err)
	}

	fmt.Printf("System Information:\n")
	fmt.Printf("Core Count: %d\n", systemInfo.CoreCount)
	fmt.Printf("Thread Count: %d\n", systemInfo.ThreadCount)
	fmt.Printf("CPU Utilization: %.2f%%\n", systemInfo.CPUUtilization)
	fmt.Printf("Total RAM: %d Gb\n", systemInfo.TotalRAM)
	fmt.Printf("Available RAM: %d Gb\n", systemInfo.AvailableRAM)
	fmt.Printf("Total Disk: %d Gb\n", systemInfo.TotalDisk)
	fmt.Printf("Free Disk: %d Gb\n", systemInfo.FreeDisk)
}
