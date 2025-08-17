package dlt645

import (
	"fmt"
	"io"
)

// Meter 表示电表对象
// 通过串口与电表通信
// port 实现 io.ReadWriteCloser 接口, 例如 github.com/tarm/serial.Port
// addr 为 12 位 BCD 地址
//
// 示例:
//
//	c, _ := serial.OpenPort(&serial.Config{Name: "/dev/ttyUSB0", Baud: 2400})
//	m, _ := dlt645.NewMeter(c, "112233445566")
//	val, _ := m.ReadEnergy(dlt645.RatePeak)
type Meter struct {
	port io.ReadWriteCloser
	addr string
}

// NewMeter 创建 Meter 实例
func NewMeter(port io.ReadWriteCloser, addr string) *Meter {
	return &Meter{port: port, addr: addr}
}

// ReadEnergy 读取指定费率的电量(kWh)
func (m *Meter) ReadEnergy(rate Rate) (float64, error) {
	id, ok := rateDataID[rate]
	if !ok {
		return 0, fmt.Errorf("unsupported rate: %v", rate)
	}
	frame, err := BuildReadFrame(m.addr, id)
	if err != nil {
		return 0, err
	}
	if _, err := m.port.Write(frame); err != nil {
		return 0, err
	}
	buf := make([]byte, 256)
	n, err := m.port.Read(buf)
	if err != nil {
		return 0, err
	}
	return ParseEnergyResponse(buf[:n])
}
