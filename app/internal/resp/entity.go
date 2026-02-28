package resp

type RESPType interface{}

type SimpleString struct {
	Value string
}

type SimpleError struct {
	Message string
}

type Integer struct {
	Value int64
}

type BulkString struct {
	Length int64
	Value  []byte
}

type Array struct {
	Count    int64
	Elements []RESPType
}
