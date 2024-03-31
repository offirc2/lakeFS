package resp

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strconv"
)

type request struct {
	commands [][]byte
	ctx      context.Context
}

func (r *request) Context() context.Context {
	return r.ctx
}

func (r *request) SetContext(ctx context.Context) {
	r.ctx = ctx
}

func (r *request) Commands() [][]byte {
	return r.commands
}

func readCmd(reader *bufio.Reader) ([]byte, error) {
	// read descriptor
	data, err := reader.ReadBytes(CR)
	if err != nil {
		return nil, ErrProtocol
	}
	if lf, err := reader.ReadByte(); err != nil || lf != LF {
		return nil, ErrProtocol
	}
	line := data[0 : len(data)-1] // remove CR
	if line[0] != '$' {
		return nil, ErrProtocol
	}
	strSize, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, ErrProtocol
	}
	buf := make([]byte, strSize+len(CRLF))
	n, err := reader.Read(buf)
	if err != nil || n < len(buf) {
		return nil, ErrProtocol
	}
	return buf[0 : len(buf)-len(CRLF)], nil
}

func parseInlineCommand(line []byte) ([][]byte, error) {
	const singleQuote = '\''
	const doubleQuote = '"'
	var inSingleQuotes, inDoubleQuotes bool
	cmds := make([][]byte, 0)
	var currentCmd []byte
	for _, c := range line {
		switch c {
		case SPACE:
			if !inSingleQuotes && !inDoubleQuotes {
				cmds = append(cmds, currentCmd)
				currentCmd = nil
			} else if inSingleQuotes || inDoubleQuotes {
				currentCmd = append(currentCmd, c)
			}
		case singleQuote:
			inSingleQuotes = !inSingleQuotes
		case doubleQuote:
			inDoubleQuotes = !inDoubleQuotes
		default:
			currentCmd = append(currentCmd, c)
		}
	}

	if inDoubleQuotes {
		return nil, fmt.Errorf("%w: unbalanced double quotes", ErrProtocol)
	}
	if inSingleQuotes {
		return nil, fmt.Errorf("%w: unbalanced single quotes", ErrProtocol)
	}
	if currentCmd != nil {
		cmds = append(cmds, currentCmd)
	}
	return cmds, nil
}

func ParseRequest(ctx context.Context, r io.Reader) (Request, error) {
	/*
		Example input, taken from:
		 https://redis.io/docs/reference/protocol-spec/#sending-commands-to-a-redis-server
		  C: *2\r\n       *2 = array of size 2
		  C: $4\r\n       $4 = bulk string of size 4
		  C: LLEN\r\n     bulk string itself
		  C: $6\r\n       $6 = bulk string of size 6
		  C: mylist\r\n   bulk string itself #2
	*/
	reader := bufio.NewReader(r)
	data, err := reader.ReadBytes(CR)
	if err != nil {
		return nil, ErrProtocol
	}
	if lf, err := reader.ReadByte(); err != nil || lf != LF {
		return nil, ErrProtocol
	}
	line := data[0 : len(data)-1] // remove CR
	// read cmd size
	instruction := line[0]
	if instruction != '*' {
		// let's attempt to parse an inline command:
		// https://redis.io/docs/reference/protocol-spec/#inline-commands
		cmds, err := parseInlineCommand(line)
		if err != nil {
			return nil, err
		}
		return &request{
			commands: cmds,
			ctx:      ctx,
		}, nil
	}
	cmdSize, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, ErrProtocol
	}
	cmds := make([][]byte, cmdSize)
	// start reading cmds
	for i := 0; i < len(cmds); i++ {
		var err error
		cmds[i], err = readCmd(reader)
		if err != nil {
			return nil, err
		}
	}
	return &request{commands: cmds, ctx: ctx}, nil
}
