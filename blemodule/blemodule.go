// Package blemodule 提供对移远 HCM111Z BLE 模块的封装接口
// 支持初始化、广播、GATT、透传、连接等核心功能
package blemodule

import (
	"fmt"
	"main/interfaces"
	"strings"
	"time"
)

// BLEModule 是 BLE 模块的主控制器
// 负责模块控制、BLE 初始化、广播、GATT、透传等功能
// 内部通过 ATCommander 发送 AT 命令
//
// 详细 API 参考 HCM111Z AT 命令手册
type BLEModule struct {
	Port       interfaces.SerialPort
	at         *ATCommander
	urcHandler func(urc string) // URC 处理函数（可选）
}

// New 创建 BLE 模块控制器
func New(port interfaces.SerialPort) *BLEModule {
	return &BLEModule{
		Port: port,
		at:   NewATCommander(port),
	}
}

// RegisterURCHandler 注册 URC 回调处理函数
func (b *BLEModule) RegisterURCHandler(handler func(urc string)) {
	b.urcHandler = handler
}

// StartURCListener 启动 URC 异步监听协程
func (b *BLEModule) StartURCListener() {
	go func() {
		for {
			line, err := b.Port.ReadLine(0)
			if err != nil || len(line) == 0 {
				continue
			}
			if strings.HasPrefix(line, "+QBLE") {
				if b.urcHandler != nil {
					b.urcHandler(line)
				} else {
					fmt.Println("[URC]", line)
				}
			}
		}
	}()
}

// --- 通用模块控制 ---
// 模块重启
func (b *BLEModule) Restart() error {
	_, err := b.at.Send("AT+QRST", 2*time.Second)
	return err
}

// 查询固件版本
func (b *BLEModule) GetVersion() (string, error) {
	return b.at.Send("AT+QVERSION", 2*time.Second)
}

// 设置串口波特率
func (b *BLEModule) SetBaud(baud int) error {
	_, err := b.at.Send(fmt.Sprintf("AT+QSETBAUD=%d", baud), 2*time.Second)
	return err
}

// --- BLE 初始化与配置 ---
// 初始化 BLE 栈
func (b *BLEModule) Init(role int) error {
	_, err := b.at.Send(fmt.Sprintf("AT+QBLEINIT=%d", role), 2*time.Second)
	return err
}

// 设置设备名称
func (b *BLEModule) SetDeviceName(name string) error {
	_, err := b.at.Send(fmt.Sprintf("AT+QBLENAME=\"%s\"", name), 2*time.Second)
	return err
}

// 查询 BLE MAC 地址
func (b *BLEModule) QueryAddress() (string, error) {
	return b.at.Send("AT+QBLEADDR?", 2*time.Second)
}

// --- 广播控制 ---
// 启动广播
func (b *BLEModule) StartAdvertising() error {
	_, err := b.at.Send("AT+QBLEADVSTART", 2*time.Second)
	return err
}

// 停止广播
func (b *BLEModule) StopAdvertising() error {
	_, err := b.at.Send("AT+QBLEADVSTOP", 2*time.Second)
	return err
}

// --- GATT 服务端 ---
// 添加服务
func (b *BLEModule) AddService(uuid string) error {
	_, err := b.at.Send(fmt.Sprintf("AT+QBLEGATTSSRV=\"%s\"", uuid), 2*time.Second)
	return err
}

// 添加特征
func (b *BLEModule) AddCharacteristic(uuid string, prop int) error {
	_, err := b.at.Send(fmt.Sprintf("AT+QBLEGATTSCHAR=\"%s\",%d", uuid, prop), 2*time.Second)
	return err
}

// 提交服务定义
func (b *BLEModule) FinishGATTServer() error {
	_, err := b.at.Send("AT+QBLEGATTSSRVDONE", 2*time.Second)
	return err
}

// 发送 Notify
func (b *BLEModule) SendNotify(connIdx, handle int, value string) error {
	_, err := b.at.Send(fmt.Sprintf("AT+QBLEGATTSNTFY=%d,%d,\"%s\"", connIdx, handle, value), 2*time.Second)
	return err
}

// --- 透传模式 ---
// 进入透传模式
func (b *BLEModule) EnterTransparent(uuid string) error {
	_, err := b.at.Send(fmt.Sprintf("AT+QBLETRANMODE=\"%s\"", uuid), 2*time.Second)
	return err
}

// 退出透传模式
func (b *BLEModule) ExitTransparent() error {
	time.Sleep(1 * time.Second)
	err := b.Port.Write("+++")
	time.Sleep(1 * time.Second)
	return err
}

// 发送透传数据
func (b *BLEModule) SendData(data string) error {
	return b.Port.Write(data)
}

func (b *BLEModule) SendLine(data string) error {
	return b.Port.Write(data + "\r\n")
}

// --- 私有方法 ---

func (b *BLEModule) sendAT(cmd string) error {
	return b.Port.Write(cmd + "\r\n")
}

func (b *BLEModule) sendATRead(cmd string) (string, error) {
	err := b.sendAT(cmd)
	if err != nil {
		return "", err
	}
	resp, err := b.Port.ReadLine(2 * time.Second)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp), nil
}
