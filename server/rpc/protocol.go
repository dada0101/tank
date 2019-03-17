package rpc

import (
	"encoding/binary"
	"math"
)

func Float32ToByte(f float32) []byte {
	bits := math.Float32bits(f)
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, bits)
	return b
}

func ByteToFloat32(b []byte) float32 {
	bits := binary.LittleEndian.Uint32(b)
	return math.Float32frombits(bits)
}

func Uint32ToByte(n uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)
	return b
}

func ByteToUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

/*
*\note
* +--+--+--+--+
* |len of buf |
* +--+--+--+--+
* |len of str | -> str : mathod-name
* +--+--+--+--+
* | str-data  |
* |			  |
* +--+--+--+--+
* | otherdata | -> int32 float32 string(about len(4) and string-data) or not
* +--+--+--+--+
 */

 type ProtocolBytes struct {
 	buf []byte
 	pos int32
 }

func NewProtocolBytes (buf []byte) *ProtocolBytes {
	return &ProtocolBytes{
		buf,
		0,
	}
}

 func(br *ProtocolBytes)GetPos() int{
 	return int(br.pos)
 }

 //@Warning this functions have many memory-copy action, only for demo

//@Decode
func(br *ProtocolBytes) DecodeInt32() int32 {
	if br.pos + 4 > int32(len(br.buf)) {
		return 0
	}
	var result = ByteToUint32(br.buf[br.pos:br.pos +4])
	br.pos += 4
	return int32(result)
}

func(br *ProtocolBytes) DecodeFloat32() float32 {
	if br.pos + 4 > int32(len(br.buf)) {
		return 0.0
	}
	var result = ByteToFloat32(br.buf[br.pos:br.pos +4])
	br.pos += 4
	return result
}

func (br *ProtocolBytes) DecodeString() (str string) {
	lens := br.DecodeInt32()
	if br.pos + lens > int32(len(br.buf)) {
		return
	}
	str = string(br.buf[br.pos:br.pos + lens])
	br.pos += lens
	return
}
//@Decode

//@Encode
func(br *ProtocolBytes) EncodeInt32(n int32) {
	bytes := Uint32ToByte(uint32(n))
	br.buf = append(br.buf, bytes...)
}

func (br *ProtocolBytes) EncodeFloat32(f float32) {
	bytes := Float32ToByte(f)
	br.buf = append(br.buf, bytes...)
}

func (br *ProtocolBytes) EncodeString(str string) {
	bytes := []byte(str)
	br.EncodeInt32(int32(len(bytes)))
	br.buf = append(br.buf, bytes...)
}
//@Encode

// \note it will set the first 4 bytes
func (br *ProtocolBytes) SetLength() {
	len := uint32(len(br.buf) - 4)
	b := Uint32ToByte(len)
	br.buf[0], br.buf[1], br.buf[2], br.buf[3] = b[0], b[1], b[2], b[3]
}

func (br *ProtocolBytes)GetBuf() []byte {
	return br.buf
}

//