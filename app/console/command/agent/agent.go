package agent

import (
	"os"

	"github.com/26huitailang/octopus/app/provider/agent"
	"github.com/26huitailang/yogo/framework/cobra"
)

func InitAgentCommand() *cobra.Command {
	AgentCommand.AddCommand(AgentRegisterLocalCommand)
	AgentCommand.AddCommand(AgentStatusCommand)
	AgentCommand.AddCommand(AgentRunScriptCommand)
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

var AgentRunScriptCommand = &cobra.Command{
	Use:   "script",
	Short: "run a specific script in exec.Command, ./yogo agent script ./demo.sh",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			cmd.Help()
			return nil
		}
		path := args[0]
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		container := cmd.GetContainer()
		writer := os.Stdout
		agentService := container.MustMake(agent.AgentKey).(agent.IService)
		if err := agentService.RunScript(content, writer); err != nil {
			return err
		}
		return nil
	},
}
