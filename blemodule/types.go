package blemodule

// types.go 定义 BLE 相关结构体，便于管理和扩展

type BLEConnectionInfo struct {
	ConnIdx   int    // 连接索引
	MAC       string // 设备 MAC
	Name      string // 设备名称
	Connected bool   // 是否已连接
}

type BLEService struct {
	UUID string
}

type BLECharacteristic struct {
	UUID       string
	Properties int
	Handle     int
}

// 更多结构体可按需扩展
