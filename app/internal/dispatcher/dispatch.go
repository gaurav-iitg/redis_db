package dispatcher

import "github.com/codecrafters-io/redis-starter-go/app/internal/resp"

type CommandFunc func(args []string) resp.RESPType

var commandTable = map[string]CommandFunc{
	COMMAND_PING: handlePing,
	COMMAND_SET:  handleSet,
	COMMAND_GET:  handleGet,
	COMMAND_ECHO: handleEcho,
}

func handlePing(args []string) resp.RESPType {
	return resp.SimpleString{Value: "PONG"}
}

func handleSet(args []string) resp.RESPType {
	return resp.SimpleString{Value: "OK"}
}

func handleGet(args []string) resp.RESPType {
	return resp.BulkString{Length: 0, Value: nil}
}

func handleEcho(args []string) resp.RESPType {
	if len(args) == 0 {
		return resp.BulkString{Length: -1, Value: []byte{}}
	}
	return resp.BulkString{Length: int64(len(args[0])), Value: []byte(args[0])}
}

func Execute(cmd resp.RESPType) (resp.RESPType, error) {

	cmd_array, ok := cmd.(resp.Array)
	if !ok {
		return resp.SimpleError{Message: "ERR invalid command format"}, nil
	}

	cmd_string, ok := cmd_array.Elements[0].(resp.BulkString)
	if !ok {
		return resp.SimpleError{Message: "ERR invalid command format"}, nil
	}

	commandFunc, ok := commandTable[string(cmd_string.Value)]
	if !ok {
		return resp.SimpleError{Message: "ERR unknown command"}, nil
	}
	var args []string
	for i := 1; i < len(cmd_array.Elements); i++ {
		args = append(args, string(cmd_array.Elements[i].(resp.BulkString).Value))
	}
	return commandFunc(args), nil
}
