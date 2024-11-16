package stats

import (
	"bufio"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os/exec"
	"runtime"
	"vpngui/internal/app/command"
	"vpngui/internal/app/models"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

type Traffic struct{}

func NewTraffic() *Traffic {
	return &Traffic{}
}

var CurrentTraffic, OldTraffic models.StatsTraffic

func (t *Traffic) CaptureTraffic() {
	var cmd *exec.Cmd
	cmdArgs := []string{"api", "statsquery"}
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", embed.GetTempFileName(), cmdArgs[0], cmdArgs[1])
	} else {
		cmd = exec.Command(embed.GetTempFileName(), cmdArgs...)
	}
	cmd.SysProcAttr = command.GetSysProcAttr()

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		logger.Error("Failed to get stdout pipe", zap.Error(err))
		return
	}

	if err := cmd.Start(); err != nil {
		logger.Error("Failed to start xray api to get statistics", zap.Error(err))
		return
	}

	jsonOutput := t.ScannerStdout(stdoutPipe)
	if jsonOutput != "" {
		OldTraffic = CurrentTraffic
		CurrentTraffic = models.StatsTraffic{}

		var statsResponse models.StatsTrafficJSON
		if err := json.Unmarshal([]byte(jsonOutput), &statsResponse); err != nil {
			logger.Error("Failed to deserialize JSON", zap.Error(err))
			return
		}

		for i := range statsResponse.Stat {
			switch statsResponse.Stat[i].Name {
			case "outbound>>>proxy>>>traffic>>>uplink":
				CurrentTraffic.ProxyUplink += statsResponse.Stat[i].Value
			case "outbound>>>proxy>>>traffic>>>downlink":
				CurrentTraffic.ProxyDownlink += statsResponse.Stat[i].Value
			case "outbound>>>direct>>>traffic>>>uplink":
				CurrentTraffic.DirectUplink += statsResponse.Stat[i].Value
			case "outbound>>>direct>>>traffic>>>downlink":
				CurrentTraffic.DirectDownlink += statsResponse.Stat[i].Value
			}
		}
		logger.Debug("Traffic capture completed", zap.Int64("ProxyUplink", CurrentTraffic.ProxyUplink),
			zap.Int64("ProxyDownlink", CurrentTraffic.ProxyDownlink),
			zap.Int64("DirectUplink", CurrentTraffic.DirectUplink),
			zap.Int64("DirectDownlink", CurrentTraffic.DirectDownlink))
	}
}

func (t *Traffic) GetTraffic(typeTraffic, typeChannel string) int64 {
	switch typeTraffic {
	case "proxy":
		switch typeChannel {
		case "uplink":
			return CurrentTraffic.ProxyUplink - OldTraffic.ProxyUplink
		case "downlink":
			return CurrentTraffic.ProxyDownlink - OldTraffic.ProxyDownlink
		}
	case "direct":
		switch typeChannel {
		case "uplink":
			return CurrentTraffic.DirectUplink - OldTraffic.DirectUplink
		case "downlink":
			return CurrentTraffic.DirectDownlink - OldTraffic.DirectDownlink
		}
	}
	return 0
}

func (t *Traffic) ScannerStdout(stdoutPipe io.ReadCloser) string {
	var jsonOutput string
	scanner := bufio.NewScanner(stdoutPipe)
	defer stdoutPipe.Close()

	for scanner.Scan() {
		line := scanner.Text()
		jsonOutput += line
	}

	return jsonOutput
}
