package xray_api

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"regexp"
	"slices"
	"strings"
	"time"
	"vpngui/internal/app/config"
	"vpngui/internal/app/models"
	"vpngui/internal/app/repository"
	"vpngui/pkg/logger"
)

type RoutesXrayAPI struct {
	run *RunXrayAPI
	rr  *repository.RoutesRepository
}

func NewRoutes(run *RunXrayAPI, rr *repository.RoutesRepository) *RoutesXrayAPI {
	return &RoutesXrayAPI{
		run: run,
		rr:  rr,
	}
}

func (x *RoutesXrayAPI) DisableRoutes() error {
	logger.Info("Disabling routes...")
	return x.toggleRoutes(false)
}

func (x *RoutesXrayAPI) EnableRoutes() error {
	logger.Info("Enabling routes...")
	return x.toggleRoutes(true)
}

func (x *RoutesXrayAPI) toggleRoutes(enable bool) error {
	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Failed to get config", zap.Error(err))
		return err
	}

	if err := x.updateRoutingConfig(getConfig.ListMode); err != nil {
		logger.Error("Failed to update routing config", zap.Error(err))
		return err
	}

	if err := x.run.cr.UpdateDisableRoutes(enable); err != nil {
		logger.Error("Failed to update disable routes", zap.Error(err))
		return err
	}

	outbound1, outbound2 := "proxy", "direct"
	if enable {
		outbound1 = x.switchModeToTag(getConfig.ListMode)
		if outbound1 == "proxy" {
			outbound1, outbound2 = outbound2, outbound1
		}
	}
	if err := x.SwapOutbounds(&config.Xray.Outbounds, outbound1, outbound2); err != nil {
		logger.Error("Failed to swap outbounds", zap.Error(err))
		return err
	}

	if err := config.Save(); err != nil {
		logger.Error("Failed to save config", zap.Error(err))
		return err
	}

	return x.restartVPNIfActive(getConfig.ActiveVPN)
}

func (x *RoutesXrayAPI) EnableBlackList() error {
	logger.Info("Enabling blacklist...")
	return x.toggleListMode("blacklist", "direct", "proxy")
}

func (x *RoutesXrayAPI) DisableBlackList() error {
	logger.Info("Enabling whitelist...")
	return x.toggleListMode("whitelist", "proxy", "direct")
}

func (x *RoutesXrayAPI) toggleListMode(listMode, outbound1, outbound2 string) error {
	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Failed to get config", zap.Error(err))
		return err
	}

	if err := x.updateRoutingConfig(listMode); err != nil {
		logger.Error("Failed to update routing config", zap.Error(err))
		return err
	}

	if err := x.run.cr.UpdateListMode(listMode); err != nil {
		logger.Error("Failed to update list mode", zap.Error(err))
		return err
	}

	if err := x.SwapOutbounds(&config.Xray.Outbounds, outbound1, outbound2); err != nil {
		logger.Error("Failed to swap outbounds", zap.Error(err))
		return err
	}

	if err := config.Save(); err != nil {
		logger.Error("Failed to save config", zap.Error(err))
		return err
	}

	return x.restartVPNIfActive(getConfig.ActiveVPN)
}

func (x *RoutesXrayAPI) updateRoutingConfig(listMode string) error {
	logger.Debug("Updating routing config with list mode", zap.String("listMode", listMode))
	if config.Xray.Routing == nil {
		logger.Debug("Creating new routing config...")
		config.Xray.Routing = new(models.RoutingConfig)
	}

	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		logger.Error("Failed to get routes", zap.Error(err))
		return err
	}
	*config.Xray.Routing = x.convertToRoutingConfig(getRoutes)

	logger.Debug("Routing config updated successfully.")

	return nil
}

func (x *RoutesXrayAPI) restartVPNIfActive(active bool) error {
	if !active {
		logger.Debug("VPN is not active. Skipping restart.")
		return nil
	}

	err := x.run.Kill()
	if err != nil {
		logger.Error("Failed to killed xray-api", zap.Error(err))
		return err
	}
	if err := x.waitForVPNState(false); err != nil {
		logger.Error("Failed to wait for VPN state", zap.Error(err))
		return err
	}

	err = x.run.Run()
	if err != nil {
		logger.Error("Failed to running xray-api", zap.Error(err))
		return err
	}
	if err := x.waitForVPNState(false); err != nil {
		logger.Error("Failed to wait for VPN state", zap.Error(err))
		return err
	}

	return nil
}

