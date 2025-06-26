package main

import (
	"fmt"
	"log"
	"time"

	"main/blemodule"
	"main/serialport"
)

func main() {
	// 串口设备路径（根据实际修改）
	dev := "/dev/ttyS3"
	baud := 115200

	// 打开串口
	port, err := serialport.OpenSerialPort(dev, baud)
	if err != nil {
		log.Fatalf("串口打开失败: %v", err)
	}
	ble := blemodule.New(port)

	// === 1. 模块基础测试 ===
	fmt.Println("1. 模块重启中...")
	if err := ble.Restart(); err != nil {
		log.Fatalf("Restart 失败: %v", err)
	}
	time.Sleep(2 * time.Second)

	fmt.Println("2. 查询版本号...")
	version, err := ble.GetVersion()
	if err != nil {
		log.Fatalf("获取版本失败: %v", err)
	}
	fmt.Println("版本信息:", version)

	// === 2. 初始化 BLE 模块 ===
	fmt.Println("3. 初始化 BLE peripheral...")
	if err := ble.Init(1); err != nil {
		log.Fatalf("Init 失败: %v", err)
	}

	fmt.Println("4. 设置设备名称...")
	if err := ble.SetDeviceName("HCM111Z"); err != nil {
		log.Fatalf("SetName 失败: %v", err)
	}

	addr, _ := ble.QueryAddress()
	fmt.Println("设备地址:", addr)

	// === 3. 广播测试 ===
	fmt.Println("5. 开始广播...")
	if err := ble.StartAdvertising(); err != nil {
		log.Fatalf("StartAdvertising 失败: %v", err)
	}
	time.Sleep(5 * time.Second)
	ble.StopAdvertising()

	// === 4. GATT 服务端配置 ===
	fmt.Println("6. 添加服务和特征...")
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
	fmt.Println("7. 模拟 Notify 发送（若已连接）...")
	_ = ble.SendNotify(0, 1, "48656C6C6F") // 示例数据 "Hello"

	// === 6. 透传测试 ===
	fmt.Println("8. 进入透传模式...")
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

	fmt.Println("测试完成")
}
