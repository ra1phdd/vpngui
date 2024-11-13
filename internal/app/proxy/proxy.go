package proxy

import (
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"runtime"
	"vpngui/pkg/logger"
)

func Enable() error {
	logger.Info("Enabling proxy settings based on OS")

	var err error
	switch runtime.GOOS {
	case "darwin":
		err = setMacOSProxy("127.0.0.1", "2080")
	case "linux":
		err = setLinuxProxy("127.0.0.1", "2080")
	case "windows":
		err = setWindowsProxy("127.0.0.1", "2080")
	}
	if err != nil {
		logger.Error("Failed to enable proxy settings", zap.String("os", runtime.GOOS), zap.Error(err))
	} else {
		logger.Info("Proxy settings enabled successfully", zap.String("os", runtime.GOOS))
	}

	return nil
}

func Disable() error {
	logger.Info("Disabling proxy settings based on OS")

	var err error
	switch runtime.GOOS {
	case "darwin":
		err = clearMacOSProxy()
	case "linux":
		err = clearLinuxProxy()
	case "windows":
		err = clearWindowsProxy()
	}
	if err != nil {
		logger.Error("Failed to disable proxy settings", zap.String("os", runtime.GOOS), zap.Error(err))
	} else {
		logger.Info("Proxy settings disabled successfully", zap.String("os", runtime.GOOS))
	}

	return nil
}

func setMacOSProxy(host, port string) error {
	proxyCommands := [][]string{
		{"networksetup", "-setwebproxy", "Wi-Fi", host, port},
		{"networksetup", "-setsocksfirewallproxy", "Wi-Fi", host, port},
	}

	err := runCommands(proxyCommands)
	if err != nil {
		return err
	}

	return nil
}

func clearMacOSProxy() error {
	proxyCommands := [][]string{
		{"networksetup", "-setwebproxystate", "Wi-Fi", "off"},
		{"networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "off"},
	}

	err := runCommands(proxyCommands)
	if err != nil {
		return err
	}

	return nil
}

func setLinuxProxy(host, port string) error {
	proxyCommands := [][]string{
		{"sh", "-c", fmt.Sprintf("export http_proxy='https://%s:%s'", host, port)},
		{"sh", "-c", fmt.Sprintf("export all_proxy='socks5://%s:%s'", host, port)},
	}

	err := runCommands(proxyCommands)
	if err != nil {
		return err
	}

	return nil
}

func clearLinuxProxy() error {
	proxyCommands := [][]string{
		{"sh", "-c", "unset http_proxy"},
		{"sh", "-c", "unset all_proxy"},
	}

	err := runCommands(proxyCommands)
	if err != nil {
		return err
	}

	return nil
}

func setWindowsProxy(host, port string) error {
	cmd := exec.Command("netsh", "winhttp", "set", "proxy", fmt.Sprintf("http=%s:%s;socks=%s:%s", host, port, host, port))
	err := runCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func clearWindowsProxy() error {
	cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	err := runCommand(cmd)
	if err != nil {
		return err
	}

	return nil
}

func runCommands(commands [][]string) error {
	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		err := runCommand(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logger.Error("Command execution failed", zap.String("cmd", cmd.String()), zap.Error(err))
		return err
	}

	return nil
}
