package xray_core

import (
	"go.uber.org/zap"
	"os/exec"
	"strings"
	"vpngui/internal/app/command"
	"vpngui/internal/app/proxy"
	"vpngui/internal/app/repository"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

var cmd *exec.Cmd

type RunXrayCore struct {
	cr *repository.ConfigRepository
}

func NewRun(cr *repository.ConfigRepository) *RunXrayCore {
	return &RunXrayCore{
		cr: cr,
	}
}

func (x *RunXrayCore) Run() error {
	logger.Info("Starting xray-core")

	err := command.TerminateProcesses("xray-core")
	if err != nil {
		return err
	}

	cmd = exec.Command(embed.GetTempFileName("xray-core"), "run", "-c", "config/xray.json")
	err = command.RunProcess("tun2socks", cmd, x.handlerStdout, x.waitForExit)
	if err != nil {
		return err
	}

	logger.Debug("xray-api started successfully")
	return nil
}

func (x *RunXrayCore) handlerStdout(line string) {
	if strings.Contains(line, "started") {
		if err := proxy.Enable(); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}

		if err := x.cr.UpdateActiveVPN(true); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}
	} else if strings.Contains(line, "address already in use") || strings.Contains(line, "failed to listen address") {
		err := command.TerminateProcesses("xray-core")
		if err != nil {
			logger.Error("Failed to terminate xray-core processes", zap.Error(err))

			if err := proxy.Disable(); err != nil {
				logger.Error("Failed to update VPN state", zap.Error(err))
			}

			if err := x.cr.UpdateActiveVPN(false); err != nil {
				logger.Error("Failed to update VPN state", zap.Error(err))
				return
			}
		}
	}
}

func (x *RunXrayCore) waitForExit() {
	if err := cmd.Wait(); err != nil && !strings.Contains(err.Error(), "exit status 255") {
		logger.Error("xray API exited with an error", zap.Error(err))

		if err := proxy.Disable(); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
		}

		if err := x.cr.UpdateActiveVPN(false); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}
	}
}

func (x *RunXrayCore) Kill(updateActiveVPN bool) error {
	logger.Info("Stopping xray-core")

	if updateActiveVPN {
		if err := x.cr.UpdateActiveVPN(false); err != nil {
			logger.Error("Failed to update active VPN state", zap.Error(err))
		}
	}
	if err := proxy.Disable(); err != nil {
		logger.Error("Failed to update VPN state", zap.Error(err))
	}

	err := command.KillProcess("xray-core", cmd)
	if err != nil {
		logger.Error("Failed to kill xray-core", zap.Error(err))
		return err
	}

	logger.Debug("xray-core stopped successfully")
	return nil
}
