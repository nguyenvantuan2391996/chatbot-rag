package utils

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"
)

func DecodeUnsafeF32(bs []byte) []float32 {
	return unsafe.Slice((*float32)(unsafe.Pointer(&bs[0])), len(bs)/4)
}

func Float32SliceToBase64(floats []float32) string {
	bytes := make([]byte, 4*len(floats))
	for i, f := range floats {
		binary.LittleEndian.PutUint32(bytes[i*4:], math.Float32bits(f))
	}
	return base64.StdEncoding.EncodeToString(bytes)
}

func Base64ToFloat32Slice(s string) ([]float32, error) {
	bytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	if len(bytes)%4 != 0 {
		return nil, fmt.Errorf("invalid byte length %d", len(bytes))
	}

	floats := make([]float32, len(bytes)/4)
	for i := 0; i < len(floats); i++ {
		bits := binary.LittleEndian.Uint32(bytes[i*4:])
		floats[i] = math.Float32frombits(bits)
	}
	return floats, nil
}
