package xray_api

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"vpngui/internal/app/proxy"
	"vpngui/internal/app/repository"
	"vpngui/pkg/embed"
)

var cmd *exec.Cmd

type RunXrayAPI struct {
	cr       *repository.ConfigRepository
	sigChan  chan os.Signal
	doneChan chan bool
}

func NewRun(cr *repository.ConfigRepository) *RunXrayAPI {
	return &RunXrayAPI{
		cr:       cr,
		sigChan:  make(chan os.Signal, 1),
		doneChan: make(chan bool, 1),
	}
}

func (x *RunXrayAPI) Run() bool {
	cmd = exec.Command(embed.GetTempFileName(), "run", "-c", "config/xray.json")
	cmd.Stderr = os.Stderr

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("ошибка получения stdout", err)
		return false
	}

	signal.Notify(x.sigChan, syscall.SIGINT, syscall.SIGTERM)

	if err := cmd.Start(); err != nil {
		fmt.Println("ошибка запуска xray-api", err)
		return false
	}

	go x.handleStdout(stdoutPipe)
	go x.handleSignals()
	go x.waitForExit()

	return true
}

func (x *RunXrayAPI) Kill() bool {
	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
			fmt.Println("ошибка при отправке сигнала SIGTERM:", err)
			return false
		}
	}

	proxy.Disable()
	if err := x.cr.UpdateActiveVPN(false); err != nil {
		fmt.Println("ошибка обновления VPN состояния:", err)
		return false
	}

	return true
}

func (x *RunXrayAPI) handleStdout(stdoutPipe io.ReadCloser) {
	scanner := bufio.NewScanner(stdoutPipe)
	defer stdoutPipe.Close()

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		if strings.Contains(line, "started") {
			proxy.Enable()
			if err := x.cr.UpdateActiveVPN(true); err != nil {
				fmt.Println("ошибка обновления VPN состояния:", err)
				return
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("ошибка чтения stdout:", err)
	}
}

func (x *RunXrayAPI) waitForExit() {
	if err := cmd.Wait(); err != nil {
		fmt.Println("xray-api завершился с ошибкой:", err)
	}
	x.doneChan <- true
}

func (x *RunXrayAPI) handleSignals() {
	select {
	case <-x.sigChan:
		x.Kill()
	case <-x.doneChan:
		signal.Stop(x.sigChan)
	}
}
