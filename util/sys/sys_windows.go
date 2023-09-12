package sys

import (
	"os/exec"
	"strings"
)

func GetTCPCount() (int, error) {
	output, err := exec.Command("cmd", "/C", "netstat", "-n", "-p", "tcp").Output()
	if err != nil {
		return 0, err
	}

	// Count the number of TCP connections
	count := 0
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "TCP") {
			count++
		}
	}

	return count, nil
}

func GetUDPCount() (int, error) {
	output, err := exec.Command("cmd", "/C", "netstat", "-n", "-p", "udp").Output()
	if err != nil {
		return 0, err
	}

	// Count the number of UDP connections
	count := 0
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "UDP") {
			count++
		}
	}

	return count, nil
}
