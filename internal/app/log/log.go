package log

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
	"vpngui/pkg/logger"
)

type Log struct{}

func New() *Log {
	return &Log{}
}

func (l *Log) CaptureStdout() {
	reader, writer, err := os.Pipe()
	if err != nil {
		logger.Warn("Error creating pipe for capturing program logs (WARNING! The 'Log' tab won't work without this)")
		return
	}

	oldStdout := os.Stdout
	os.Stdout = writer
	os.Stderr = writer

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			if line != "" {
				logger.MuLog.Lock()
				logger.OutLog = append(logger.OutLog, line)
				logger.MuLog.Unlock()
			}

			_, err = fmt.Fprintln(oldStdout, line)
			if err != nil {
				logger.Error("Error restoring output to original stdout", zap.Error(err))
				return
			}
		}
	}()
}

func (l *Log) GetLogs() string {
	logger.MuLog.Lock()
	defer logger.MuLog.Unlock()
	return strings.Join(logger.OutLog, "\n")
}
