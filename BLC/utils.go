package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

//int64 => []byte
func IntToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	//BigEndian指定大端小端
	//binary.Write是将数据的二进制格式写入字节缓冲区中
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}
