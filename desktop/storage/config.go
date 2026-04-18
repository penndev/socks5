// Package storage 提供 Bolt 持久化及与前端 JSON 字段一致的数据模型。
package storage

// Settings 与前端 settings store 持久化 JSON 一致。
type Settings struct {
	Proxy  ProxySettings  `json:"proxy"`
	System SystemSettings `json:"system"`
}

type ProxySettings struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SystemSettings struct {
	Language           string `json:"language"`
	ThemeMode          string `json:"themeMode"`
	StartupOnBoot      bool   `json:"startupOnBoot"`
	EnableLogRecording bool   `json:"enableLogRecording"`
}

// ServerEntry 与前端服务器列表项一致（不含运行时延迟字段）。
type ServerEntry struct {
	ID       string `json:"id"`
	Host     string `json:"host"`
	Remark   string `json:"remark"`
	Username string `json:"username"`
	Password string `json:"password"`
	Protocol string `json:"protocol"`
}

const (
	KeySettings = "settings"
	KeyServers  = "servers"
)
