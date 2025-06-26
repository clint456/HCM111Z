package serialport

import (
	"strings"
	"time"

	"main/interfaces"

	"go.bug.st/serial"
)

// serialPortImpl 是 interfaces.SerialPort 的实现
// 不导出，仅包内使用
type serialPortImpl struct {
	port serial.Port
}

// OpenSerialPort 打开串口设备，返回 interfaces.SerialPort
func OpenSerialPort(device string, baud int) (interfaces.SerialPort, error) {
	mode := &serial.Mode{
		BaudRate: baud,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
		DataBits: 8,
	}
	p, err := serial.Open(device, mode)
	if err != nil {
		return nil, err
	}
	return &serialPortImpl{port: p}, nil
}

func (s *serialPortImpl) Write(data string) error {
	_, err := s.port.Write([]byte(data))
	return err
}

func (s *serialPortImpl) ReadLine(timeout time.Duration) (string, error) {
	s.port.SetReadTimeout(timeout)
	buf := make([]byte, 256)
	n, err := s.port.Read(buf)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(buf[:n])), nil
}
