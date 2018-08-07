package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
)

// 将int64转换为[]uint8 字节数组(十进制转为256进制？）
//	uint8       the set of all unsigned  8-bit integers (0 to 255)
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	/*
	big endian：最高字节在地址最低位，最低字节在地址最高位，依次排列。
	little endian：最低字节在最低位，最高字节在最高位，反序排列。
	 */
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()// [0 0 0 0 0 0 1 0]
}
