package BLC

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/json"
	"encoding/gob"
	"fmt"
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

	return buff.Bytes() // [0 0 0 0 0 0 1 0]
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

//对象序列化
func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

//对象序列化(包含接口)
func gobEncodeWithRegister(data interface{}, inter interface{}) []byte {
	var buff bytes.Buffer
	gob.Register(inter)
	encoder := gob.NewEncoder(&buff)
	err := encoder.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

// 反序列化：将字节数组反序列化为block对象
func gobDecode(blockBytes []byte, o interface{}) {
	decoder := gob.NewDecoder(bytes.NewReader(blockBytes))
	err := decoder.Decode(o)
	if err != nil {
		log.Panic(err)
	}
}

//将命令转换成字节数组
func commandToBytes(command string) []byte {
	var bytes [COMMAND_LENGTH]byte //设置为12字节长度

	//将命令string转换为byte格式并放进bytes, 剩余位置用0填充
	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

//将字节数组转换成命令
func bytesToCommand(commandBtyes []byte) string {
	var command []byte

	//去掉commandBytes中的0
	for _, b := range commandBtyes {
		if b != 0x0 {
			command = append(command, b)
		}
	}
	//fmt.Println("commandBytes:", commandBtyes, "command: ", command)
	return fmt.Sprintf("%s", command)
}
