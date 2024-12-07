package models

type Xray struct {
	Log LogConfig `json:"log"`
	API APIConfig `json:"api"`
	//	DNS         DNSConfig        `json:"dns"`
	Inbounds  []InboundConfig  `json:"inbounds"`
	Outbounds []OutboundConfig `json:"outbounds"`
	Policy    PolicyConfig     `json:"policy"`
	//	Reverse   ReverseConfig    `json:"reverse"`
	Routing *RoutingConfig `json:"routing,omitempty"`
	//	Transport   interface{}      `json:"transport"`
	Stats interface{} `json:"stats"`
	//	FakeDNS     interface{}      `json:"fakedns"`
	//	Metrics     interface{}      `json:"metrics"`
	//	Observatory interface{}      `json:"observatory"`
	//	BurstObs    interface{}      `json:"burstObservatory"`
}

type LogConfig struct {
	//	Access      string `json:"access"`
	//	Error       string `json:"error"`
	LogLevel string `json:"loglevel"`
	//	DnsLog      bool   `json:"dnsLog"`
	//	MaskAddress string `json:"maskAddress"`
}

type APIConfig struct {
	Tag      string   `json:"tag"`
	Listen   string   `json:"listen"`
	Services []string `json:"services"`
}

type DNSConfig struct {
	Hosts                  map[string]interface{} `json:"hosts"`
	Servers                []DNSServerConfig      `json:"servers"`
	ClientIP               string                 `json:"clientIp"`
	QueryStrategy          string                 `json:"queryStrategy"`
	DisableCache           bool                   `json:"disableCache"`
	DisableFallback        bool                   `json:"disableFallback"`
	DisableFallbackIfMatch bool                   `json:"disableFallbackIfMatch"`
	Tag                    string                 `json:"tag"`
}

type DNSServerConfig struct {
	Address      string   `json:"address"`
	Port         *int     `json:"port,omitempty"`
	Domains      []string `json:"domains"`
	ExpectIPs    []string `json:"expectIPs"`
	SkipFallback bool     `json:"skipFallback"`
	ClientIP     string   `json:"clientIP,omitempty"`
}

type InboundConfig struct {
	Listen         string                 `json:"listen"`
	Port           int                    `json:"port"`
	Protocol       string                 `json:"protocol"`
	Settings       map[string]interface{} `json:"settings,omitempty"`
	StreamSettings map[string]interface{} `json:"streamSettings,omitempty"`
	Tag            string                 `json:"tag,omitempty"`
	Sniffing       SniffingConfig         `json:"sniffing,omitempty"`
	Allocate       AllocateConfig         `json:"allocate,omitempty"`
}

type SniffingConfig struct {
	Enabled      bool     `json:"enabled,omitempty"`
	DestOverride []string `json:"destOverride,omitempty"`
	RouteOnly    bool     `json:"routeOnly,omitempty"`
}

type AllocateConfig struct {
	Strategy    string `json:"strategy,omitempty"`
	Refresh     int    `json:"refresh,omitempty"`
	Concurrency int    `json:"concurrency,omitempty"`
}

type OutboundConfig struct {
	SendThrough    string                 `json:"sendThrough,omitempty"`
	Protocol       string                 `json:"protocol"`
	Settings       map[string]interface{} `json:"settings,omitempty"`
	Tag            string                 `json:"tag"`
	StreamSettings map[string]interface{} `json:"streamSettings,omitempty"`
	//ProxySettings  ProxySettingsConfig    `json:"proxySettings,omitempty"`
	Mux interface{} `json:"mux,omitempty"`
}

type ProxySettingsConfig struct {
	Tag string `json:"tag,omitempty"`
}

type PolicyConfig struct {
	Levels map[string]PolicyLevelConfig `json:"levels,omitempty"`
	System PolicySystemConfig           `json:"system,omitempty"`
}

type PolicyLevelConfig struct {
	Handshake         int  `json:"handshake,omitempty"`
	ConnIdle          int  `json:"connIdle,omitempty"`
	UplinkOnly        int  `json:"uplinkOnly,omitempty"`
	DownlinkOnly      int  `json:"downlinkOnly,omitempty"`
	StatsUserUplink   bool `json:"statsUserUplink,omitempty"`
	StatsUserDownlink bool `json:"statsUserDownlink,omitempty"`
	BufferSize        int  `json:"bufferSize,omitempty"`
}

type PolicySystemConfig struct {
	StatsInboundUplink    bool `json:"statsInboundUplink,omitempty"`
	StatsInboundDownlink  bool `json:"statsInboundDownlink,omitempty"`
	StatsOutboundUplink   bool `json:"statsOutboundUplink,omitempty"`
	StatsOutboundDownlink bool `json:"statsOutboundDownlink,omitempty"`
}

type ReverseConfig struct {
	Bridges []ReverseBridgeConfig `json:"bridges,omitempty"`
	Portals []ReversePortalConfig `json:"portals,omitempty"`
}

type ReverseBridgeConfig struct {
	Tag    string `json:"tag,omitempty"`
	Domain string `json:"domain,omitempty"`
}

type ReversePortalConfig struct {
	Tag    string `json:"tag,omitempty"`
	Domain string `json:"domain,omitempty"`
}

type RoutingConfig struct {
	DomainStrategy string        `json:"domainStrategy"`
	DomainMatcher  string        `json:"domainMatcher"`
	Rules          []RoutingRule `json:"rules,omitempty"`
	Balancers      []string      `json:"balancers,omitempty"`
}

type RoutingRule struct {
	Type        string   `json:"type"`
	Domain      []string `json:"domain,omitempty"`
	IP          []string `json:"ip,omitempty"`
	Port        string   `json:"port,omitempty"`
	OutboundTag string   `json:"outboundTag"`
}
