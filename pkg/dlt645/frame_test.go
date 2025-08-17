package dlt645

import (
	"bytes"
	"math"
	"testing"
)

func TestBuildReadFrame(t *testing.T) {
	frame, err := BuildReadFrame("112233445566", rateDataID[RatePeak])
	if err != nil {
		t.Fatalf("BuildReadFrame error: %v", err)
	}
	expected := []byte{0x68, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x68, 0x11, 0x04, 0x33, 0x34, 0x33, 0x33, 0x17, 0x16}
	if !bytes.Equal(frame, expected) {
		t.Fatalf("unexpected frame: %x", frame)
	}
}

func TestParseEnergyResponse(t *testing.T) {
	frame := []byte{0x68, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x68, 0x91, 0x08, 0x33, 0x34, 0x33, 0x33, 0x89, 0x67, 0x45, 0x33, 0x03, 0x16}
	val, err := ParseEnergyResponse(frame)
	if err != nil {
		t.Fatalf("parse error: %v", err)
	}
	if math.Abs(val-1234.56) > 0.01 {
		t.Fatalf("expect 1234.56 got %f", val)
	}
}

// stubPort 用于模拟串口设备
type stubPort struct {
	written []byte
	reply   []byte
}

func (s *stubPort) Write(p []byte) (int, error) {
	s.written = append(s.written, p...)
	return len(p), nil
}

func (s *stubPort) Read(p []byte) (int, error) {
	copy(p, s.reply)
	return len(s.reply), nil
}

func (s *stubPort) Close() error { return nil }

func TestMeterReadEnergy(t *testing.T) {
	resp := []byte{0x68, 0x66, 0x55, 0x44, 0x33, 0x22, 0x11, 0x68, 0x91, 0x08, 0x33, 0x34, 0x33, 0x33, 0x89, 0x67, 0x45, 0x33, 0x03, 0x16}
	sp := &stubPort{reply: resp}
	m := NewMeter(sp, "112233445566")
	val, err := m.ReadEnergy(RatePeak)
	if err != nil {
		t.Fatalf("read energy: %v", err)
	}
	if math.Abs(val-1234.56) > 0.01 {
		t.Fatalf("expect 1234.56 got %f", val)
	}
	// 校验请求报文
	expectReq, _ := BuildReadFrame("112233445566", rateDataID[RatePeak])
	if !bytes.Equal(expectReq, sp.written) {
		t.Fatalf("unexpected request: %x", sp.written)
	}
}