func (x *RoutesXrayAPI) SwapOutbounds(outbounds *[]models.OutboundConfig, tag1, tag2 string) error {
	index1, index2 := -1, -1

	for i, outbound := range *outbounds {
		if outbound.Tag == tag1 {
			index1 = i
		}
		if outbound.Tag == tag2 {
			index2 = i
		}
	}

	if index1 == -1 || index2 == -1 {
		logger.Warn("One or both tags not found.")
		return errors.New("one or both tags not found")
	} else if index1 == 0 && index2 == 1 {
		return nil
	}

	(*outbounds)[index1], (*outbounds)[index2] = (*outbounds)[index2], (*outbounds)[index1]

	return nil
}

func (x *RoutesXrayAPI) GetDomain(listMode string) string {
	logger.Debug("Fetching domains...", zap.String("listMode", listMode))
	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		logger.Error("Error fetching routes", zap.String("listMode", listMode), zap.Error(err))
		return ""
	}

	var domains []string
	for _, rule := range getRoutes.Rules {
		if rule.RuleType == "domain" {
			domains = append(domains, rule.RuleValue)
		}
	}
	logger.Debug("Fetched domains", zap.String("listMode", listMode), zap.Int("domainCount", len(domains)))

	return strings.Join(domains, "\n")
}

func (x *RoutesXrayAPI) GetIP(listMode string) string {
	logger.Debug("Fetching IPs...", zap.String("listMode", listMode))
	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		logger.Error("Error fetching routes", zap.String("listMode", listMode), zap.Error(err))
		return ""
	}

	var ips []string
	for _, rule := range getRoutes.Rules {
		if rule.RuleType == "ip" {
			ips = append(ips, rule.RuleValue)
		}
	}
	logger.Debug("Fetched IPs", zap.String("listMode", listMode), zap.Int("ipCount", len(ips)))

	return strings.Join(ips, "\n")
}

func (x *RoutesXrayAPI) GetPort(listMode string) string {
	logger.Debug("Fetching ports...", zap.String("listMode", listMode))
	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		logger.Error("Error fetching routes", zap.String("listMode", listMode), zap.Error(err))
		return ""
	}

	var ports []string
	for _, rule := range getRoutes.Rules {
		if rule.RuleType == "port" {
			ports = append(ports, rule.RuleValue)
		}
	}
	logger.Debug("Fetched ports", zap.String("listMode", listMode), zap.Int("portCount", len(ports)))

	return strings.Join(ports, ", ")
}

func (x *RoutesXrayAPI) AddDomain(listMode string, domain string) {
	if !x.isValidDomain(domain) {
		logger.Warn("Невалидный домен", zap.String("domain", domain), zap.String("listMode", listMode))
		return
	}

	err := x.isFirstItem(listMode)
	if err != nil {
		logger.Error("Ошибка проверки на отсутствие соответствующего поля в Routing", zap.String("domain", domain), zap.String("listMode", listMode), zap.Error(err))
		return
	}

	err = x.rr.AddRule(listMode, "domain", domain)
	if err != nil {
		logger.Error("Ошибка добавления правила в БД", zap.String("domain", domain), zap.String("listMode", listMode), zap.Error(err))
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == x.switchModeToTag(listMode) {
			config.Xray.Routing.Rules[i].Domain = append(config.Xray.Routing.Rules[i].Domain, domain)
		}
	}

	if err := config.Save(); err != nil {
		logger.Error("Ошибка сохранения JSON-конфига", zap.String("domain", domain), zap.String("listMode", listMode), zap.Error(err))
		return
	}

	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Ошибка получения конфига из БД", zap.Error(err))
		return
	}

	err = x.restartVPNIfActive(getConfig.ActiveVPN)
	if err != nil {
		logger.Error("Ошибка перезагрузки VPN", zap.Error(err))
		return
	}
}

func (x *RoutesXrayAPI) AddIP(listMode string, ip string) {
	if !x.isValidIP(ip) {
		logger.Warn("Невалидный IP-адрес", zap.String("ip", ip), zap.String("listMode", listMode))
		return
	}

	err := x.isFirstItem(listMode)
	if err != nil {
		return
	}

	err = x.rr.AddRule(listMode, "ip", ip)
	if err != nil {
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == x.switchModeToTag(listMode) {
			config.Xray.Routing.Rules[i].IP = append(config.Xray.Routing.Rules[i].IP, ip)
		}
	}

	if err := config.Save(); err != nil {
		return
	}

	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Ошибка получения конфига из БД", zap.Error(err))
		return
	}

	err = x.restartVPNIfActive(getConfig.ActiveVPN)
	if err != nil {
		logger.Error("Ошибка перезагрузки VPN", zap.Error(err))
		return
	}
}

