package network

import (
	"bytes"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net"
	"os/exec"
	"runtime"
	"strings"
	"vpngui/pkg/logger"
)

var (
	DefaultIP        string
	DefaultGateway   string
	DefaultInterface string
)

type Network struct{}

func New() *Network {
	return &Network{}
}

func (n *Network) Init() error {
	var err error
	DefaultInterface, err = n.GetDefaultInterface()
	if err != nil {
		logger.Error("Failed get default interface", zap.Error(err))
		return err
	}
	DefaultGateway, err = n.GetDefaultGateway()
	if err != nil {
		logger.Error("Failed get default gateway", zap.Error(err))
		return err
	}
	DefaultIP, err = n.GetDefaultIP(DefaultInterface)
	if err != nil {
		logger.Error("Failed get default ip", zap.Error(err))
		return err
	}

	return nil
}

func (n *Network) GetDefaultInterface() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("sh", "-c", "route -n get default | grep 'interface' | awk '{print $2}'")
	case "linux":
		cmd = exec.Command("sh", "-c", "ip route show default | awk '/default/ {print $5}'")
	case "windows":
		cmd = exec.Command("cmd", "/C", `powershell -Command "Get-NetRoute -DestinationPrefix '0.0.0.0/0' | Select-Object -ExpandProperty InterfaceAlias"`)
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Error("Failed to execute command", zap.Error(err), zap.String("stdout", stdout.String()), zap.String("stderr", stderr.String()))
		return "", err
	}

	for _, line := range strings.Split(stdout.String(), "\n") {
		if iface := strings.TrimSpace(line); iface != "" {
			return iface, nil
		}
	}

	return "", errors.New("no default interface found")
}

func (n *Network) GetDefaultIP(interfaceName string) (string, error) {
	iface, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", fmt.Errorf("failed to get interface %s: %w", interfaceName, err)
	}

	addrs, err := iface.Addrs()
	if err != nil {
		return "", fmt.Errorf("failed to get addresses for interface %s: %w", interfaceName, err)
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP

		if ip.To4() == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			continue
		}

		return ip.String(), nil
	}

	return "", fmt.Errorf("no suitable IP address found for interface %s", interfaceName)
}

func (n *Network) GetDefaultGateway() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("sh", "-c", "route -n get default | grep 'gateway' | awk '{print $2}'")
	case "linux":
		cmd = exec.Command("sh", "-c", "ip route show default | awk '/default/ {print $3}'")
	case "windows":
		cmd = exec.Command("cmd", "/C", "for /f \"tokens=3\" %a in ('route print ^| findstr \"\\<0.0.0.0\\>\"') do @echo %a")
	default:
		return "", fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		logger.Error("Failed to execute command", zap.Error(err), zap.String("stdout", stdout.String()), zap.String("stderr", stderr.String()))
		return "", err
	}

	for _, line := range strings.Split(stdout.String(), "\n") {
		if gateway := strings.TrimSpace(line); gateway != "" {
			return gateway, nil
		}
	}

	return "", errors.New("no default gateway found")
}
