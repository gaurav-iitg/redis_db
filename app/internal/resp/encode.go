package resp

import (
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
	return []byte("$" + strconv.FormatInt(v.Length, 10) + "\r\n" + string(v.Value) + "\r\n"), nil
}

func encodeArray(v Array) ([]byte, error) {
	// Array format: *<count>\r\n<element1><element2>...<elementN>
	resp := "*" + strconv.FormatInt(v.Count, 10) + "\r\n"
	for _, element := range v.Elements {
		encodedElement, err := EncodeResp(element)
		if err != nil {
			return nil, err
		}
		resp += string(encodedElement)
	}
	return []byte(resp), nil
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
