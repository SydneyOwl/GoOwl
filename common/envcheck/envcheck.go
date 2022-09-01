package envcheck

import (
	"net"
	"os"
	"strings"
	"time"

	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// CheckOS returns OS info of the host.
func CheckOS() string {
	hostInfo, _ := host.Info()
	return hostInfo.OS
}

// CheckArch returns Arch info of the host.
func CheckArch() string {
	hostInfo, _ := host.Info()
	return hostInfo.KernelArch
}

// CheckPlatform returns platform info and kernel version of the host.
func CheckPlatform() string {
	hostInfo, _ := host.Info()
	return hostInfo.Platform + "," + hostInfo.KernelVersion
}

// CheckIsRoot check whether program running with root permission.
func CheckIsRoot() bool {
	uid := os.Getuid()
	return uid == 0
}

// CheckDiskSpace check if there's enough space to run GoOwl; Minimum is 2G.
func CheckDiskSpace() bool {
	pwd, _ := os.Getwd()
	diskInfo, _ := disk.Usage(pwd)
	return float64(diskInfo.Free)/1024/1024/1024 > 2
}

// CheckMemory check if there's enough memory to run GoOwl; Minimum is 128M.
func CheckMemory() bool {
	memoryInfo, _ := mem.VirtualMemory()
	return float64(memoryInfo.Total)/1024/1024 > 128
}

// CheckDocker check if GoOwl running in docker environment; Minimum is 128M. Return -1 if unknown;0 if false and 1 if true
func CheckDocker() int {
	isDockerEnv := false
	if _, err := os.Stat("/.dockerenv"); err == nil {
		isDockerEnv = true
	}
	cgroup, err := os.ReadFile("/proc/self/cgroup")
	if err != nil {
		return -1
	}
	isDockerCGroup := strings.Contains(string(cgroup), "docker")
	if isDockerCGroup || isDockerEnv {
		return 1
	}
	return 0
}

// CheckConn check if address is reachable.
func CheckConn(addr string) bool {
	conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
	if err != nil {
		return false
	} else {
		if conn != nil {
			return true
		} else {
			return false
		}
	}
}
