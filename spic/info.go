package spic

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
)

type TESABInfo struct {
	ESAMNumber                    string `json:"esam_number"`                      //ESAM序列号
	ESAMVersion                   string `json:"esam_version"`                     //ESAM版本号
	SymmetricKeyVersion           string `json:"symmetric_key_version"`            //对称密钥版本
	MainStationCertificateVersion byte   `json:"main_station_certificate_version"` //主站证书版本号
	TerminalCertificateVersion    byte   `json:"terminal_certificate_version"`     //终端证书版本号
	SessionTimeLimit              uint32 `json:"session_time_limit"`               //会话时效门限
	SessionTimeRemainingTime      uint32 `json:"session_time_remaining_time"`      //会话时效剩余时间
	ASCTR                         uint32 `json:"ASCTR"`                            //单地址应用协商计数器
	ARCTR                         uint32 `json:"ARCTR"`                            //主动上报计数器
	AGSEQ                         uint32 `json:"AGSEQ"`                            //应用广播通信序列号
	TerminalCertificateNumber     string `json:"terminal_certificate_number"`      //终端证书序列号
	MainStationCertificateNumber  string `json:"main_station_certificate_number"`  //主站证书序列号
}

func (t *TESABInfo) Decode(buf *bytes.Reader) error {
	//ESAM序列号
	ESAMNumber := make([]byte, 8)
	err := binary.Read(buf, binary.BigEndian, &ESAMNumber)
	if err != nil {
		return err
	}
	t.ESAMNumber = hex.EncodeToString(ESAMNumber)
	//ESAM版本号
	ESAMVersion := make([]byte, 4)
	err = binary.Read(buf, binary.BigEndian, &ESAMVersion)
	if err != nil {
		return err
	}
	t.ESAMVersion = hex.EncodeToString(ESAMVersion)
	//对称密钥版本
	SymmetricKeyVersion := make([]byte, 16)
	err = binary.Read(buf, binary.BigEndian, &SymmetricKeyVersion)
	if err != nil {
		return err
	}
	t.SymmetricKeyVersion = hex.EncodeToString(SymmetricKeyVersion)
	//主站证书版本号
	err = binary.Read(buf, binary.BigEndian, &t.MainStationCertificateVersion)
	if err != nil {
		return err
	}
	//终端证书版本号
	err = binary.Read(buf, binary.BigEndian, &t.TerminalCertificateVersion)
	if err != nil {
		return err
	}
	//会话时效门限
	err = binary.Read(buf, binary.BigEndian, &t.SessionTimeLimit)
	if err != nil {
		return err
	}
	//会话时效剩余时间
	err = binary.Read(buf, binary.BigEndian, &t.SessionTimeRemainingTime)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &t.ASCTR)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &t.ARCTR)
	if err != nil {
		return err
	}
	err = binary.Read(buf, binary.BigEndian, &t.AGSEQ)
	if err != nil {
		return err
	}
	TerminalCertificateNumber := make([]byte, 16)
	err = binary.Read(buf, binary.BigEndian, &TerminalCertificateNumber)
	if err != nil {
		return err
	}
	t.TerminalCertificateNumber = hex.EncodeToString(TerminalCertificateNumber)
	MainStationCertificateNumber := make([]byte, 16)
	err = binary.Read(buf, binary.BigEndian, &MainStationCertificateNumber)
	if err != nil {
		return err
	}
	t.MainStationCertificateNumber = hex.EncodeToString(MainStationCertificateNumber)
	return nil
}
