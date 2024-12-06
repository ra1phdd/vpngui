package tun

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"syscall"
	"vpngui/internal/app/command"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

var DefaultInterface, DefaultIP, DefaultGW string
var cmd *exec.Cmd

func Enable() error {
	logger.Info("Enabling TUN settings based on OS")

	if runtime.GOOS == "linux" {
		commands := [][]string{
			{"ip", "tuntap", "add", "mode", "tun", "dev", "tun0"},
			{"ip", "addr", "add", "198.18.0.1/15", "dev", "tun0"},
			{"ip", "link", "set", "dev", "tun0", "up"},
		}

		err := runCommands(commands, false)
		if err != nil {
			return err
		}
	} else if runtime.GOOS == "windows" {
		commands := [][]string{
			{"netsh", "interface", "set", "interface", "name=\"wintun\"", "admin=disabled"},
		}
		err := runCommands(commands, true)
		if err != nil {
			return err
		}
	}

	var err error
	if runtime.GOOS != "windows" {
		DefaultInterface, DefaultIP, err = GetDefaultInterface()
		if err != nil {
			return err
		}
		DefaultGW, err = GetDefaultGateway()
		if err != nil {
			return err
		}
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

func setMacOSTun() error {
	err := clearMacOSTun()
	if err != nil {
		return err
	}

	commands := [][]string{
		{"sudo", "ifconfig", "utun100", "198.18.0.1", "198.18.0.1", "up"},
		{"sudo", "route", "add", "-net", "1.0.0.0/8", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "2.0.0.0/7", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "4.0.0.0/6", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "8.0.0.0/5", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "16.0.0.0/4", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "32.0.0.0/3", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "64.0.0.0/2", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "128.0.0.0/1", "198.18.0.1"},
		{"sudo", "route", "add", "-net", "198.18.0.0/15", "198.18.0.1"},
	}

	return runCommands(commands, false)
}

func clearMacOSTun() error {
	commands := [][]string{
		{"sudo", "ifconfig", "utun100", "198.18.0.1", "198.18.0.1", "down"},
		{"sudo", "route", "delete", "default"},
		{"sudo", "route", "delete", "1.0.0.0/8"},
		{"sudo", "route", "delete", "2.0.0.0/7"},
		{"sudo", "route", "delete", "4.0.0.0/6"},
		{"sudo", "route", "delete", "8.0.0.0/5"},
		{"sudo", "route", "delete", "16.0.0.0/4"},
		{"sudo", "route", "delete", "32.0.0.0/3"},
		{"sudo", "route", "delete", "64.0.0.0/2"},
		{"sudo", "route", "delete", "128.0.0.0/1"},
		{"sudo", "route", "delete", "198.18.0.0/15"},
		{"sudo", "route", "add", "default", DefaultGW},
		{"sudo", "route", "add", "-net", fmt.Sprintf("%s/32", DefaultIP), "-interface", DefaultInterface},
	}

	return runCommands(commands, true)
}

func setLinuxTun() error {
	commands := [][]string{
		{"ip", "route", "del", "default"},
		{"ip", "route", "add", "default", "via", "198.18.0.1", "dev", "tun0", "metric", "1"},
		{"ip", "route", "add", "default", "via", "172.17.0.1", "dev", DefaultInterface, "metric", "10"},
	}

	return runCommands(commands, false)
}

func clearLinuxTun() error {
	commands := [][]string{
		{"ip", "route", "del", "default"},
	}

	return runCommands(commands, false)
}

func setWindowsTun() error {
	err := clearWindowsTun()
	if err != nil {
		return err
	}

	commands := [][]string{
		{"netsh", "interface", "ipv4", "set", "address", "name=\"wintun\"", "source=static", "addr=192.168.123.1", "mask=255.255.255.0"},
		{"netsh", "interface", "ipv4", "set", "dnsservers", "name=\"wintun\"", "static", "address=8.8.8.8", "register=none", "validate=no"},
		{"netsh", "interface", "ipv4", "add", "route", "0.0.0.0/0", "\"wintun\"", "192.168.123.1", "metric=1"},
	}

	return runCommands(commands, false)
}

func clearWindowsTun() error {
	commands := [][]string{
		{"netsh", "interface", "ipv4", "delete", "route", "0.0.0.0/0", "\"wintun\""},
	}

	return runCommands(commands, true)
}

func runCommands(commands [][]string, ignoreErr bool) error {
	for _, args := range commands {
		logger.Debug("Executing command", zap.String("cmd", strings.Join(args, " ")))
		cmd := exec.Command(args[0], args[1:]...)
		err := runCommand(cmd, ignoreErr)
		if err != nil {
			return err
		}
	}

	return nil
}

func runCommand(cmd *exec.Cmd, ignoreErr bool) error {
	cmd.SysProcAttr = command.GetSysProcAttr()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil && !ignoreErr {
		logger.Error("Command execution failed", zap.String("cmd", cmd.String()), zap.Error(err))
		return err
	}

	return nil
}

func RunTun2socks() error {
	logger.Info("Starting tun2socks")

	err := TerminateProcesses()
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

	if runtime.GOOS == "windows" {
		cmd = exec.Command(embed.GetTempFileName("tun2socks"), "-device", device, "-proxy", "socks5://127.0.0.1:2080")
	} else {
		cmd = exec.Command("sudo", embed.GetTempFileName("tun2socks"), "-device", device, "-proxy", "socks5://127.0.0.1:2080")
	}
	cmd.SysProcAttr = command.GetSysProcAttr()
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

	go handleStdout(stdoutPipe)
	go waitForExit()

	logger.Debug("tun2socks started successfully")

	return nil
}

func KillTun2socks() error {
	logger.Info("Stopping tun2socks")
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

	logger.Debug("tun2socks stopped successfully")

	return nil
}

func TerminateProcesses() error {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("taskkill", "/FI", "IMAGENAME eq tun2socks*", "/F")
	} else {
		cmd = exec.Command("pkill", "-f", "tun2socks")
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("Failed to terminate tun2socks processes", zap.Error(err), zap.String("output", string(output)))
		return err
	}

	logger.Info("Successfully terminated tun2socks processes")
	return nil
}

func handleStdout(stdoutPipe io.ReadCloser) {
	logger.Info("Handling stdout for tun2socks")
	scanner := bufio.NewScanner(stdoutPipe)
	defer stdoutPipe.Close()

	for scanner.Scan() {
		line := scanner.Text()
		logger.Info("Received line from tun2socks", zap.String("line", line))
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
			} else {
				logger.Info("TUN settings enabled successfully", zap.String("os", runtime.GOOS))
			}
		}
	}
}

func waitForExit() {
	if err := cmd.Wait(); err != nil {
		logger.Error("tun2socks exited with an error", zap.Error(err))
	}
}

func GetDefaultInterface() (string, string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return "", "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp != 0 && iface.Flags&net.FlagLoopback == 0 {
			addrs, _ := iface.Addrs()
			for _, addr := range addrs {
				if ip, ok := addr.(*net.IPNet); ok && ip.IP.To4() != nil {
					return iface.Name, ip.IP.String(), nil
				}
			}
		}
	}

	return "", "", errors.New("no default interface found")
}

func GetDefaultGateway() (string, error) {
	cmd := exec.Command("sh", "-c", "route -n get default | grep 'gateway' | awk '{print $2}'")

	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		logger.Error("Failed to execute command", zap.Error(err))
		return "", err
	}

	gateway := strings.TrimSpace(out.String())
	if gateway == "" {
		return "", errors.New("no default gateway found")
	}

	return gateway, nil
}
