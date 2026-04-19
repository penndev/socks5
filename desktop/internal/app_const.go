package internal

type appConst struct {
	// 应用运行状态日志名称
	LogTypeName_STATUS string
	// 连接日志列表的消息名称
	LogTypeName_LOG string
	// 语言切换事件（payload 为语言标识字符串）
	EventNameLocaleChanged string
}

var AppConfig = appConst{
	LogTypeName_STATUS:     "logServerStatus",
	LogTypeName_LOG:        "logProxyList",
	EventNameLocaleChanged: "localeChanged",
}

type AppConst struct{}

func (a AppConst) AppConfig() appConst {
	return AppConfig
}

func (a AppConst) ProxyScheme() []string {
	return []string{
		"Socks5",
		"Socks5OverTLS",
		"Http",
		"HttpOverTLS",
	}
}
