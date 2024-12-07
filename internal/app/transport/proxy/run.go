package proxy

import (
	"go.uber.org/zap"
	"runtime"
	"vpngui/internal/app/runner"
	"vpngui/pkg/logger"
)

type Proxy struct {
	rc *runner.Command
}

func New(rc *runner.Command) *Proxy {
	return &Proxy{
		rc: rc,
	}
}

func (p *Proxy) Enable() error {
	logger.Info("Enabling proxy settings based on OS")

	var err error
	switch runtime.GOOS {
	case "darwin":
		err = p.setMacOSProxy("127.0.0.1", "2080")
	case "linux":
		err = p.setLinuxProxy("127.0.0.1", "2080")
	case "windows":
		err = p.setWindowsProxy("127.0.0.1", "2080")
	}
	if err != nil {
		logger.Error("Failed to enable proxy settings", zap.String("os", runtime.GOOS), zap.Error(err))
	} else {
		logger.Info("Proxy settings enabled successfully", zap.String("os", runtime.GOOS))
	}

	return nil
}

func (p *Proxy) Disable() error {
	logger.Info("Disabling proxy settings based on OS")

	var err error
	switch runtime.GOOS {
	case "darwin":
		err = p.clearMacOSProxy()
	case "linux":
		err = p.clearLinuxProxy()
	case "windows":
		err = p.clearWindowsProxy()
	}
	if err != nil {
		logger.Error("Failed to disable proxy settings", zap.String("os", runtime.GOOS), zap.Error(err))
	} else {
		logger.Info("Proxy settings disabled successfully", zap.String("os", runtime.GOOS))
	}

	return nil
}
