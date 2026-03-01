package resp

import (
	"bytes"
	"strconv"
)

func EncodeResp(v RESPType) ([]byte, error) {
	switch t := v.(type) {
	case SimpleString:
		return encodeSimpleString(t)
	case SimpleError:
		return encodeSimpleError(t)
	case Integer:
		return encodeInteger(t)
	case BulkString:
		return encodeBulkString(t)
	case Array:
		return encodeArray(t)
	default:
		return nil, ErrUnknownType
	}
}

func encodeBulkString(v BulkString) ([]byte, error) {
	// Bulk string format: $<length>\r\n<data>\r\n
	var buf bytes.Buffer
	if v.Value == nil {
		buf.WriteString("$-1\r\n")
	} else {
		buf.WriteString("$")
		buf.WriteString(strconv.FormatInt(v.Length, 10))
		buf.WriteString("\r\n")
		buf.Write(v.Value)
		buf.WriteString("\r\n")
	}
	return buf.Bytes(), nil
}

func encodeArray(v Array) ([]byte, error) {
	// Array format: *<count>\r\n<element1><element2>...<elementN>
	var buf bytes.Buffer
	buf.WriteString("*")
	buf.WriteString(strconv.FormatInt(v.Count, 10))
	buf.WriteString("\r\n")
	for _, element := range v.Elements {
		encodedElement, err := EncodeResp(element)
		if err != nil {
			return nil, err
		}
		buf.Write(encodedElement)
	}
	return buf.Bytes(), nil
}

func encodeSimpleString(v SimpleString) ([]byte, error) {
	// Simple string format: +<string>\r\n
	return []byte("+" + v.Value + "\r\n"), nil
}

func encodeSimpleError(v SimpleError) ([]byte, error) {
	// Simple error format: -<message>\r\n
	return []byte("-" + v.Message + "\r\n"), nil
}

func encodeInteger(v Integer) ([]byte, error) {
	// Integer format: :<number>\r\n
	return []byte(":" + strconv.FormatInt(v.Value, 10) + "\r\n"), nil
}