func (x *RoutesXrayAPI) AddPort(listMode string, port string) {
	if !x.isValidPort(port) {
		logger.Warn("Невалидный порт", zap.String("port", port), zap.String("listMode", listMode))
		return
	}

	err := x.isFirstItem(listMode)
	if err != nil {
		return
	}

	err = x.rr.AddRule(listMode, "port", port)
	if err != nil {
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag != x.switchModeToTag(listMode) {
			continue
		}

		if len(config.Xray.Routing.Rules[i].Port) == 0 {
			config.Xray.Routing.Rules[i].Port = port
		} else {
			config.Xray.Routing.Rules[i].Port += fmt.Sprintf(",%s", port)
		}
	}

	if err := config.Save(); err != nil {
		return
	}

	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Ошибка получения конфига из БД", zap.Error(err))
		return
	}

	err = x.restartVPNIfActive(getConfig.ActiveVPN)
	if err != nil {
		logger.Error("Ошибка перезагрузки VPN", zap.Error(err))
		return
	}
}

func (x *RoutesXrayAPI) DelDomain(listMode string, domain string) {
	if !x.isValidDomain(domain) {
		logger.Warn("Невалидный домен", zap.String("domain", domain), zap.String("listMode", listMode))
		return
	}

	err := x.rr.DeleteRule(listMode, "domain", domain)
	if err != nil {
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == x.switchModeToTag(listMode) {
			idx := slices.Index(config.Xray.Routing.Rules[i].Domain, domain)
			if idx != -1 {
				config.Xray.Routing.Rules[i].Domain = slices.Delete(config.Xray.Routing.Rules[i].Domain, idx, idx+1)
			}
		}
	}

	err = x.isLastItem(listMode)
	if err != nil {
		return
	}

	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Ошибка получения конфига из БД", zap.Error(err))
		return
	}

	err = x.restartVPNIfActive(getConfig.ActiveVPN)
	if err != nil {
		logger.Error("Ошибка перезагрузки VPN", zap.Error(err))
		return
	}
}

func (x *RoutesXrayAPI) DelIP(listMode string, ip string) {
	if !x.isValidIP(ip) {
		logger.Warn("Невалидный IP-адрес", zap.String("ip", ip), zap.String("listMode", listMode))
		return
	}

	err := x.rr.DeleteRule(listMode, "ip", ip)
	if err != nil {
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == x.switchModeToTag(listMode) {
			idx := slices.Index(config.Xray.Routing.Rules[i].IP, ip)
			if idx != -1 {
				config.Xray.Routing.Rules[i].IP = slices.Delete(config.Xray.Routing.Rules[i].IP, idx, idx+1)
			}
		}
	}

	err = x.isLastItem(listMode)
	if err != nil {
		return
	}

	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Ошибка получения конфига из БД", zap.Error(err))
		return
	}

	err = x.restartVPNIfActive(getConfig.ActiveVPN)
	if err != nil {
		logger.Error("Ошибка перезагрузки VPN", zap.Error(err))
		return
	}
}

func (x *RoutesXrayAPI) DelPort(listMode string, port string) {
	if !x.isValidPort(port) {
		logger.Warn("Невалидный порт", zap.String("port", port), zap.String("listMode", listMode))
		return
	}

	err := x.rr.DeleteRule(listMode, "port", port)
	if err != nil {
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag != x.switchModeToTag(listMode) {
			continue
		}

		if len(config.Xray.Routing.Rules[i].Port) == 0 {
			break
		}

		re := regexp.MustCompile(fmt.Sprintf(",?%s,?", port))
		config.Xray.Routing.Rules[i].Port = re.ReplaceAllString(config.Xray.Routing.Rules[i].Port, "")
	}

	err = x.isLastItem(listMode)
	if err != nil {
		return
	}

	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Ошибка получения конфига из БД", zap.Error(err))
		return
	}

	err = x.restartVPNIfActive(getConfig.ActiveVPN)
	if err != nil {
		logger.Error("Ошибка перезагрузки VPN", zap.Error(err))
		return
	}
}

func (x *RoutesXrayAPI) isValidDomain(domain string) bool {
	regex := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,63}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(domain)
}

func (x *RoutesXrayAPI) isValidIP(ip string) bool {
	regex := `^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])$`
	re := regexp.MustCompile(regex)
	return re.MatchString(ip)
}

func (x *RoutesXrayAPI) isValidPort(port string) bool {
	regex := `^(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{0,3})$`
	re := regexp.MustCompile(regex)
	return re.MatchString(port)
}

func (x *RoutesXrayAPI) isFirstItem(listMode string) error {
	if config.Xray.Routing != nil {
		config.Xray.Routing = new(models.RoutingConfig)
	}

	outboundTag := x.switchModeToTag(listMode)

	found := false
	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			found = true
		}
	}

	if !found {
		newRules := models.RoutingRule{
			Type:        "field",
			OutboundTag: outboundTag,
		}

		config.Xray.Routing.Rules = append(config.Xray.Routing.Rules, newRules)
	}

	err := config.Save()
	if err != nil {
		return err
	}

	return nil
}

