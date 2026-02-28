package dispatcher

type Command string

const (
	COMMAND_PING Command = "PING"
	COMMAND_SET  Command = "SET"
	COMMAND_GET  Command = "GET"
	COMMAND_ECHO Command = "ECHO"
)
