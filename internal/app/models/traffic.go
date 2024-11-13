package models

type StatsTrafficJSON struct {
	Stat []Traffic `json:"stat"`
}
type Traffic struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type StatsTraffic struct {
	ProxyUplink    int64
	ProxyDownlink  int64
	DirectUplink   int64
	DirectDownlink int64
}
