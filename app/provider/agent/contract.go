package agent

import "io"

const AgentKey = "yogo:agent"

type IService interface {
	Register() error // 注册agent为本机服务
	Status() error   // 查询agent服务状态
	RunScript(script []byte, writer io.Writer) error
	Tail(path string, writer io.Writer) error
}

type Agent struct {
	Ip      string // 本机ip
	Port    string // 本机端口
	HostUrl string // 主机地址
	Status  bool   // 状态
}
