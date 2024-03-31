package resp

import (
	"bytes"
	"context"
	"errors"
	"net"

	"sync/atomic"
)

const (
	SPACE = ' '
	CR    = '\r'
	LF    = '\n'
)

var (
	CRLF = []byte{CR, LF}
)

type Request interface {
	Commands() [][]byte
	Context() context.Context
	SetContext(ctx context.Context)
}

type ResponseWriter interface {
	WriteSimpleString(string)
	WriteVerbatimString(string)
	WriteBulkString([]byte)
	WriteSimpleError(ErrorPrefix, string)
	WriteBulkError(ErrorPrefix, string)
	WriteInteger(int)
	WriteArray(size int)
	WriteNullArray()
	WriteMap(size int)
	WriteSet(size int)
	WritePush(size int)
	WriteNull()
	WriteBool(bool)
	WriteDouble(float64)
	WriteBigint(int64)
	WriteOK()
	Disconnect()
}

type Handler interface {
	Handle(Request, ResponseWriter)
}

type Server struct {
	listener net.Listener
	handler  Handler
	running  *atomic.Bool
}

func (s *Server) Run() error {
	s.running.Store(true)
	for s.running.Load() {
		// Accept incoming connections
		conn, err := s.listener.Accept()
		if err != nil {
			// TODO (ozk): log
		} else {
			// Handle client connection in a goroutine
			go s.handleClient(conn)
		}
	}
	return nil
}

func (s *Server) Shutdown(_ context.Context) error {
	s.running.Store(false)
	return s.listener.Close()
}

func (s *Server) handleClient(conn net.Conn) {
	accepting := &atomic.Bool{}
	accepting.Store(true)
	ctx := context.Background()
	for accepting.Load() {
		// parse request
		req, err := ParseRequest(ctx, conn)
		w := &responseWriter{accepting: accepting, buffer: bytes.NewBuffer(nil)}
		if errors.Is(err, ErrProtocol) {
			w.WriteSimpleError(ErrorPrefixGeneric, err.Error())
			_, _ = w.FlushTo(conn)
			break
		}

		s.handler.Handle(req, w)
		_, err = w.FlushTo(conn)
		if err != nil {
			break
		}
		ctx = req.Context()
	}
	_ = conn.Close()
}

func NewServer(listener net.Listener, handler Handler) *Server {
	return &Server{
		listener: listener,
		handler:  handler,
		running:  &atomic.Bool{},
	}
}
