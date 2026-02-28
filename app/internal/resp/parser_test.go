package resp

import (
	"bufio"
	"bytes"
	"reflect"
	"testing"
)

type testCase struct {
	name    string
	line    []byte
	want    RESPType
	wantErr bool
}

func TestParseRESP(t *testing.T) {

	invalidPrefixTests := []testCase{
		{
			name:    "invalid prefix",
			line:    []byte("&ERR"),
			want:    nil,
			wantErr: true,
		},
	}

	simpleStringTests := []testCase{
		{
			name:    "simple string",
			line:    []byte("+OK"),
			want:    SimpleString{Value: "OK"},
			wantErr: false,
		},
		{
			name:    "empty simple string",
			line:    []byte("+"),
			want:    SimpleString{Value: ""},
			wantErr: false,
		},
	}

	simpleErrorTests := []testCase{
		{
			name:    "simple error",
			line:    []byte("-ERR unknown command"),
			want:    SimpleError{Message: "ERR unknown command"},
			wantErr: false,
		},
	}

	integerTests := []testCase{
		{
			name:    "integer",
			line:    []byte(":1000"),
			want:    Integer{Value: 1000},
			wantErr: false,
		},
		{
			name:    "negative integer",
			line:    []byte(":-1000"),
			want:    Integer{Value: -1000},
			wantErr: false,
		},
		{
			name:    "invalid integer",
			line:    []byte(":abc"),
			want:    nil,
			wantErr: true,
		},
	}
	var tests []testCase
	tests = append(tests, invalidPrefixTests...)
	tests = append(tests, simpleStringTests...)
	tests = append(tests, simpleErrorTests...)
	tests = append(tests, integerTests...)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRESP(tt.line)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseRESP(%q) error = %v, wantErr %v", tt.line, err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRESP(%q) = %#v, want %#v", tt.line, got, tt.want)
			}
		})
	}
}

func TestReadRESP(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    RESPType
		wantErr bool
	}{
		{
			name:    "simple string",
			input:   "+OK\r\n",
			want:    SimpleString{Value: "OK"},
			wantErr: false,
		},
		{
			name:    "missing CRLF",
			input:   "+NOCRLF",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "bulk string",
			input:   "$5\r\nhello\r\n",
			want:    BulkString{Value: []byte("hello")},
			wantErr: false,
		},
		{
			name:    "null bulk string",
			input:   "$-1\r\n",
			want:    BulkString{Value: nil},
			wantErr: false,
		},
		{
			name:    "empty bulk string",
			input:   "$0\r\n\r\n",
			want:    BulkString{Value: []byte("")},
			wantErr: false,
		},
		{
			name:    "integer",
			input:   ":1000\r\n",
			want:    Integer{Value: 1000},
			wantErr: false,
		},
		{
			name:    "array of bulk strings",
			input:   "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			want:    Array{Elements: []RESPType{BulkString{Value: []byte("hello")}, BulkString{Value: []byte("world")}}},
			wantErr: false,
		},
		{
			name:    "simple error",
			input:   "-ERR unknown command\r\n",
			want:    SimpleError{Message: "ERR unknown command"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufio.NewReader(bytes.NewBufferString(tt.input))
			got, err := ReadRESP(r)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ReadRESP(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadRESP(%q) = %#v, want %#v", tt.input, got, tt.want)
			}
		})
	}
}
