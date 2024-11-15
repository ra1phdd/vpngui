package xray_api

import (
	"bufio"
	"go.uber.org/zap"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"vpngui/internal/app/proxy"
	"vpngui/internal/app/repository"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

var cmd *exec.Cmd

type RunXrayAPI struct {
	cr      *repository.ConfigRepository
	errorCh chan error
}

func NewRun(cr *repository.ConfigRepository) *RunXrayAPI {
	return &RunXrayAPI{
		cr:      cr,
		errorCh: make(chan error, 1),
	}
}

func (x *RunXrayAPI) Run() error {
	logger.Info("Starting xray-api")

	cmd = exec.Command(embed.GetTempFileName(), "run", "-c", "config/xray.json")
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

	logger.Debug("xray-api started successfully")
	go x.handleStdout(stdoutPipe)
	go x.waitForExit()

	return nil
}

func (x *RunXrayAPI) Kill() error {
	logger.Info("Stopping xray API")
	if cmd != nil && cmd.Process != nil {
		var err error
		if runtime.GOOS == "windows" {
			err = cmd.Process.Kill()
			if err != nil {
				logger.Error("Failed to kill process", zap.Error(err))
				return err
			}
			logger.Debug("Process killed")
		} else {
			err = cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				logger.Error("Failed to send SIGTERM to process", zap.Error(err))
				return err
			}
			logger.Debug("Sent SIGTERM signal to process")
		}
	}

	if err := proxy.Disable(); err != nil {
		logger.Error("Failed to update VPN state", zap.Error(err))
		return err
	}
	if err := x.cr.UpdateActiveVPN(false); err != nil {
		logger.Error("Failed to update active VPN state", zap.Error(err))
		return err
	}
	logger.Debug("xray-api stopped successfully")

	return nil
}

func (x *RunXrayAPI) KillOnClose() error {
	logger.Info("Stopping xray API")
	if cmd != nil && cmd.Process != nil {
		var err error
		if runtime.GOOS == "windows" {
			err = cmd.Process.Kill()
			if err != nil {
				logger.Error("Failed to kill process", zap.Error(err))
				return err
			}
			logger.Debug("Process killed")
		} else {
			err = cmd.Process.Signal(syscall.SIGTERM)
			if err != nil {
				logger.Error("Failed to send SIGTERM to process", zap.Error(err))
				return err
			}
			logger.Debug("Sent SIGTERM signal to process")
		}
	}

	if err := proxy.Disable(); err != nil {
		logger.Error("Failed to update VPN state", zap.Error(err))
		return err
	}
	logger.Debug("xray-api stopped successfully")

	return nil
}

func (x *RunXrayAPI) handleStdout(stdoutPipe io.ReadCloser) {
	logger.Info("Handling stdout for xray API")
	scanner := bufio.NewScanner(stdoutPipe)
	defer stdoutPipe.Close()

	for scanner.Scan() {
		line := scanner.Text()
		logger.Info("Received line from stdout", zap.String("line", line))
		if strings.Contains(line, "started") {
			if err := proxy.Enable(); err != nil {
				logger.Error("Failed to update VPN state", zap.Error(err))
				return
			}

			if err := x.cr.UpdateActiveVPN(true); err != nil {
				logger.Error("Failed to update VPN state", zap.Error(err))
				return
			}
		}
	}
}

func (x *RunXrayAPI) waitForExit() {
	if err := cmd.Wait(); err != nil && !strings.Contains(err.Error(), "exit status 255") {
		logger.Error("xray API exited with an error", zap.Error(err))

		if err := proxy.Disable(); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}

		if err := x.cr.UpdateActiveVPN(true); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}
	}
}
