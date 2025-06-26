package blemodule

import (
	"fmt"
	"main/interfaces"
	"strings"
	"time"
)

// atcmd.go 封装所有 AT 命令的发送与响应处理
// 参考 HCM111Z AT 手册，便于维护和查找

type ATCommander struct {
	port interfaces.SerialPort
}

func NewATCommander(port interfaces.SerialPort) *ATCommander {
	return &ATCommander{port: port}
}

// 发送 AT 命令并等待响应
func (a *ATCommander) Send(cmd string, timeout time.Duration) (string, error) {
	fullCmd := fmt.Sprintf("%s\r\n", cmd)
	err := a.port.Write(fullCmd)
	if err != nil {
		return "", err
	}
	resp, err := a.port.ReadLine(timeout)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}

// TODO: 按照手册继续添加各类命令的封装方法
