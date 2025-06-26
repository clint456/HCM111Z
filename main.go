package main

import (
	"fmt"
	"log"
	"time"

	"main/blemodule"
	"main/serialport"
)

func main() {
	// === 0. 串口初始化 ===
	dev := "/dev/ttyS3" // 串口设备路径（请根据实际修改）
	baud := 115200
	port, err := serialport.OpenSerialPort(dev, baud)
	if err != nil {
		log.Fatalf("串口打开失败: %v", err)
	}
	ble := blemodule.New(port)

	// === 1. 通用控制命令演示 ===
	fmt.Println("[1] 模块重启...")
	if err := ble.Restart(); err != nil {
		log.Fatalf("Restart 失败: %v", err)
	}
	time.Sleep(2 * time.Second)

	fmt.Println("[2] 查询固件版本...")
	version, err := ble.GetVersion()
	if err != nil {
		log.Fatalf("获取版本失败: %v", err)
	}
	fmt.Println("版本信息:", version)

	// === 2. BLE 初始化与配置 ===
	fmt.Println("[3] 初始化 BLE peripheral...")
	if err := ble.Init(1); err != nil {
		log.Fatalf("Init 失败: %v", err)
	}

	fmt.Println("[4] 设置设备名称...")
	if err := ble.SetDeviceName("HCM111Z"); err != nil {
		log.Fatalf("SetName 失败: %v", err)
	}

	addr, err := ble.QueryAddress()
	if err != nil {
		log.Fatalf("获取地址失败: %v", err)
	}
	fmt.Println("设备地址:", addr)

	// === 3. 广播控制 ===
	fmt.Println("[5] 启动广播...请用手机 App 连接设备")
	if err := ble.StartAdvertising(); err != nil {
		log.Fatalf("StartAdvertising 失败: %v", err)
	}

	// 自动监听 URC，等待 +QBLESTAT: 1（已连接）
	connected := make(chan struct{})
	ble.RegisterURCHandler(func(urc string) {
		evt := blemodule.ParseURC(urc)
		if evt != nil && evt.Type == "QBLESTAT" && len(evt.Fields) > 0 && evt.Fields[0] == "1" {
			fmt.Println("[URC] 设备已连接！")
			select {
			case connected <- struct{}{}:
			default:
			}
		}
	})
	ble.StartURCListener()

	// 等待连接
	select {
	case <-connected:
		fmt.Println("检测到 BLE 连接，继续后续测试...")
	case <-time.After(60 * time.Second):
		log.Fatalf("60 秒内未检测到 BLE 连接，测试终止")
	}

	ble.StopAdvertising()

	// === 4. GATT 服务端配置 ===
	fmt.Println("[6] 添加服务和特征...")
	if err := ble.AddService("FFF0"); err != nil {
		log.Fatalf("AddService 失败: %v", err)
	}
	if err := ble.AddCharacteristic("FFF1", 0x10); err != nil {
		log.Fatalf("AddCharacteristic 失败: %v", err)
	}
	if err := ble.FinishGATTServer(); err != nil {
		log.Fatalf("FinishGATTServer 失败: %v", err)
	}

	// === 5. 发送 Notify 示例 ===
	fmt.Println("[7] 发送 Notify（如有连接）...")
	_ = ble.SendNotify(0, 1, "48656C6C6F") // "Hello"

	// === 6. 透传模式测试 ===
	fmt.Println("[8] 进入透传模式...")
	if err := ble.EnterTransparent("FFF1"); err != nil {
		log.Fatalf("EnterTransparent 失败: %v", err)
	}
	time.Sleep(1 * time.Second)

	fmt.Println("发送透传数据...")
	_ = ble.SendData("Hello BLE透传\n")

	fmt.Println("退出透传模式...")
	if err := ble.ExitTransparent(); err != nil {
		log.Fatalf("退出透传失败: %v", err)
	}

	fmt.Println("[9] 测试完成")
}
