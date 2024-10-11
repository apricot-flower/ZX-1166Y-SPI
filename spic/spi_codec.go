package spic

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"github.com/ecc1/spi"
)

type SPICodec struct {
	device *spi.Device
	Dev    string
	Mode   uint8
	Speed  int
}

func (s *SPICodec) Open() error {
	device, err := spi.Open(s.Dev, s.Speed, 0)
	if err != nil {
		return err
	}
	s.device = device
	err = s.device.SetMaxSpeed(s.Speed)
	if err != nil {
		return err
	}
	err = s.device.SetMode(s.Mode)
	return err
}

func (s *SPICodec) Close() error {
	if s.device == nil {
		return nil
	}
	return s.device.Close()
}

// SessionKeyConnect 建立应用连接（会话密钥协商）
// ucSessionData：服务器随机数，48 字节
// ucSign:服务器签名信息，Len-48 字节
func (s *SPICodec) SessionKeyConnect(ucOutSessionInit []byte, ucOutSign []byte) (ucSessionData string, ucSign string, err error) {
	if len(ucOutSessionInit) == 32 {
		return "", "", errors.New("ucOutSessionInit长度应为32！")
	}
	buf := new(bytes.Buffer)
	tx := `80020000`
	txArr, err := hex.DecodeString(tx)
	if err != nil {
		return "", "", err
	}
	_, err = buf.Write(txArr)
	if err != nil {
		return "", "", err
	}
	length := uint16(len(ucOutSessionInit) + len(ucOutSign))
	err = binary.Write(buf, binary.BigEndian, length)
	if err != nil {
		return "", "", err
	}
	err = binary.Write(buf, binary.BigEndian, ucOutSessionInit)
	if err != nil {
		return "", "", err
	}
	err = binary.Write(buf, binary.BigEndian, ucOutSign)
	if err != nil {
		return "", "", err
	}
	//计算cs
	txArr = buf.Bytes()
	buf.Reset()
	cs := Cs(txArr)
	txArr = append([]byte{0x55}, txArr...)
	txArr = append(txArr, cs)
	data, err := s.TransferBytes(txArr, 1024)
	if err != nil {
		return "", "", err
	}
	ucSessionData = hex.EncodeToString(data[:48])
	ucSign = hex.EncodeToString(data[48:])
	return ucSessionData, ucSign, nil
}

// SelectMasterStationCertificate 获取主站证书
func (s *SPICodec) SelectMasterStationCertificate() (string, error) {
	tx := `558036000C000045`
	data, err := s.TransferString(tx, 4096)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// SelectTerminalCertificate 获取终端证书
func (s *SPICodec) SelectTerminalCertificate() (string, error) {
	tx := `558036000B000042`
	data, err := s.TransferString(tx, 4096)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(data), nil
}

// SelectESAMInfos 主站下发获取ESAM信息命令
func (s *SPICodec) SelectESAMInfos() (*TESABInfo, error) {
	tx := `55803600FF0000B6`
	data, err := s.TransferString(tx, 1024)
	if err != nil {
		return nil, err
	}
	info := &TESABInfo{}
	err = info.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (s *SPICodec) TransferString(tx string, length int) ([]byte, error) {
	txArr, err := hex.DecodeString(tx)
	if err != nil {
		return nil, err
	}
	return s.TransferBytes(txArr, length)
}

func (s *SPICodec) TransferBytes(tx []byte, length int) ([]byte, error) {
	txArr := make([]byte, length)
	copy(txArr, tx)
	rxArr := make([]byte, length)
	err := s.device.Transfer(txArr, rxArr)
	if err != nil {
		return nil, err
	}
	rxArr, err = decode(rxArr)
	if err != nil {
		return nil, err
	}
	return rxArr, nil
}

// 拆包器
func decode(array []byte) ([]byte, error) {
	buf := bytes.NewReader(array)
	var err error
	var start byte
	for start != 0x55 {
		err = binary.Read(buf, binary.BigEndian, &start)
		if err != nil {
			return nil, err
		}
	}
	//读取状态码
	status := make([]byte, 2)
	err = binary.Read(buf, binary.BigEndian, &status)
	if err != nil {
		return nil, err
	}
	//验证错误码
	if hex.EncodeToString(status) != "9000" {
		return nil, errors.New("不是正确报文，返回了错误码，错误码为：" + hex.EncodeToString(status))
	}
	//解析长度
	length := make([]byte, 2)
	err = binary.Read(buf, binary.BigEndian, &length)
	if err != nil {
		return nil, err
	}
	data := make([]byte, uint16(length[0])<<8|uint16(length[1]))
	err = binary.Read(buf, binary.BigEndian, &data)
	if err != nil {
		return nil, err
	}
	var csValue byte
	err = binary.Read(buf, binary.BigEndian, &csValue)
	if err != nil {
		return nil, err
	}
	//todo 计算cs
	csData := append(status, length...)
	csData = append(csData, data...)
	if Cs(csData) != csValue {
		return nil, errors.New("cs错误！")
	}
	return data, nil
}

// Cs 计算校验值
func Cs(data []byte) byte {
	var xorSum byte = 0
	for _, b := range data {
		xorSum ^= b
	}

	// 取反最终的异或值
	invertedXorSum := ^xorSum
	return invertedXorSum
}
