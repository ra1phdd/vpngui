package xray_core

import (
	"go.uber.org/zap"
	"os/exec"
	"strings"
	"vpngui/internal/app/repository"
	"vpngui/internal/app/runner"
	"vpngui/internal/app/transport"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

var cmd *exec.Cmd

type RunXrayCore struct {
	cr *repository.ConfigRepository
	rp *runner.Process
	ts *transport.Transport
}

func NewRun(cr *repository.ConfigRepository, rp *runner.Process, ts *transport.Transport) *RunXrayCore {
	return &RunXrayCore{
		cr: cr,
		rp: rp,
		ts: ts,
	}
}

func (rx *RunXrayCore) Run() error {
	logger.Info("Starting xray-core")

	err := rx.rp.Terminate("xray-core")
	if err != nil {
		return err
	}

	cmd = exec.Command(embed.GetTempFileName("xray-core"), "run", "-c", "config/xray.json")
	err = rx.rp.Run("tun2socks", cmd, rx.handlerStdout, rx.waitForExit)
	if err != nil {
		return err
	}

	logger.Debug("xray-api started successfully")
	return nil
}

func (rx *RunXrayCore) handlerStdout(line string) {
	if strings.Contains(line, "started") {
		if err := rx.ts.Enable(); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}

		if err := rx.cr.UpdateActiveVPN(true); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}
	} else if strings.Contains(line, "address already in use") || strings.Contains(line, "failed to listen address") {
		err := rx.rp.Terminate("xray-core")
		if err != nil {
			logger.Error("Failed to terminate xray-core processes", zap.Error(err))

			if err := rx.ts.Disable(); err != nil {
				logger.Error("Failed to update VPN state", zap.Error(err))
				return
			}

			if err := rx.cr.UpdateActiveVPN(false); err != nil {
				logger.Error("Failed to update VPN state", zap.Error(err))
				return
			}
		}
	}
}

func (rx *RunXrayCore) waitForExit() {
	if err := cmd.Wait(); err != nil && !strings.Contains(err.Error(), "exit status 255") {
		logger.Error("xray API exited with an error", zap.Error(err))

		if err := rx.ts.Disable(); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}

		if err := rx.cr.UpdateActiveVPN(false); err != nil {
			logger.Error("Failed to update VPN state", zap.Error(err))
			return
		}
	}
}

func (rx *RunXrayCore) Kill(updateActiveVPN bool) error {
	logger.Info("Stopping xray-core")

	if updateActiveVPN {
		if err := rx.cr.UpdateActiveVPN(false); err != nil {
			logger.Error("Failed to update active VPN state", zap.Error(err))
		}
	}
	if err := rx.ts.Disable(); err != nil {
		logger.Error("Failed to update VPN state", zap.Error(err))
	}

	err := rx.rp.Kill("xray-core", cmd)
	if err != nil {
		logger.Error("Failed to kill xray-core", zap.Error(err))
		return err
	}

	logger.Debug("xray-core stopped successfully")
	return nil
}
