package main

import (
//	"fmt"
	"math"
	"encoding/binary"
	"bytes"
)

const (
	CHECK_BYTE = 123
	LOG_ALLOC_ID = 1
	LOG_IN = 2
	LOG_ONLINE = 3
	LOG_OFF = 4
)

type InternetPackage struct {
	check		byte
	ptype		byte
	state		byte
	reserve		byte
	id		int32
	X		float32
	Y		float32
	Z		float32
	toX		float32
	toY		float32
	toZ		float32
}

func BytesToFloat32(b []byte) float32 {
	bytesBuffer := bytes.NewBuffer(b)
	var tmp float32
	binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return tmp
}

func BytesToInt32(b []byte) int32 {
    bytesBuffer := bytes.NewBuffer(b)
    var tmp int32
    binary.Read(bytesBuffer, binary.LittleEndian, &tmp)
	return tmp
}

func Float32ToBytes(f float32) []byte {
	bits := math.Float32bits(f)
	bytesBuffer := make([]byte, 4)
	binary.BigEndian.PutUint32(bytesBuffer, bits)
	return bytesBuffer
}

func Int32ToBytes(i int32) []byte {
	bytesBuffer := new(bytes.Buffer)
	binary.Write(bytesBuffer, binary.BigEndian, i)
	return bytesBuffer.Bytes()
}

//代码待优化
func EntityToBytes(p InternetPackage) []byte {
	bytesBuffer := make([]byte, 32)
	bytesBuffer[0] = p.check
	bytesBuffer[1] = p.ptype
	bytesBuffer[2] = p.state
	bytesBuffer[3] = p.reserve
	tempBuffer  := make([]byte, 4)
	tempBuffer   = Int32ToBytes(p.id)
	bytesBuffer[4] = tempBuffer[0]
	bytesBuffer[5] = tempBuffer[1]
	bytesBuffer[6] = tempBuffer[2]
	bytesBuffer[7] = tempBuffer[3]
	tempBuffer   = Float32ToBytes(p.X)
	bytesBuffer[8] = tempBuffer[0]
	bytesBuffer[9] = tempBuffer[1]
	bytesBuffer[10] = tempBuffer[2]
	bytesBuffer[11] = tempBuffer[3]
	tempBuffer   = Float32ToBytes(p.Y)
	bytesBuffer[12] = tempBuffer[0]
	bytesBuffer[13] = tempBuffer[1]
	bytesBuffer[14] = tempBuffer[2]
	bytesBuffer[15] = tempBuffer[3]
	tempBuffer   = Float32ToBytes(p.Z)
	bytesBuffer[16] = tempBuffer[0]
	bytesBuffer[17] = tempBuffer[1]
	bytesBuffer[18] = tempBuffer[2]
	bytesBuffer[19] = tempBuffer[3]
	tempBuffer   = Float32ToBytes(p.toX)
	bytesBuffer[20] = tempBuffer[0]
	bytesBuffer[21] = tempBuffer[1]
	bytesBuffer[22] = tempBuffer[2]
	bytesBuffer[23] = tempBuffer[3]
	tempBuffer   = Float32ToBytes(p.toY)
	bytesBuffer[24] = tempBuffer[0]
	bytesBuffer[25] = tempBuffer[1]
	bytesBuffer[26] = tempBuffer[2]
	bytesBuffer[27] = tempBuffer[3]
	tempBuffer   = Float32ToBytes(p.toZ)
	bytesBuffer[28] = tempBuffer[0]
	bytesBuffer[29] = tempBuffer[1]
	bytesBuffer[30] = tempBuffer[2]
	bytesBuffer[31] = tempBuffer[3]
	return bytesBuffer
}
