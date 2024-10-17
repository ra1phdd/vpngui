package xray_api

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"vpngui/config"
	"vpngui/internal/app/proxy"
	"vpngui/pkg/embed"
)

var cmd *exec.Cmd

type XrayAPI struct{}

func New() *XrayAPI {
	return &XrayAPI{}
}

func (x *XrayAPI) Run() {
	cmd = exec.Command(embed.GetTempFileName(), "run", "-c", "config/xray.json")

	cmd.Stderr = os.Stderr

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("ошибка получения stdout", err)
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := cmd.Start(); err != nil {
		fmt.Println("ошибка запуска xray-api", err)
	}

	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Println(line)

			if strings.Contains(line, "started") {
				proxy.Enable()

				config.JSON.ActiveVPN = true
				err = config.SaveConfig()
				if err != nil {
					return
				}
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Println("ошибка чтения stdout", err)
		}
	}()

	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("xray-api завершился с ошибкой", err)
		}
	}()

	<-sigChan

	proxy.Disable()

	if err := cmd.Process.Kill(); err != nil {
		fmt.Println("ошибка завершения xray-api", err)
	}
}

func (x *XrayAPI) Kill() {
	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			fmt.Println("ошибка при отправке сигнала SIGTERM", err)
		}
	}
	proxy.Disable()

	config.JSON.ActiveVPN = false
	err := config.SaveConfig()
	if err != nil {
		return
	}
}
