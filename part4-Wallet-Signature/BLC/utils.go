package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
)

// 将int64转换为[]uint8 字节数组(十进制转为256进制？）
//	uint8       the set of all unsigned  8-bit integers (0 to 255)
func IntToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	/*
	big endian：最高字节在地址最低位，最低字节在地址最高位，依次排列。
	little endian：最低字节在最低位，最高字节在最高位，反序排列。
	 */
	//将二进制数据写入w
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()// [0 0 0 0 0 0 1 0]
}


// 标准的JSON字符串转数组
func JSONToArray(jsonString string) []string {

	//json 到 []string
	var sArr []string
	if err := json.Unmarshal([]byte(jsonString), &sArr); err != nil {
		log.Panic(err)
	}
	return sArr
}

//字节数组反转
func ReverseBytes(data []byte) {
	for i, j := 0, len(data)-1; i < j; i, j = i+1, j-1 {
		data[i], data[j] = data[j], data[i]
	}
}