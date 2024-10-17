package proxy

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func Enable() {
	switch runtime.GOOS {
	case "darwin":
		setMacOSProxy("127.0.0.1", "2080", "2081", "2082")
	case "linux":
		setLinuxProxy("127.0.0.1", "2080", "2081", "2082")
	case "windows":
		setWindowsProxy("127.0.0.1", "2080", "2081")
	}
}

func Disable() {
	switch runtime.GOOS {
	case "darwin":
		clearMacOSProxy()
	case "linux":
		clearLinuxProxy()
	case "windows":
		clearWindowsProxy()
	}
}

func setMacOSProxy(host, httpPort, httpsPort, socksPort string) {
	proxyCommands := [][]string{
		//{"networksetup", "-setwebproxy", "Wi-Fi", host, httpPort},
		//{"networksetup", "-setsecurewebproxy", "Wi-Fi", host, httpsPort},
		{"networksetup", "-setsocksfirewallproxy", "Wi-Fi", host, socksPort},
	}

	runCommands(proxyCommands)
}

func clearMacOSProxy() {
	proxyCommands := [][]string{
		//{"networksetup", "-setwebproxystate", "Wi-Fi", "off"},
		//{"networksetup", "-setsecurewebproxystate", "Wi-Fi", "off"},
		{"networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "off"},
	}

	runCommands(proxyCommands)
}

func setLinuxProxy(host, httpPort, httpsPort, socksPort string) {
	proxyCommands := [][]string{
		//{"sh", "-c", fmt.Sprintf("export http_proxy='http://%s:%s'", host, httpPort)},
		//{"sh", "-c", fmt.Sprintf("export https_proxy='https://%s:%s'", host, httpsPort)},
		{"sh", "-c", fmt.Sprintf("export all_proxy='socks5://%s:%s'", host, socksPort)},
	}

	runCommands(proxyCommands)
}

func clearLinuxProxy() {
	proxyCommands := [][]string{
		//{"sh", "-c", "unset http_proxy"},
		//{"sh", "-c", "unset https_proxy"},
		{"sh", "-c", "unset all_proxy"},
	}

	runCommands(proxyCommands)
}

func setWindowsProxy(host, httpPort, httpsPort string) {
	cmd := exec.Command("netsh", "winhttp", "set", "proxy", fmt.Sprintf("http=%s:%s;https=%s:%s", host, httpPort, host, httpsPort))
	runCommand(cmd)
}

func clearWindowsProxy() {
	cmd := exec.Command("netsh", "winhttp", "reset", "proxy")
	runCommand(cmd)
}

func runCommands(commands [][]string) {
	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		runCommand(cmd)
	}
}

func runCommand(cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Error executing command: %v\n", err)
	}
}
