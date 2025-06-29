# HCM111Z BLE 模块控制项目

## 项目简介
本项目用于通过串口控制移远 HCM111Z BLE 模块，支持 BLE 初始化、广播、GATT 服务、透传等核心功能。适合嵌入式、物联网 BLE 外设开发。

## 目录结构
```
├── main.go                // 程序入口，演示 BLE 模块主要功能（自动检测连接）
├── interfaces/            // 只定义串口等接口
│   └── serial.go
├── serialport/            // 串口接口的具体实现
│   └── serialport.go
├── blemodule/             // BLE 模块控制逻辑
│   ├── blemodule.go       // BLE API 封装
│   ├── atcmd.go           // AT 命令底层封装
│   ├── urc.go             // URC 解析与事件
│   └── types.go           // 结构体定义
├── go.mod, go.sum         // Go 依赖管理文件
```

## 依赖说明
- Go 1.18 及以上
- [go.bug.st/serial](https://github.com/bugst/go-serial) 串口库

## 快速开始
1. 安装依赖：
   ```bash
   go mod tidy
   ```
2. 修改 `main.go` 中串口设备路径（如 `/dev/ttyS3`）和波特率。
3. 编译并运行：
   ```bash
   go run main.go
   ```

## 主要功能
- BLE 模块初始化与重启
- 设备名称设置与地址查询
- 广播控制（自动监听连接）
- GATT 服务与特征添加
- Notify 发送
- 透传模式数据收发
- URC 监听与事件驱动

## 使用流程与注意事项
1. 程序启动后会自动初始化 BLE 并进入广播。
2. **请用手机 App（如 nRF Connect）或 PC 端扫描并连接设备。**
3. 程序会自动监听串口 URC，检测到 `+QBLESTAT: 1`（已连接）后，才会继续后续 GATT、Notify、透传等测试。
4. 若 60 秒内未检测到连接，程序会自动终止。
5. 适合自动化测试和开发调试。

## 代码解耦说明
- `interfaces` 只定义接口，便于扩展和测试
- `serialport` 实现接口，便于后续替换为其他实现
- `blemodule` 只依赖接口，主逻辑清晰，易维护

## 许可证
MIT 