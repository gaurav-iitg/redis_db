package dispatcher

import (
	"strings"

	"github.com/redis-go/app/internal/datastore"
	"github.com/redis-go/app/internal/resp"
)

func (d *Dispatcher) handlePing(args [][]byte) resp.RESPType {
	if len(args) == 0 {
		return resp.SimpleString{Value: "PONG"}
	}
	if len(args) == 1 {
		return resp.BulkString{
			Length: int64(len(args[0])),
			Value:  args[0],
		}
	}
	return resp.SimpleError{Message: "ERR wrong number of arguments for 'ping'"}
}

type CommandFunc func(d *Dispatcher, args [][]byte) resp.RESPType

func (d *Dispatcher) handleSet(args [][]byte) resp.RESPType {
	if len(args) != 2 {
		return resp.SimpleError{Message: "ERR wrong number of arguments for 'set'"}
	}

	d.dataStore.Set(string(args[0]), string(args[1]))
	return resp.SimpleString{Value: "OK"}
}

func (d *Dispatcher) handleGet(args [][]byte) resp.RESPType {
	if len(args) != 1 {
		return resp.SimpleError{Message: "ERR wrong number of arguments for 'get'"}
	}

	value, exists := d.dataStore.Get(string(args[0]))
	if !exists {
		return resp.BulkString{Length: -1, Value: nil}
	}

	return resp.BulkString{
		Length: int64(len(value)),
		Value:  []byte(value),
	}
}

func (d *Dispatcher) handleEcho(args [][]byte) resp.RESPType {
	if len(args) == 0 {
		return resp.BulkString{Length: -1, Value: []byte{}}
	}
	return resp.BulkString{Length: int64(len(args[0])), Value: args[0]}
}

type Dispatcher struct {
	dataStore *datastore.DataStore
	commands  map[Command]CommandFunc
}

func New() *Dispatcher {
	d := &Dispatcher{
		dataStore: datastore.New(),
		commands:  make(map[Command]CommandFunc),
	}

	d.commands[COMMAND_PING] = (*Dispatcher).handlePing
	d.commands[COMMAND_SET] = (*Dispatcher).handleSet
	d.commands[COMMAND_GET] = (*Dispatcher).handleGet
	d.commands[COMMAND_ECHO] = (*Dispatcher).handleEcho

	return d
}

func (d *Dispatcher) Execute(cmd resp.RESPType) (resp.RESPType, error) {

	arr, ok := cmd.(resp.Array)
	if !ok || len(arr.Elements) == 0 {
		return resp.SimpleError{Message: "ERR invalid command format"}, nil
	}

	cmdBulk, ok := arr.Elements[0].(resp.BulkString)
	if !ok {
		return resp.SimpleError{Message: "ERR invalid command format"}, nil
	}

	commandName := strings.ToUpper(string(cmdBulk.Value))
	commandEnum := Command(commandName)

	commandFunc, ok := d.commands[commandEnum]
	if !ok {
		return resp.SimpleError{Message: "ERR unknown command"}, nil
	}

	args := make([][]byte, 0, len(arr.Elements)-1)
	for i := 1; i < len(arr.Elements); i++ {
		bulk, ok := arr.Elements[i].(resp.BulkString)
		if !ok {
			return resp.SimpleError{Message: "ERR invalid argument type"}, nil
		}
		args = append(args, bulk.Value)
	}

	return commandFunc(d, args), nil
}
