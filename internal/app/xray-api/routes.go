package xray_api

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"vpngui/config"
	"vpngui/internal/app/models"
	"vpngui/internal/app/repository"
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
	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		return err
	}

	config.Xray.Routing = nil

	if err := x.run.cr.UpdateDisableRoutes(true); err != nil {
		return err
	}

	if err := x.SwapOutbounds(&config.Xray.Outbounds, "proxy", "direct"); err != nil {
		return err
	}

	if err := config.SaveConfig(); err != nil {
		return err
	}

	return x.restartVPNIfActive(getConfig.ActiveVPN)
}

func (x *RoutesXrayAPI) EnableRoutes() error {
	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		return err
	}

	if err := x.updateRoutingConfig(getConfig.ListMode); err != nil {
		return err
	}

	if err := x.run.cr.UpdateDisableRoutes(false); err != nil {
		return err
	}

	swapDirection := "direct"
	if getConfig.ListMode == "blacklist" {
		swapDirection = "proxy"
	}
	if err := x.SwapOutbounds(&config.Xray.Outbounds, swapDirection, "direct"); err != nil {
		return err
	}

	if err := config.SaveConfig(); err != nil {
		return err
	}

	return x.restartVPNIfActive(getConfig.ActiveVPN)
}

func (x *RoutesXrayAPI) EnableBlackList() error {
	return x.toggleListMode("blacklist", "direct", "proxy")
}

func (x *RoutesXrayAPI) DisableBlackList() error {
	return x.toggleListMode("whitelist", "proxy", "direct")
}

func (x *RoutesXrayAPI) toggleListMode(listMode, outbound1, outbound2 string) error {
	getConfig, err := x.run.cr.GetConfig()
	if err != nil {
		return err
	}

	if err := x.updateRoutingConfig(listMode); err != nil {
		return err
	}

	if err := x.run.cr.UpdateListMode(listMode); err != nil {
		return err
	}

	if err := x.SwapOutbounds(&config.Xray.Outbounds, outbound1, outbound2); err != nil {
		return err
	}

	if err := config.SaveConfig(); err != nil {
		return err
	}

	return x.restartVPNIfActive(getConfig.ActiveVPN)
}

func (x *RoutesXrayAPI) updateRoutingConfig(listMode string) error {
	if config.Xray.Routing == nil {
		config.Xray.Routing = new(models.RoutingConfig)
	}

	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		return err
	}
	*config.Xray.Routing = x.convertToRoutingConfig(getRoutes)

	return nil
}

func (x *RoutesXrayAPI) restartVPNIfActive(active bool) error {
	if !active {
		return nil
	}

	x.run.Kill()
	if err := x.waitForVPNState(false); err != nil {
		return err
	}

	go x.run.Run()
	if err := x.waitForVPNState(false); err != nil {
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
		return fmt.Errorf("one or both tags not found")
	} else if index1 == 0 && index2 == 1 {
		return nil
	}

	(*outbounds)[index1], (*outbounds)[index2] = (*outbounds)[index2], (*outbounds)[index1]

	return nil
}

func (x *RoutesXrayAPI) GetDomain(listMode string) string {
	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		return ""
	}

	var domains []string
	for _, rule := range getRoutes.Rules {
		if rule.RuleType == "domain" {
			domains = append(domains, rule.RuleValue)
		}
	}

	return strings.Join(domains, "\n")
}

func (x *RoutesXrayAPI) GetIP(listMode string) string {
	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		return ""
	}

	var ips []string
	for _, rule := range getRoutes.Rules {
		if rule.RuleType == "ip" {
			ips = append(ips, rule.RuleValue)
		}
	}

	return strings.Join(ips, "\n")
}

func (x *RoutesXrayAPI) GetPort(listMode string) string {
	getRoutes, err := x.rr.GetRoutes(listMode)
	if err != nil {
		return ""
	}

	var ports []string
	for _, rule := range getRoutes.Rules {
		if rule.RuleType == "port" {
			ports = append(ports, rule.RuleValue)
		}
	}

	return strings.Join(ports, ", ")
}

func (x *RoutesXrayAPI) AddDomain(listMode string, domain string) {
	if !isValidDomain(domain) {
		fmt.Println("Невалидный домен")
		return
	}

	err := x.rr.AddRule(listMode, "domain", domain)
	if err != nil {
		return
	}
}

func (x *RoutesXrayAPI) AddIP(listMode string, ip string) {
	if !isValidIP(ip) {
		fmt.Println("Невалидный IP-адрес")
		return
	}

	err := x.rr.AddRule(listMode, "ip", ip)
	if err != nil {
		return
	}
}

func (x *RoutesXrayAPI) AddPort(listMode string, port string) {
	if !isValidPort(port) {
		fmt.Println("Невалидный IP-адрес")
		return
	}

	err := x.rr.AddRule(listMode, "port", port)
	if err != nil {
		return
	}
}

func (x *RoutesXrayAPI) DelDomain(listMode string, domain string) {
	isLastItem(listMode)
	if !isValidDomain(domain) {
		fmt.Println("Невалидный домен")
		return
	}

	err := x.rr.DeleteRule(listMode, "domain", domain)
	if err != nil {
		return
	}
}

func (x *RoutesXrayAPI) DelIP(listMode string, ip string) {
	isLastItem(listMode)
	if !isValidIP(ip) {
		fmt.Println("Невалидный IP-адрес")
		return
	}

	err := x.rr.DeleteRule(listMode, "ip", ip)
	if err != nil {
		return
	}
}

func (x *RoutesXrayAPI) DelPort(listMode string, port string) {
	isLastItem(listMode)
	if !isValidPort(port) {
		fmt.Println("Невалидный порт")
		return
	}

	err := x.rr.DeleteRule(listMode, "port", port)
	if err != nil {
		return
	}
}

func isValidDomain(domain string) bool {
	regex := `^([a-zA-Z0-9-]+\.)+[a-zA-Z]{2,63}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(domain)
}

func isValidIP(domain string) bool {
	regex := `^((25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])\.){3}(25[0-5]|2[0-4][0-9]|1[0-9]{2}|[1-9]?[0-9])$`
	re := regexp.MustCompile(regex)
	return re.MatchString(domain)
}

func isValidPort(domain string) bool {
	regex := `^(6553[0-5]|655[0-2][0-9]|65[0-4][0-9]{2}|6[0-4][0-9]{3}|[1-5][0-9]{4}|[1-9][0-9]{0,3})$`
	re := regexp.MustCompile(regex)
	return re.MatchString(domain)
}

func isLastItem(outboundTag string) {
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

	err := config.SaveConfig()
	if err != nil {
		return
	}
}

func (x *RoutesXrayAPI) convertToRoutingConfig(listConfig models.ListConfig) models.RoutingConfig {
	var routingConfig models.RoutingConfig

	routingConfig.DomainStrategy = listConfig.DomainStrategy
	routingConfig.DomainMatcher = listConfig.DomainMatcher

	outboundTag := "proxy"
	if listConfig.Type == "whitelist" {
		outboundTag = "direct"
	}

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
