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

func CheckOS() string {
	hostInfo, _ := host.Info()
	return hostInfo.OS
}
func CheckArch() string {
	hostInfo, _ := host.Info()
	return hostInfo.KernelArch
}

//Check platform
func CheckPlatform() string {
	hostInfo, _ := host.Info()
	return hostInfo.Platform + "," + hostInfo.KernelVersion
}

//Check whether program r is running in root
func CheckIsRoot() bool {
	uid := os.Getuid()
	return uid == 0
}

//Check if theres enough space to run GoOwl;Default 2G
func CheckDiskSpace() bool {
	pwd, _ := os.Getwd()
	diskInfo, _ := disk.Usage(pwd)
	return float64(diskInfo.Free)/1024/1024/1024 > 2
}

//Check if memory>1G
func CheckMemory() bool {
	memoryInfo, _ := mem.VirtualMemory()
	return float64(memoryInfo.Total)/1024/1024/1024 > 1
}

//Return -1 if unknown;0 if fasle andf 1 if true
func CheckDocker() int {
	isDockerEnv := false
	if _, err := os.Stat("/.dockerenv"); err == nil {
		isDockerEnv = true
	}
	isDockerEnv = false
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

//Check address is reachable
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
