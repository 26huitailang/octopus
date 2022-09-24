package agent

import (
	"bufio"
	"bytes"
	"os"
	"path"

	"github.com/26huitailang/octopus/app/provider/agent"
	agentService "github.com/26huitailang/octopus/app/provider/agent"
	"github.com/26huitailang/yogo/framework/gin"
	"github.com/26huitailang/yogo/framework/util"
)

type AgentApi struct {
}

func Register(r *gin.Engine) error {
	api := NewAgentApi()
	r.Bind(&agentService.AgentProvider{})

	r.POST("/agent/script/run", api.ApiAgentRunScript)
	r.GET("/agent/scripts", api.ApiAgentScripts)
	r.GET("/agent/scripts/:name", api.ApiAgentDetailScript)
	return nil
}

func NewAgentApi() *AgentApi {
	return &AgentApi{}
}

// ApiAgentScirpts list scripts
// @Summary list scripts
// @Description list scripts
// @Produce json
// @Tags agent
// @Success 200 json map[string]interface{}
// @Router /agent/scripts [get]
func (api *AgentApi) ApiAgentScripts(c *gin.Context) {
	agentService := c.MustMake(agent.AgentKey).(agent.IService)
	scripts, err := agentService.ListScript()
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(200, map[string]interface{}{"scripts": scripts})
}

// ApiAgentDetailScript get scirpt info
// @Summary get scirpt info
// @Description get scirpt info
// @Produce json
// @Tags agent
// @Success 200 json map[string]interface{}
// @Router /agent/scripts/:name [get]
func (api *AgentApi) ApiAgentDetailScript(c *gin.Context) {
	name := c.Param("name")
	appService := c.MustMakeApp()
	scriptPath := path.Join(appService.StorageFolder(), "scripts")
	file := path.Join(scriptPath, name)
	if !util.Exists(file) {
		c.JSON(400, map[string]interface{}{"detail": "not exists"})
	}

	content, err := os.ReadFile(file)
	if err != nil {
		c.AbortWithError(500, err)
	}
	c.JSON(200, map[string]interface{}{"content": string(content)})
}

// ApiAgentRunScript run a command post to server
// @Summary run a command post to server
// @Description run a command post to server
// @Produce json
// @Tags agent
// @Success 200 json map[string]interface{}
// @Router /agent/script/run [post]
func (api *AgentApi) ApiAgentRunScript(c *gin.Context) {
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
