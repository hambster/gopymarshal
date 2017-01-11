package gopymarshal

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

var (
	ErrType = errors.New("unsupport type")
)

func EmptyMap() (ret []byte) {
	return []byte{CODE_DICT, CODE_STOP}
}

func Marshal(data interface{}) (ret []byte, retErr error) {
	if nil == data {
		return
	}

	var buffer bytes.Buffer
	if retErr = marshal(&buffer, data); nil != retErr {
		return
	}

	ret = buffer.Bytes()
	return
}

func marshal(buffer *bytes.Buffer, data interface{}) (ret error) {
	if nil == data {
		ret = buffer.WriteByte(CODE_NONE)
		return
	}

	switch d := data.(type) {
	case int32:
		ret = writeInt32(buffer, d)
	case string:
		ret = writeString(buffer, d)
	case []byte:
		ret = writeBytes(buffer, d)
	case float64:
		ret = writeFloat64(buffer, d)
	case []interface{}:
		ret = writeList(buffer, d)
	case map[interface{}]interface{}:
		ret = writeDict(buffer, d)
	case map[string]interface{}:
		ret = writeDictStrInter(buffer, d)
	default:
		ret = ErrType
	}

	return
}

func writeInt32(buffer *bytes.Buffer, data int32) (ret error) {
	if ret = buffer.WriteByte(CODE_INT); nil != ret {
		return
	}

	ret = binary.Write(buffer, binary.LittleEndian, data)
	return
}

func writeBytes(buffer *bytes.Buffer, data []byte) (ret error) {
	if ret = buffer.WriteByte(CODE_TSTRING); nil != ret {
		return
	}

	if ret = binary.Write(buffer, binary.LittleEndian, int32(len(data))); nil == ret {
		_, ret = buffer.Write(data)
	}

	return
}

func writeString(buffer *bytes.Buffer, data string) (ret error) {
	if ret = buffer.WriteByte(CODE_UNICODE); nil != ret {
		return
	}

	if ret = binary.Write(buffer, binary.LittleEndian, int32(len(data))); nil == ret {
		_, ret = buffer.WriteString(data)
	}

	return
}

func writeFloat64(buffer *bytes.Buffer, data float64) (ret error) {
	if ret = buffer.WriteByte(CODE_FLOAT); nil != ret {
		return
	}

	bits := math.Float64bits(data)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, bits)
	_, ret = buffer.Write(buf)
	return
}

func writeList(buffer *bytes.Buffer, data []interface{}) (ret error) {
	listLen := len(data)

	if ret = buffer.WriteByte(CODE_LIST); nil != ret {
		return
	}
	if ret = writeInt32(buffer, int32(listLen)); nil != ret {
		return
	}

	for idx := 0; idx < listLen; idx++ {
		if ret = marshal(buffer, data[idx]); nil != ret {
			break
		}
	}

	return
}

func writeNone(buffer *bytes.Buffer) (ret error) {
	ret = buffer.WriteByte(CODE_NONE)
	return
}

func writeDictStrInter(buffer *bytes.Buffer, data map[string]interface{}) (ret error) {
	if ret = buffer.WriteByte(CODE_DICT); nil != ret {
		return
	}

	for k, v := range data {
		if ret = marshal(buffer, k); nil != ret {
			return
		}

		if ret = marshal(buffer, v); nil != ret {
			return
		}
	}

	ret = buffer.WriteByte(CODE_STOP)
	return
}

func writeDict(buffer *bytes.Buffer, data map[interface{}]interface{}) (ret error) {
	if ret = buffer.WriteByte(CODE_DICT); nil != ret {
		return
	}

	for k, v := range data {
		if ret = marshal(buffer, k); nil != ret {
			return
		}

		if ret = marshal(buffer, v); nil != ret {
			return
		}
	}

	ret = buffer.WriteByte(CODE_STOP)
	return
}
