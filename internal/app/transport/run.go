package transport

import (
	"go.uber.org/zap"
	"vpngui/internal/app/repository"
	"vpngui/internal/app/transport/proxy"
	"vpngui/internal/app/transport/tun"
	"vpngui/pkg/logger"
)

type Transport struct {
	cr  *repository.ConfigRepository
	p   *proxy.Proxy
	tun *tun.Tun
}

func New(cr *repository.ConfigRepository, p *proxy.Proxy, tun *tun.Tun) *Transport {
	return &Transport{
		cr:  cr,
		p:   p,
		tun: tun,
	}
}

func (t *Transport) Enable() error {
	getConfig, err := t.cr.GetConfig()
	if err != nil {
		logger.Error("Failed to get config", zap.Error(err))
		return err
	}

	if getConfig.VpnMode == "proxy" {
		return t.p.Enable()
	}
	return t.tun.Enable()
}

func (t *Transport) Disable() error {
	getConfig, err := t.cr.GetConfig()
	if err != nil {
		logger.Error("Failed to get config", zap.Error(err))
		return err
	}

	if getConfig.VpnMode == "proxy" {
		return t.p.Disable()
	}
	return t.tun.Disable()
}
