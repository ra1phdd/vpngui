package tun

import (
	"go.uber.org/zap"
	"os/exec"
	"runtime"
	"strings"
	"vpngui/internal/app/runner"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

var cmd *exec.Cmd

type Tun struct {
	rc *runner.Command
	rp *runner.Process
}

func New(rc *runner.Command, rp *runner.Process) *Tun {
	return &Tun{
		rc: rc,
		rp: rp,
	}
}

func (t *Tun) Enable() error {
	logger.Info("Enabling TUN settings based on OS")

	if runtime.GOOS == "linux" {
		commands := [][]string{
			{"ip", "tuntap", "add", "mode", "tun", "dev", "tun0"},
			{"ip", "addr", "add", "198.18.0.1/15", "dev", "tun0"},
			{"ip", "link", "set", "dev", "tun0", "up"},
		}

		err := t.rc.RunCommands(commands, false)
		if err != nil {
			return err
		}
	} else if runtime.GOOS == "windows" {
		commands := [][]string{
			{"netsh", "interface", "set", "interface", "name=\"wintun\"", "admin=disabled"},
		}
		err := t.rc.RunCommands(commands, true)
		if err != nil {
			return err
		}
	}

	err := t.RunTun2socks()
	if err != nil {
		return err
	}

	return nil
}

func (t *Tun) Disable() error {
	logger.Info("Disabling TUN settings based on OS")

	err := t.KillTun2socks()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "darwin":
		err = t.clearMacOSTun()
	case "linux":
		err = t.clearLinuxTun()
	case "windows":
		err = t.clearWindowsTun()
	}
	if err != nil {
		logger.Error("Failed to disable TUN settings", zap.String("os", runtime.GOOS), zap.Error(err))
	} else {
		logger.Info("TUN settings disabled successfully", zap.String("os", runtime.GOOS))
	}

	return nil
}

func (t *Tun) RunTun2socks() error {
	logger.Info("Starting tun2socks")

	err := t.rp.Terminate("tun2socks")
	if err != nil {
		return err
	}

	var device string
	switch runtime.GOOS {
	case "darwin":
		device = "utun100"
	case "linux":
		device = "tun0"
	case "windows":
		device = "wintun"
	}

	if runtime.GOOS != "windows" {
		cmd = exec.Command("sudo", embed.GetTempFileName("tun2socks"), "-device", device, "-proxy", "socks5://127.0.0.1:2080")
	} else {
		cmd = exec.Command(embed.GetTempFileName("tun2socks"), "-device", device, "-proxy", "socks5://127.0.0.1:2080")
	}
	err = t.rp.Run("tun2socks", cmd, t.handlerStdout, t.waitForExit)
	if err != nil {
		return err
	}

	logger.Debug("tun2socks started successfully")
	return nil
}

func (t *Tun) handlerStdout(line string) {
	if strings.Contains(line, "[STACK] tun://") {
		var err error
		switch runtime.GOOS {
		case "darwin":
			err = t.setMacOSTun()
		case "linux":
			err = t.setLinuxTun()
		case "windows":
			err = t.setWindowsTun()
		}
		if err != nil {
			logger.Error("Failed to enable TUN settings", zap.String("os", runtime.GOOS), zap.Error(err))
			return
		}
		logger.Info("TUN settings enabled successfully", zap.String("os", runtime.GOOS))
	}
}

func (t *Tun) waitForExit() {
	if err := cmd.Wait(); err != nil {
		logger.Error("tun2socks exited with an error", zap.Error(err))
	}
}

func (t *Tun) KillTun2socks() error {
	logger.Info("Stopping tun2socks")

	err := t.rp.Kill("tun2socks", cmd)
	if err != nil {
		return err
	}

	logger.Info("tun2socks stopped successfully")
	return nil
}
