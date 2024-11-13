package log

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"os"
	"strings"
	"sync"
	"vpngui/pkg/logger"
)

type Log struct{}

func New() *Log {
	return &Log{}
}

var (
	OutLog []string
	mu     sync.Mutex
)

func (a *Log) CaptureStdout() {
	reader, writer, err := os.Pipe()
	if err != nil {
		logger.Warn("Error creating pipe for capturing program logs (WARNING! The 'Log' tab won't work without this)")
		return
	}

	oldStdout := os.Stdout
	os.Stdout = writer

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			mu.Lock()
			OutLog = append(OutLog, line)
			mu.Unlock()

			_, err = fmt.Fprintln(oldStdout, line)
			if err != nil {
				logger.Error("Error restoring output to original stdout", zap.Error(err))
				return
			}
		}
	}()
}

func (a *Log) GetLogs() string {
	mu.Lock()
	defer mu.Unlock()
	return strings.Join(OutLog, "\n")
}
