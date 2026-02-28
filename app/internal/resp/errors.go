package resp

import "errors"

var ErrUnknownType = errors.New("unknown RESP value type")
var ErrInvalidTerminator = errors.New("line does not end with \\r\\n")
var ErrInvalidRESPType = errors.New("invalid RESP type")
