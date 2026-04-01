package internal

type appConst struct {
	// 应用运行状态日志名称
	LogTypeName_STATUS string
	// 连接日志列表的消息名称
	LogTypeName_LOG string
}

var AppConfig = appConst{
	LogTypeName_STATUS: "logServerStatus",
	LogTypeName_LOG:    "logProxyList",
}

type AppConst struct{}

func (a AppConst) AppConfig() appConst {
	return AppConfig
}
