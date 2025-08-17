package dlt645

import (
	"encoding/hex"
	"errors"
)

// Rate 表示费率类型
// 1-尖峰、2-平段、3-谷段
// 此处仅实现常见的三个费率
// 如需扩展可继续添加枚举
// 枚举值对应 DL/T 645 标准中的数据标识
// 0x00000100 费率1(峰) ,0x00000200 费率2(平) ,0x00000300 费率3(谷)
type Rate int

const (
	RatePeak   Rate = 1 // 峰
	RateFlat   Rate = 2 // 平
	RateValley Rate = 3 // 谷
)

var rateDataID = map[Rate]uint32{
	RatePeak:   0x00000100,
	RateFlat:   0x00000200,
	RateValley: 0x00000300,
}

// BuildReadFrame 生成读取数据的报文
// addr 为 12 位表地址(以 BCD 格式给出, 如 "112233445566")
// dataID 为数据标识
func BuildReadFrame(addr string, dataID uint32) ([]byte, error) {
	addrBytes, err := encodeAddress(addr)
	if err != nil {
		return nil, err
	}
	// 构建帧
	frame := []byte{0x68}
	frame = append(frame, addrBytes[:]...) // 地址
	frame = append(frame, 0x68)
	frame = append(frame, 0x11) // 控制码: 读数据
	frame = append(frame, 0x04) // 数据长度: 4 字节数据标识
	did := encodeDataID(dataID)
	frame = append(frame, did...)
	cs := checksum(frame)
	frame = append(frame, cs, 0x16)
	return frame, nil
}

// ParseEnergyResponse 解析电表返回的能量数据
// frame 为完整的16进制帧
// 返回值为 kWh
func ParseEnergyResponse(frame []byte) (float64, error) {
	if len(frame) < 16 { // 最小长度
		return 0, errors.New("frame too short")
	}
	if frame[0] != 0x68 || frame[len(frame)-1] != 0x16 {
		return 0, errors.New("frame format error")
	}
	if checksum(frame[:len(frame)-2]) != frame[len(frame)-2] {
		return 0, errors.New("checksum error")
	}
	// 跳过固定头部: 68 + 地址6 + 68 + 控制码 + 长度
	// 索引从第10字节开始即数据区
	data := frame[10 : len(frame)-2]
	if len(data) < 8 {
		return 0, errors.New("data field too short")
	}
	// 数据字段前4字节为数据标识,后4字节为数值
	decoded := make([]byte, len(data))
	for i, b := range data {
		decoded[i] = b - 0x33
	}
	val := bcdToFloat(decoded[4:8])
	return val, nil
}

// encodeAddress 将 12 位 BCD 地址编码并倒序
func encodeAddress(addr string) ([6]byte, error) {
	var res [6]byte
	if len(addr) != 12 {
		return res, errors.New("address must be 12 digits")
	}
	for i := 0; i < 6; i++ {
		pair := addr[2*(5-i) : 2*(5-i)+2]
		b, err := hex.DecodeString(pair)
		if err != nil {
			return res, err
		}
		res[i] = b[0]
	}
	return res, nil
}

// encodeDataID 编码数据标识并 +0x33
func encodeDataID(id uint32) []byte {
	bs := []byte{byte(id), byte(id >> 8), byte(id >> 16), byte(id >> 24)}
	for i := range bs {
		bs[i] += 0x33
	}
	return bs
}

// checksum 计算校验和
func checksum(data []byte) byte {
	var cs byte
	for _, b := range data {
		cs += b
	}
	return cs
}

// bcdToFloat 将4字节 BCD(低位在前) 转换为浮点数,假定小数位2
func bcdToFloat(bs []byte) float64 {
	var digits [8]byte
	for i := 0; i < 4; i++ {
		b := bs[i]
		digits[2*i] = b & 0x0F
		digits[2*i+1] = b >> 4
	}
	// 逆序排列
	n := 0
	for i := 7; i >= 0; i-- {
		n = n*10 + int(digits[i])
	}
	return float64(n) / 100.0
}
