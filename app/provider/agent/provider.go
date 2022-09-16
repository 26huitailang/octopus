package agent

import (
	"github.com/26huitailang/yogo/framework"
)

type AgentProvider struct {
	framework.ServiceProvider

	c framework.Container
}

func (sp *AgentProvider) Name() string {
	return AgentKey
}

func (sp *AgentProvider) Register(c framework.Container) framework.NewInstance {
	return NewService
}

func (sp *AgentProvider) IsDefer() bool {
	return true
}

func (sp *AgentProvider) Params(c framework.Container) []interface{} {
	return []interface{}{sp.c}
}

func (sp *AgentProvider) Boot(c framework.Container) error {
	sp.c = c
	return nil
}
