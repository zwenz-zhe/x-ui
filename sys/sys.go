package sys

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strings"
)

// GetTCPCount 返回当前活动的 TCP 连接数
func GetTCPCount() (int, error) {
	// 实现获取 TCP 连接数的逻辑
}

// GetUDPCount 返回当前活动的 UDP 连接数
func GetUDPCount() (int, error) {
	switch runtime.GOOS {
	case "windows":
		return getWindowsUDPCount()
	case "linux":
		return getLinuxUDPCount()
	default:
		return 0, fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
}

// getWindowsUDPCount Windows 系统下获取 UDP 连接数
func getWindowsUDPCount() (int, error) {
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	// 解析 netstat 输出，统计 UDP 连接数
	count := 0
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 5 && fields[1] == "UDP" {
			count++
		}
	}

	return count, nil
}

// getLinuxUDPCount Linux 系统下获取 UDP 连接数
func getLinuxUDPCount() (int, error) {
	// 使用 net 包监听所有 UDP 连接，然后统计连接数
	conns, err := net.ListenPacket("udp", "0.0.0.0:0")
	if err != nil {
		return 0, err
	}
	defer conns.Close()

	count := 0
	buffer := make([]byte, 4096)
	for {
		_, _, err := conns.ReadFrom(buffer)
		if err != nil {
			// 在没有新连接时返回统计结果
			if opErr, ok := err.(*net.OpError); ok && opErr.Op == "read" {
				break
			}
			return 0, err
		}
		count++
	}

	return count, nil
}
