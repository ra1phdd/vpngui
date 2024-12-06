package command

import (
	"bufio"
	"fmt"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"vpngui/pkg/logger"
)

func RunProcess(name string, cmd *exec.Cmd, handler func(line string), waitForExit func()) error {
	logger.Info("Starting process", zap.String("name", name))

	cmd.SysProcAttr = GetSysProcAttr()
	cmd.Stderr = os.Stderr

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("Failed to get stdout pipe", zap.Error(err))
		return err
	}

	if err := cmd.Start(); err != nil {
		logger.Error("Failed to start command", zap.Error(err))
		return err
	}

	go handleStdout(name, stdoutPipe, handler)
	go waitForExit()

	logger.Debug("tun2socks started successfully")
	return nil
}

func handleStdout(name string, stdoutPipe io.ReadCloser, handler func(line string)) {
	logger.Info(fmt.Sprintf("Handling stdout for %s", name))

	scanner := bufio.NewScanner(stdoutPipe)
	defer stdoutPipe.Close()

	for scanner.Scan() {
		line := scanner.Text()
		logger.Info(fmt.Sprintf("Received line from %s", name), zap.String("line", line))
		handler(line)
	}
}

func KillProcess(name string, cmd *exec.Cmd) error {
	logger.Info("Stopping process", zap.String("name", name))

	if cmd != nil && cmd.Process != nil {
		var err error
		if runtime.GOOS == "windows" {
			err = cmd.Process.Kill()
		} else {
			err = cmd.Process.Signal(syscall.SIGTERM)
		}
		if err != nil {
			logger.Error("Failed to kill process", zap.Error(err))
			return err
		}

		logger.Debug("Process killed")
	}

	logger.Debug("Process stopped successfully")
	return nil
}

func TerminateProcesses(name string) error {
	logger.Info("Terminate processes", zap.String("name", name))

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/FI", fmt.Sprintf("IMAGENAME eq %s*", name), "/F")
	} else {
		cmd = exec.Command("pkill", "-f", name)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Failed to terminate processes", zap.Error(err), zap.String("output", string(output)))
		return err
	}

	logger.Info("Successfully terminated processes")
	return nil
}

func RunCommands(commands [][]string, ignoreErr bool) error {
	for _, args := range commands {
		logger.Debug("Executing command", zap.String("cmd", strings.Join(args, " ")))
		cmd := exec.Command(args[0], args[1:]...)
		err := RunCommand(cmd, ignoreErr)
		if err != nil {
			return err
		}
	}

	return nil
}

func RunCommand(cmd *exec.Cmd, ignoreErr bool) error {
	cmd.SysProcAttr = GetSysProcAttr()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil && !ignoreErr {
		logger.Error("Command execution failed", zap.String("cmd", cmd.String()), zap.Error(err))
		return err
	}

	return nil
}