func (x *RoutesXrayAPI) isLastItem(listMode string) error {
	outboundTag := x.switchModeToTag(listMode)

	found := false
	var item int
	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			if len(config.Xray.Routing.Rules[i].Domain) == 0 && len(config.Xray.Routing.Rules[i].IP) == 0 && len(config.Xray.Routing.Rules[i].Port) == 0 {
				found = true
				item = i
			}
		}
	}

	if found {
		config.Xray.Routing.Rules = append(config.Xray.Routing.Rules[:item], config.Xray.Routing.Rules[item+1:]...)
	}

	err := config.Save()
	if err != nil {
		return err
	}

	return nil
}

func (x *RoutesXrayAPI) convertToRoutingConfig(listConfig models.ListConfig) models.RoutingConfig {
	if len(listConfig.Rules) == 0 {
		return models.RoutingConfig{
			DomainStrategy: "AsIs",
			DomainMatcher:  "hybrid",
			Rules:          []models.RoutingRule{},
			Balancers:      nil,
		}
	}

	var routingConfig models.RoutingConfig

	routingConfig.DomainStrategy = listConfig.DomainStrategy
	routingConfig.DomainMatcher = listConfig.DomainMatcher

	outboundTag := x.switchModeToTag(listConfig.Type)
	var domains, ips, ports []string
	for _, rule := range listConfig.Rules {
		switch rule.RuleType {
		case "domain":
			domains = append(domains, rule.RuleValue)
		case "ip":
			ips = append(ips, rule.RuleValue)
		case "port":
			ports = append(ports, rule.RuleValue)
		}
	}

	var routingRule models.RoutingRule
	routingRule.Type = "field"
	routingRule.OutboundTag = outboundTag
	routingRule.Domain = domains
	routingRule.IP = ips
	routingRule.Port = strings.Join(ports, ",")
	routingConfig.Rules = append(routingConfig.Rules, routingRule)

	return routingConfig
}

func (x *RoutesXrayAPI) waitForVPNState(expectedState bool) error {
	for {
		getConfig, err := x.run.cr.GetConfig()
		if err != nil {
			return err
		}
		if getConfig.ActiveVPN == expectedState {
			break
		}
		time.Sleep(100 * time.Millisecond)
	}
	return nil
}

func (x *RoutesXrayAPI) switchModeToTag(listMode string) string {
	var outboundTag string
	switch listMode {
	case "whitelist":
		outboundTag = "direct"
	case "blacklist":
		outboundTag = "proxy"
	default:
		outboundTag = "proxy"
	}

	return outboundTag
}

func (x *RoutesXrayAPI) ActualizeConfig() {
	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		logger.Error("Failed to get config", zap.Error(err))
		return
	}

	getRoutes, err := x.rr.GetRoutes(getConfig.ListMode)
	if err != nil {
		logger.Error("Error fetching routes", zap.String("listMode", getConfig.ListMode), zap.Error(err))
		return
	}

	*config.Xray.Routing = x.convertToRoutingConfig(getRoutes)

	// ВРЕМЯНКА - НАДО ПИЛИТЬ БД
	outbound := models.OutboundConfig{
		Protocol: "vless",
		Tag:      "proxy",
		Settings: map[string]interface{}{
			"vnext": []map[string]interface{}{
				{
					"address":      "5.252.178.253",
					"country_code": "RO",
					"port":         443,
					"users": []map[string]interface{}{
						{
							"encryption": "none",
							"flow":       "",
							"id":         "eca25c92-e209-4e0c-acba-2364948d2b60",
						},
					},
				},
			},
		},
		StreamSettings: map[string]interface{}{
			"network": "tcp",
			"realitySettings": map[string]interface{}{
				"fingerprint": "random",
				"publicKey":   "ebOVmspzPxXxK05suE8N81pVMfmDh4y8wvm_l5VSPik",
				"serverName":  "twitch.tv",
				"shortId":     "5f63",
				"spiderX":     "/",
			},
			"security": "reality",
		},
	}
	config.Xray.Outbounds = append(config.Xray.Outbounds, outbound)
	if getConfig.ListMode == "whitelist" || getConfig.DisableRoutes == true || len(getRoutes.Rules) == 0 {
		err = x.SwapOutbounds(&config.Xray.Outbounds, "proxy", "direct")
		if err != nil {
			logger.Error("Error swapping outbound rules", zap.Error(err))
		}
	}

	err = config.Save()
	if err != nil {
		logger.Error("Failed to save config", zap.Error(err))
		return
	}
}
