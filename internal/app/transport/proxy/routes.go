package proxy

import (
	"fmt"
	"go.uber.org/zap"
	"vpngui/pkg/logger"
)

func (p *Proxy) setMacOSProxy(host, port string) error {
	commands := [][]string{
		{"networksetup", "-setwebproxy", "Wi-Fi", host, port},
		{"networksetup", "-setsocksfirewallproxy", "Wi-Fi", host, port},
	}

	err := p.rc.RunCommands(commands, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) clearMacOSProxy() error {
	commands := [][]string{
		{"networksetup", "-setwebproxystate", "Wi-Fi", "off"},
		{"networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "off"},
	}

	err := p.rc.RunCommands(commands, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) setLinuxProxy(host, port string) error {
	commands := [][]string{
		{"sh", "-c", fmt.Sprintf("export http_proxy='https://%s:%s'", host, port)},
		{"sh", "-c", fmt.Sprintf("export all_proxy='socks5://%s:%s'", host, port)},
	}

	err := p.rc.RunCommands(commands, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) clearLinuxProxy() error {
	commands := [][]string{
		{"sh", "-c", "unset http_proxy"},
		{"sh", "-c", "unset all_proxy"},
	}

	err := p.rc.RunCommands(commands, false)
	if err != nil {
		return err
	}

	return nil
}

func (p *Proxy) setWindowsProxy(host, port string) error {
	proxy := fmt.Sprintf("%s:%s", host, port)

	commands := [][]string{
		{"reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f"},
		{"reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyServer", "/t", "REG_SZ", "/d", proxy, "/f"},
		{"netsh", "winhttp", "set", "proxy", fmt.Sprintf("http=%s:%s", host, port)},
	}

	err := p.rc.RunCommands(commands, false)
	if err != nil {
		return err
	}

	if err := notifySettingsChange(); err != nil {
		logger.Error("Failed to notify settings change", zap.Error(err))
		return err
	}

	return nil
}

func (p *Proxy) clearWindowsProxy() error {
	commands := [][]string{
		{"reg", "add", `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f"},
		{"netsh", "winhttp", "reset", "proxy"},
	}

	err := p.rc.RunCommands(commands, false)
	if err != nil {
		return err
	}

	if err := notifySettingsChange(); err != nil {
		logger.Error("Failed to notify settings change", zap.Error(err))
		return err
	}

	return nil
}
