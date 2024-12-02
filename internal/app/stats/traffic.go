package stats

import (
	"bufio"
	"encoding/json"
	"go.uber.org/zap"
	"io"
	"os/exec"
	"strings"
	"vpngui/internal/app/command"
	"vpngui/internal/app/models"
	"vpngui/internal/app/proxy"
	"vpngui/internal/app/repository"
	"vpngui/pkg/embed"
	"vpngui/pkg/logger"
)

type Traffic struct {
	cr *repository.ConfigRepository
}

func NewTraffic(cr *repository.ConfigRepository) *Traffic {
	return &Traffic{
		cr: cr,
	}
}

var (
	OldTraffic     models.StatsTraffic
	CurrentTraffic models.StatsTraffic
)

func (t *Traffic) CaptureTraffic() {
	cmd := exec.Command(embed.GetTempFileName("xray-core"), "api", "statsquery")
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
		if strings.Contains(line, "resource temporarily unavailable") {
			err := embed.Init()
			if err != nil {
				logger.Error("Failed to create xray file", zap.Error(err))

				if err := proxy.Disable(); err != nil {
					logger.Error("Failed to update VPN state", zap.Error(err))
				}

				if err := t.cr.UpdateActiveVPN(false); err != nil {
					logger.Error("Failed to update VPN state", zap.Error(err))
				}
			}
		}
		jsonOutput += line
	}

	return jsonOutput
}
