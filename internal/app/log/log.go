package log

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Log struct{}

func New() *Log {
	return &Log{}
}

var OutLog []string

func (a *Log) CaptureStdout() {
	reader, writer, err := os.Pipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка создания пайпа: %v\n", err)
		return
	}

	oldStdout := os.Stdout
	os.Stdout = writer

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()
			OutLog = append(OutLog, line)
			fmt.Fprintln(oldStdout, line)
		}
	}()
}

func (a *Log) GetLogs() string {
	return strings.Join(OutLog, "\n")
}
