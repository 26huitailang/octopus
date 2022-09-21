package agent

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/26huitailang/yogo/framework"
	"github.com/26huitailang/yogo/framework/contract"
	"github.com/google/shlex"
)

type Service struct {
	c framework.Container
	IService
}

// 初始化实例的方法
func NewService(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)
	return &Service{c: c}, nil
}

func (s *Service) Register() error {
	appService := s.c.MustMake(contract.AppKey).(contract.App)
	baseFolder := appService.BaseFolder()
	binPath := path.Join(baseFolder, "octopus")

	tmplData := struct {
		BaseFolder string
		BinPath    string
	}{
		BaseFolder: baseFolder,
		BinPath:    binPath,
	}
	// 检查是否有systemd服务
	systemctlPath, err := exec.LookPath("systemctl")
	if err != nil {
		fmt.Println("no systemctl in this machine")
		return err
	}
	// stop服务
	cmd := exec.Command(systemctlPath, "stop", "octopus-agent.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("=========== systemctl stop failed ===========")
		fmt.Println(string(out))
		fmt.Println("=========== systemctl stop failed ===========")
		return err
	}
	// 生成配置
	path := "/etc/systemd/system/octopus-agent.service"
	f, err := os.Create(path)
	if err != nil {
		fmt.Printf("open %s failed: %v", path, err)
		return err
	}
	defer f.Close()
	t := template.Must(template.New("systemctl").Parse(systemdServiceTmpl))
	if err = t.Execute(f, tmplData); err != nil {
		fmt.Println("generate octopus-agent.service file failed")
		return err
	}
	// daemon-reload 配置
	cmd = exec.Command(systemctlPath, "daemon-reload")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("=========== systemctl aemon-reload failed ===========")
		fmt.Println(string(out))
		fmt.Println("=========== systemctl aemon-reload failed ===========")
		return err
	}
	// restart服务
	cmd = exec.Command(systemctlPath, "restart", "octopus-agent.service")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("=========== systemctl start failed ===========")
		fmt.Println(string(out))
		fmt.Println("=========== systemctl start failed ===========")
		return err
	}
	// enable服务
	cmd = exec.Command(systemctlPath, "enable", "octopus-agent.service")
	out, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println("=========== systemctl enable failed ===========")
		fmt.Println(string(out))
		fmt.Println("=========== systemctl enable failed ===========")
		return err
	}
	fmt.Println("register octopus-agent.serivce sucess")
	return nil
}

func (s *Service) Status() error {
	// 检查是否有systemd服务
	systemctlPath, err := exec.LookPath("systemctl")
	if err != nil {
		fmt.Println("no systemctl in this machine")
		return err
	}
	// stop服务
	cmd := exec.Command(systemctlPath, "status", "octopus-agent.service")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("=========== systemctl stop failed ===========")
		fmt.Println(string(out))
		fmt.Println("=========== systemctl stop failed ===========")
		return err
	}
	fmt.Println(string(out))
	return nil
}

func (s *Service) RunScript(script []byte, writer io.Writer) error {
	logger := s.c.MustMake(contract.LogKey).(contract.Log)
	logger.Debug(context.TODO(), "got script, start parse", map[string]interface{}{"scirpt": string(script)})
	cmds := bytes.Split(script, []byte("\n"))
	for _, cmd := range cmds {
		ret, err := shlex.Split(string(cmd))
		if err != nil {
			return err
		}
		var command *exec.Cmd
		logger.Debug(context.TODO(), "get command", map[string]interface{}{"ret": ret})
		if len(ret) == 0 {
			continue
		} else if len(ret) == 1 {
			command = exec.Command(ret[0])
		} else {
			command = exec.Command(ret[0], ret[1:]...)
		}
		out, err := command.CombinedOutput()
		if err != nil {
			logger.Error(context.TODO(), "out:", map[string]interface{}{"out": string(out), "cmd": cmd})
			return err
		}
		n, err := writer.Write(out)
		logger.Debug(context.TODO(), "out:", map[string]interface{}{"out": string(out), "cmd": string(cmd), "n": n})
		if err != nil {
			return err
		}
	}
	return nil
}

var systemdServiceTmpl = `
[Unit]
Description=development agent
After=network.target
StartLimitIntervalSec=0
[Service]
Type=simple
Restart=always
RestartSec=5
User=root
WorkingDirectory={{.BaseFolder}}
ExecStart={{.BinPath}} app start

[Install]
WantedBy=multi-user.target
`
