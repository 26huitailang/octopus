package agent

import (
	"bufio"
	"bytes"

	"github.com/26huitailang/octopus/app/provider/agent"
	agentService "github.com/26huitailang/octopus/app/provider/agent"
	"github.com/26huitailang/yogo/framework/gin"
)

type AgentApi struct {
}

func Register(r *gin.Engine) error {
	api := NewAgentApi()
	r.Bind(&agentService.AgentProvider{})

	r.POST("/agent/script", api.ApiAgentScript)
	return nil
}

func NewAgentApi() *AgentApi {
	return &AgentApi{}
}

// ApiAgentScirpt run a command post to server
// @Summary run a command post to server
// @Description run a command post to server
// @Produce json
// @Tags agent
// @Success 200 json map[string]interface{}
// @Router /agent/script [post]
func (api *AgentApi) ApiAgentScript(c *gin.Context) {
	type Req struct {
		Script string
	}
	req := &Req{}
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithError(500, err)
	}
	agentService := c.MustMake(agent.AgentKey).(agent.IService)
	buffer := bytes.NewBuffer([]byte(""))
	writer := bufio.NewWriter(buffer)
	if err = agentService.RunScript([]byte(req.Script), writer); err != nil {
		c.AbortWithError(500, err)
	}
	writer.Flush()
	c.JSON(200, map[string]interface{}{"out": buffer.String()})
}
