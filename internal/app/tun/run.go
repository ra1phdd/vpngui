package tun

import (
	"go.uber.org/zap"
	"os/exec"
	"runtime"
	"strings"
	"vpngui/internal/app/command"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

var cmd *exec.Cmd

func Enable() error {
	logger.Info("Enabling TUN settings based on OS")

	if runtime.GOOS == "linux" {
		commands := [][]string{
			{"ip", "tuntap", "add", "mode", "tun", "dev", "tun0"},
			{"ip", "addr", "add", "198.18.0.1/15", "dev", "tun0"},
			{"ip", "link", "set", "dev", "tun0", "up"},
		}

		err := command.RunCommands(commands, false)
		if err != nil {
			return err
		}
	} else if runtime.GOOS == "windows" {
		commands := [][]string{
			{"netsh", "interface", "set", "interface", "name=\"wintun\"", "admin=disabled"},
		}
		err := command.RunCommands(commands, true)
		if err != nil {
			return err
		}
	}

	var err error
	DefaultInterface, err = GetDefaultInterface()
	if err != nil {
		return err
	}
	DefaultIP, err = GetDefaultIP(DefaultInterface)
	if err != nil {
		return err
	}
	DefaultGateway, err = GetDefaultGateway()
	if err != nil {
		return err
	}

	err = RunTun2socks()
	if err != nil {
		return err
	}

	return nil
}

func Disable() error {
	logger.Info("Disabling TUN settings based on OS")

	err := KillTun2socks()
	if err != nil {
		return err
	}

	switch runtime.GOOS {
	case "darwin":
		err = clearMacOSTun()
	case "linux":
		err = clearLinuxTun()
	case "windows":
		err = clearWindowsTun()
	}
	if err != nil {
		logger.Error("Failed to disable TUN settings", zap.String("os", runtime.GOOS), zap.Error(err))
	} else {
		logger.Info("TUN settings disabled successfully", zap.String("os", runtime.GOOS))
	}

	return nil
}

func RunTun2socks() error {
	logger.Info("Starting tun2socks")

	err := command.TerminateProcesses("tun2socks")
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
	err = command.RunProcess("tun2socks", cmd, handlerStdout, waitForExit)
	if err != nil {
		return err
	}

	logger.Debug("tun2socks started successfully")
	return nil
}

func handlerStdout(line string) {
	if strings.Contains(line, "[STACK] tun://") {
		var err error
		switch runtime.GOOS {
		case "darwin":
			err = setMacOSTun()
		case "linux":
			err = setLinuxTun()
		case "windows":
			err = setWindowsTun()
		}
		if err != nil {
			logger.Error("Failed to enable TUN settings", zap.String("os", runtime.GOOS), zap.Error(err))
			return
		}
		logger.Info("TUN settings enabled successfully", zap.String("os", runtime.GOOS))
	}
}

func waitForExit() {
	if err := cmd.Wait(); err != nil {
		logger.Error("tun2socks exited with an error", zap.Error(err))
	}
}

func KillTun2socks() error {
	logger.Info("Stopping tun2socks")

	err := command.KillProcess("tun2socks", cmd)
	if err != nil {
		return err
	}

	logger.Info("tun2socks stopped successfully")
	return nil
}
