package main

import (
	"codec-SPI/spic"
	"fmt"
)

func main() {
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
	teasmInfos, err := spiCodec.SelectMasterStationCertificate()
	if err != nil {
		fmt.Printf("采集主站下发获取 TESAM 信息命令错误：%s \n", err.Error())
		return
	}
	fmt.Printf("TESAM 信息：%s \n", teasmInfos)
}
