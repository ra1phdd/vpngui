package xray_api

import (
	"fmt"
	"regexp"
	"strings"
	"time"
	"vpngui/config"
	"vpngui/internal/app/models"
)

func (x *XrayAPI) DisableRoutes() error {
	config.Routes = *config.Xray.Routing
	config.Xray.Routing = nil
	config.JSON.DisableRoutes = true

	err := x.SwapOutbounds(&config.Xray.Outbounds, "proxy", "direct")
	if err != nil {
		return err
	}

	err = config.SaveConfig()
	if err != nil {
		return err
	}

	x.Kill()
	for config.JSON.ActiveVPN {
		time.Sleep(100 * time.Millisecond)
	}
	go x.Run()
	for !config.JSON.ActiveVPN {
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (x *XrayAPI) EnableRoutes() error {
	if config.Xray.Routing == nil {
		config.Xray.Routing = new(models.RoutingConfig)
	}

	*config.Xray.Routing = config.Routes

	config.Routes = models.RoutingConfig{}
	config.JSON.DisableRoutes = false

	if config.JSON.EnableBlackList {
		err := x.SwapOutbounds(&config.Xray.Outbounds, "direct", "proxy")
		if err != nil {
			return err
		}
	} else {
		err := x.SwapOutbounds(&config.Xray.Outbounds, "proxy", "direct")
		if err != nil {
			return err
		}
	}

	err := config.SaveConfig()
	if err != nil {
		return err
	}

	x.Kill()
	for config.JSON.ActiveVPN {
		time.Sleep(100 * time.Millisecond)
	}
	go x.Run()
	for !config.JSON.ActiveVPN {
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

func (x *XrayAPI) SwapOutbounds(outbounds *[]models.OutboundConfig, tag1, tag2 string) error {
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
func (x *XrayAPI) GetDomain(outboundTag string) string {
	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			return strings.Join(config.Xray.Routing.Rules[i].Domain, "\n")
		}
	}

	return ""
}

func (x *XrayAPI) GetIP(outboundTag string) string {
	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			return strings.Join(config.Xray.Routing.Rules[i].IP, "\n")
		}
	}

	return ""
}

func (x *XrayAPI) GetPort(outboundTag string) string {
	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			return config.Xray.Routing.Rules[i].Port
		}
	}

	return ""
}

func (x *XrayAPI) AddDomain(outboundTag string, domain string) {
	isValidTag(outboundTag)

	if !isValidDomain(domain) {
		fmt.Println("Невалидный домен")
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			config.Xray.Routing.Rules[i].Domain = append(config.Xray.Routing.Rules[i].Domain, domain)
			break
		}
	}

	err := config.SaveConfig()
	if err != nil {
		return
	}
}

func (x *XrayAPI) AddIP(outboundTag string, ip string) {
	isValidTag(outboundTag)

	if !isValidIP(ip) {
		fmt.Println("Невалидный IP-адрес")
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			config.Xray.Routing.Rules[i].IP = append(config.Xray.Routing.Rules[i].IP, ip)
			break
		}
	}

	err := config.SaveConfig()
	if err != nil {
		return
	}
}

func (x *XrayAPI) AddPort(outboundTag string, port string) {
	isValidTag(outboundTag)

	if !isValidPort(port) {
		fmt.Println("Невалидный порт")
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			if config.Xray.Routing.Rules[i].Port == port {
				break
			} else if config.Xray.Routing.Rules[i].Port == "" {
				config.Xray.Routing.Rules[i].Port = port
				break
			}

			config.Xray.Routing.Rules[i].Port += fmt.Sprintf(",%s", port)
			break
		}
	}

	err := config.SaveConfig()
	if err != nil {
		return
	}
}

func (x *XrayAPI) DelDomain(outboundTag string, domain string) {
	isLastItem(outboundTag)
	if !isValidDomain(domain) {
		fmt.Println("Невалидный домен")
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			for j, d := range config.Xray.Routing.Rules[i].Domain {
				if d == domain {
					config.Xray.Routing.Rules[i].Domain = append(config.Xray.Routing.Rules[i].Domain[:j], config.Xray.Routing.Rules[i].Domain[j+1:]...)
					break
				}
			}
			break
		}
	}

	err := config.SaveConfig()
	if err != nil {
		return
	}
}

func (x *XrayAPI) DelIP(outboundTag string, ip string) {
	isLastItem(outboundTag)
	if !isValidIP(ip) {
		fmt.Println("Невалидный IP-адрес")
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			for j, ipAddr := range config.Xray.Routing.Rules[i].IP {
				if ipAddr == ip {
					config.Xray.Routing.Rules[i].IP = append(config.Xray.Routing.Rules[i].IP[:j], config.Xray.Routing.Rules[i].IP[j+1:]...)
					break
				}
			}
			break
		}
	}

	err := config.SaveConfig()
	if err != nil {
		return
	}
}

func (x *XrayAPI) DelPort(outboundTag string, port string) {
	isLastItem(outboundTag)
	if !isValidPort(port) {
		fmt.Println("Невалидный порт")
		return
	}

	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			ports := config.Xray.Routing.Rules[i].Port
			portList := strings.Split(ports, ",")

			for j, p := range portList {
				if p == port {
					portList = append(portList[:j], portList[j+1:]...)
					break
				}
			}

			config.Xray.Routing.Rules[i].Port = strings.Join(portList, ",")
			break
		}
	}

	err := config.SaveConfig()
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

func isValidTag(outboundTag string) {
	found := false
	for i := range config.Xray.Routing.Rules {
		if config.Xray.Routing.Rules[i].OutboundTag == outboundTag {
			found = true
		}
	}

	if !found {
		newRule := models.RoutingRule{
			Type:        "field",
			OutboundTag: outboundTag,
		}

		config.Xray.Routing.Rules = append(config.Xray.Routing.Rules, newRule)
	}
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
