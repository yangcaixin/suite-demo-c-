package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"dlt645server/pkg/dlt645"
	"github.com/tarm/serial"
	"gopkg.in/yaml.v3"
)

// Config 用于读取配置文件
// 包含串口信息和 WEB 服务端口
// 样例见根目录的 config.yaml
//
// serial_port: /dev/ttyUSB0
// baud_rate: 2400
// address: 112233445566
// server_port: 8080

type Config struct {
	SerialPort string `yaml:"serial_port"`
	BaudRate   int    `yaml:"baud_rate"`
	Address    string `yaml:"address"`
	ServerPort int    `yaml:"server_port"`
}

func loadConfig(path string) (*Config, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	cfgPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	cfg, err := loadConfig(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	port, err := serial.OpenPort(&serial.Config{Name: cfg.SerialPort, Baud: cfg.BaudRate})
	if err != nil {
		log.Fatalf("open serial: %v", err)
	}
	meter := dlt645.NewMeter(port, cfg.Address)

	http.HandleFunc("/energy", func(w http.ResponseWriter, r *http.Request) {
		rateStr := r.URL.Query().Get("rate")
		var rate dlt645.Rate
		switch rateStr {
		case "peak":
			rate = dlt645.RatePeak
		case "flat":
			rate = dlt645.RateFlat
		case "valley":
			rate = dlt645.RateValley
		default:
			http.Error(w, "rate must be peak|flat|valley", http.StatusBadRequest)
			return
		}
		val, err := meter.ReadEnergy(rate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%.2f", val)
	})

	addr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
