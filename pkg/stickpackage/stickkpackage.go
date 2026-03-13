package stickpackage

import "encoding/binary"

func Pack(data []byte) []byte {
	length := len(data)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[:4], uint32(length))
	copy(buf[4:], data)
	return buf
}
func UnPack(data []byte) ([]byte, error) {
	if len(data) < 4 {
		return nil, nil
	}
	length := binary.BigEndian.Uint32(data[:4])
	if len(data) < int(4+length) {
		return nil, nil
	}
	return data[4 : 4+length], nil
}
func UnpackLength(data []byte) uint32 {
	return binary.BigEndian.Uint32(data)
}
