package utils

import (
	"bytes"
	"encoding/binary"
	"log"
)

//int64 => []byte
func Int64ToBytes(num int64) []byte {
	buff := new(bytes.Buffer)
	//BigEndian指定大端小端
	//binary.Write是将数据的二进制格式写入字节缓冲区中
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buff.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)
	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}


// ReverseBytes reverses a byte array
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}
