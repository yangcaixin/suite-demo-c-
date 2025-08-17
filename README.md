# suite-demo-c-

1.加解密文件为API目录下的DingTalkCrypt.cs文件

2.加解密demo为receiv.ashx.cs

3.vs版本为2015

4.本demo为钱姓开发者提供，在此对他表示感谢
## Go DL/T 645 示例

本仓库额外提供基于 Go 语言的 DL/T 645 协议采集示例，位于 `pkg/dlt645` 和 `cmd/dlt645server` 目录。

### 运行示例

1. 修改 `config.yaml` 中的串口、波特率、电表地址及 Web 服务端口。
2. 编译并运行服务端：

```bash
# 安装依赖并运行单元测试
cd /workspace/suite-demo-c-
go test ./...

# 启动 Web 服务
cd /workspace/suite-demo-c-/cmd/dlt645server
go run . -config ../config.yaml
```

启动后可通过 `http://<RaspberryPiIP>:8080/energy?rate=peak` 访问峰段电量，`rate` 参数支持 `peak`、`flat`、`valley`。

