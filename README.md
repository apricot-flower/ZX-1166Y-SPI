# ZX-1166Y-SPI
## 操作1166Y
## 使用教程
> 方法都在spi_codec.go中

> 需要自己组装流程

## 初始化
```go
spiCodec := spic.SPICodec{
		Dev:   "/dev/spidev0.0",
		Mode:  3,
		Speed: 25000,
	}
err := spiCodec.Open()
if err != nil {
    fmt.Printf("打开SPI通信错误：%s \n", err.Error())
    return
}
defer func() {
    err = spiCodec.Close()
    if err != nil {
        fmt.Printf("关闭SPI通信错误：%s \n", err.Error())
        return
    }
}()
```
## 举例 （获取主站证书）
```go
certificate, err := spiCodec.SelectMasterStationCertificate()
```
