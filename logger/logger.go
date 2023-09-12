package logger

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/op/go-logging"
)

var logger *logging.Logger

func init() {
	InitLogger(logging.INFO)
}

func InitLogger(level logging.Level) {
	format := logging.MustStringFormatter(
		`%{time:2006/01/02 15:04:05} %{level} - %{message}`,
	)
	newLogger := logging.MustGetLogger("x-ui")
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backendFormatter)
	backendLeveled.SetLevel(level, "")
	newLogger.SetBackend(backendLeveled)

	logger = newLogger
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Warning(args ...interface{}) {
	logger.Warning(args...)
}

func Warningf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func QueryUserTraffic(userID string) {
	// Generate tcpdump filter expression for user traffic
	filterExpr := fmt.Sprintf("host %s", userID)

	// Run tcpdump command with filter expression
	cmd := exec.Command("tcpdump", "-i", "eth0", filterExpr)
	output, err := cmd.Output()
	if err != nil {
		logger.Errorf("Failed to run tcpdump command: %s", err)
		return
	}

	// Process the tcpdump output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		// Log the tcpdump output
		logger.Infof("TCPDUMP: %s", line)
	}
}

func Trace(args ...interface{}) {
	logger.Debug(args...)
}

func Tracef(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}
