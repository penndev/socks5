package main

type AppConst struct {
	// 应用运行状态日志名称
	LogTypeName_STATUS string
	// 连接日志列表的消息名称
	LogTypeName_LOG string
}

var appConst = AppConst{
	LogTypeName_STATUS: "logServerStatus",
	LogTypeName_LOG:    "logProxyList",
}

type App struct{}

func (a App) AppConfig() AppConst {
	return appConst
}
