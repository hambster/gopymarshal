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

	switch data.(type) {
	case int32:
		ret = writeInt32(buffer, data.(int32))
	case string:
		ret = writeString(buffer, data.(string))
	case float64:
		ret = writeFloat64(buffer, data.(float64))
	case []interface{}:
		ret = writeList(buffer, data.([]interface{}))
	case map[interface{}]interface{}:
		ret = writeDict(buffer, data.(map[interface{}]interface{}))
	case map[string]interface{}:
		tmp := make(map[interface{}]interface{})
		for k, v := range data.(map[string]interface{}) {
			tmp[k] = v
		}
		ret = writeDict(buffer, tmp)
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
	listLen := int32(0)
	for idx := 0; idx < len(data); idx++ {
		if !isValidData(data[idx]) {
			continue
		}

		listLen++
	}

	if ret = buffer.WriteByte(CODE_LIST); nil != ret {
		return
	}
	if ret = writeInt32(buffer, listLen); nil != ret {
		return
	}

	for idx := 0; idx < len(data); idx++ {
		if !isValidData(data[idx]) {
			continue
		}

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

func writeDict(buffer *bytes.Buffer, data map[interface{}]interface{}) (ret error) {
	if ret = buffer.WriteByte(CODE_DICT); nil != ret {
		return
	}

	for k, v := range data {
		if isValidData(k) && isValidData(v) {
			if ret = marshal(buffer, k); nil != ret {
				return
			}

			if ret = marshal(buffer, v); nil != ret {
				return
			}
		}
	}

	ret = buffer.WriteByte(CODE_STOP)
	return
}

func isValidData(data interface{}) (ret bool) {
	if nil == data {
		ret = true
		return
	}

	switch data.(type) {
	case int32:
		ret = true
	case string:
		ret = true
	case float64:
		ret = true
	case []interface{}:
		ret = true
	case map[interface{}]interface{}:
		ret = true
	case map[string]interface{}:
		ret = true
	}

	return
}
