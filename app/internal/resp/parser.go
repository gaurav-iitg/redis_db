package resp

import (
	"bufio"
	"io"
	"strconv"
)

func readLine(r *bufio.Reader) ([]byte, error) {
	line, err := r.ReadSlice('\n')
	if err != nil {
		return nil, err
	}
	n := len(line)
	if n < 2 || line[n-2] != '\r' {
		return nil, ErrInvalidTerminator
	}
	return line[:n-2], nil
}

// ReadRESP reads a complete RESP message from the reader, handling TCP streaming
func ReadRESP(r *bufio.Reader) (RESPType, error) {
	// Read the first line to determine the RESP type and metadata
	line, err := readLine(r)
	if err != nil {
		return nil, err
	}

	if len(line) == 0 {
		return nil, ErrInvalidRESPType
	}

	switch line[0] {
	case '+':
		// Simple string: already complete on first line
		return parseSimpleString(line)
	case '-':
		// Simple error: already complete on first line
		return parseSimpleError(line)
	case ':':
		// Integer: already complete on first line
		return parseInteger(line)
	case '$':
		// Bulk string: need to read length, then read data
		return parseBulkStringStreaming(r, line)
	case '*':
		// Array: need to read count, then read each element
		return parseArrayStreaming(r, line)
	default:
		return nil, ErrInvalidRESPType
	}
}

// ParseRESP parses a single line (used when data is already buffered)
func ParseRESP(line []byte) (RESPType, error) {
	if len(line) == 0 {
		return nil, ErrInvalidRESPType
	}

	switch line[0] {
	case '+':
		return parseSimpleString(line)
	case '-':
		return parseSimpleError(line)
	case ':':
		return parseInteger(line)
	case '$':
		// Bulk strings cannot be parsed from a single line - data is on following lines
		// Use ReadRESP() instead for streaming parsing
		return nil, ErrInvalidRESPType
	case '*':
		// Arrays cannot be parsed from a single line - elements are on following lines
		// Use ReadRESP() instead for streaming parsing
		return nil, ErrInvalidRESPType
	default:
		return nil, ErrInvalidRESPType
	}
}

func parseInteger(line []byte) (RESPType, error) {
	// line is in format: :<value> (CRLF already removed by readLine)
	if len(line) < 2 {
		return nil, ErrInvalidRESPType
	}

	// Remove the leading ':'
	valueStr := string(line[1:])

	value, err := strconv.ParseInt(valueStr, 10, 64)
	if err != nil {
		return nil, err
	}

	return Integer{Value: value}, nil
}

// parseBulkStringStreaming handles bulk strings in streaming mode
// line: the "$<length>" part (without data)
func parseBulkStringStreaming(r *bufio.Reader, line []byte) (RESPType, error) {
	// line is in format: $<length>
	if len(line) < 2 {
		return nil, ErrInvalidRESPType
	}

	// Parse the length from the first line
	lengthStr := string(line[1:])
	length, err := strconv.ParseInt(lengthStr, 10, 64)
	if err != nil {
		return nil, err
	}

	// Special case: null bulk string
	if length == -1 {
		return BulkString{Length: length, Value: nil}, nil
	}

	if length < 0 {
		return nil, ErrInvalidRESPType
	}

	// Read exactly 'length' bytes for the data
	data := make([]byte, length)
	n := 0
	for n < int(length) {
		bytesRead, err := r.Read(data[n:])
		if err != nil {
			return nil, err
		}
		if bytesRead == 0 {
			return nil, io.ErrUnexpectedEOF
		}
		n += bytesRead
	}

	// Read the terminating \r\n using bufio methods
	cr, err := r.ReadByte()
	if err != nil {
		return nil, err
	}
	lf, err := r.ReadByte()
	if err != nil {
		return nil, err
	}

	if cr != '\r' || lf != '\n' {
		return nil, ErrInvalidTerminator
	}

	return BulkString{Length: length, Value: data}, nil
}

func parseSimpleError(line []byte) (RESPType, error) {
	// -<message>
	return SimpleError{Message: string(line[1:])}, nil
}

// parseArrayStreaming handles arrays in streaming mode
// line: the "*<count>" part
func parseArrayStreaming(r *bufio.Reader, line []byte) (RESPType, error) {
	// line is in format: *<count>
	if len(line) < 2 {
		return nil, ErrInvalidRESPType
	}

	// Parse the count
	countStr := string(line[1:])
	count, err := strconv.ParseInt(countStr, 10, 64)
	if err != nil {
		return nil, err
	}

	// Special case: null array
	if count == -1 {
		return Array{Count: count, Elements: nil}, nil
	}

	if count < 0 {
		return nil, ErrInvalidRESPType
	}

	// Read 'count' elements
	elements := make([]RESPType, 0, count)
	for i := 0; i < int(count); i++ {
		element, err := ReadRESP(r)
		if err != nil {
			return nil, err
		}
		elements = append(elements, element)
	}

	return Array{Count: count, Elements: elements}, nil
}

func parseSimpleString(line []byte) (RESPType, error) {
	return SimpleString{Value: string(line[1:])}, nil
}
