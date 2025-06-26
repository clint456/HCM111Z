package interfaces

import "time"

// SerialPort 是串口读写接口抽象
// 只定义接口，不包含实现
type SerialPort interface {
	Write(data string) error
	ReadLine(timeout time.Duration) (string, error)
}
