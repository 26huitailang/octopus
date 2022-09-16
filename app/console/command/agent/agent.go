package agent

import (
	"github.com/26huitailang/octopus/app/provider/agent"
	"github.com/26huitailang/yogo/framework/cobra"
)

func InitAgentCommand() *cobra.Command {
	AgentCommand.AddCommand(AgentRegisterLocalCommand)
	AgentCommand.AddCommand(AgentStatusCommand)
	return AgentCommand
}

var AgentCommand = &cobra.Command{
	Use:   "agent",
	Short: "agent",
	RunE: func(c *cobra.Command, args []string) error {
		if len(args) == 0 {
			c.Help()
		}
		return nil
	},
}

var AgentRegisterLocalCommand = &cobra.Command{
	Use:   "register",
	Short: "register to host as a service",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		agentService := container.MustMake(agent.AgentKey).(agent.IService)
		if err := agentService.Register(); err != nil {
			return err
		}
		return nil
	},
}

var AgentStatusCommand = &cobra.Command{
	Use:   "status",
	Short: "status of octopus-agent.service",
	RunE: func(cmd *cobra.Command, args []string) error {
		container := cmd.GetContainer()
		agentService := container.MustMake(agent.AgentKey).(agent.IService)
		if err := agentService.Status(); err != nil {
			return err
		}
		return nil
	},
}
