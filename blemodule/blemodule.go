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
// 包含子模块封装接口：模块控制、广播、GATT、连接、透传等
type BLEModule struct {
	Port       interfaces.SerialPort
	urcHandler func(urc string) // URC 处理函数（可选）
}

// New 创建 BLE 模块控制器
func New(port interfaces.SerialPort) *BLEModule {
	return &BLEModule{Port: port}
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

func (b *BLEModule) Restart() error {
	return b.sendAT("AT+QRST")
}

func (b *BLEModule) GetVersion() (string, error) {
	return b.sendATRead("AT+QVERSION")
}

func (b *BLEModule) SetBaud(baud int) error {
	return b.sendAT(fmt.Sprintf("AT+QSETBAUD=%d", baud))
}

// --- BLE 初始化与配置 ---

func (b *BLEModule) Init(role int) error {
	return b.sendAT(fmt.Sprintf("AT+QBLEINIT=%d", role))
}

func (b *BLEModule) SetDeviceName(name string) error {
	return b.sendAT(fmt.Sprintf("AT+QBLENAME=\"%s\"", name))
}

func (b *BLEModule) QueryAddress() (string, error) {
	return b.sendATRead("AT+QBLEADDR?")
}

// --- 广播控制 ---

func (b *BLEModule) StartAdvertising() error {
	return b.sendAT("AT+QBLEADVSTART")
}

func (b *BLEModule) StopAdvertising() error {
	return b.sendAT("AT+QBLEADVSTOP")
}

// --- GATT 服务端 ---

func (b *BLEModule) AddService(uuid string) error {
	return b.sendAT(fmt.Sprintf("AT+QBLEGATTSSRV=\"%s\"", uuid))
}

func (b *BLEModule) AddCharacteristic(uuid string, prop int) error {
	return b.sendAT(fmt.Sprintf("AT+QBLEGATTSCHAR=\"%s\",%d", uuid, prop))
}

func (b *BLEModule) FinishGATTServer() error {
	return b.sendAT("AT+QBLEGATTSSRVDONE")
}

func (b *BLEModule) SendNotify(connIdx, handle int, value string) error {
	return b.sendAT(fmt.Sprintf("AT+QBLEGATTSNTFY=%d,%d,\"%s\"", connIdx, handle, value))
}

// --- 透传模式 ---

func (b *BLEModule) EnterTransparent(uuid string) error {
	return b.sendAT(fmt.Sprintf("AT+QBLETRANMODE=\"%s\"", uuid))
}

func (b *BLEModule) ExitTransparent() error {
	time.Sleep(1 * time.Second)
	err := b.Port.Write("+++")
	time.Sleep(1 * time.Second)
	return err
}

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
