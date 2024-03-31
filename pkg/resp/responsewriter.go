package resp

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync/atomic"
)

const (
	OK = "OK"
)

type responseWriter struct {
	buffer    *bytes.Buffer
	accepting *atomic.Bool
}

func (r *responseWriter) write(p []byte) {
	_, _ = r.buffer.Write(p)
}

func (r *responseWriter) FlushTo(writer io.Writer) (int, error) {
	data := r.buffer.Bytes()
	r.buffer.Reset()
	return writer.Write(data)
}

func (r *responseWriter) Disconnect() {
	r.accepting.Store(false)
}

func (r *responseWriter) writeAll(slices ...[]byte) {
	for _, slice := range slices {
		r.write(slice)
	}
}

func (r *responseWriter) WriteSimpleString(s string) {
	s = strings.ReplaceAll(s, string(CRLF), " ")
	r.writeAll([]byte{'+'}, []byte(s), CRLF)
}

func (r *responseWriter) WriteVerbatimString(s string) {
	data := []byte(s)
	length := len(data) + 4
	r.writeAll([]byte{'='}, []byte(strconv.Itoa(length)), CRLF, []byte("txt:"), data, CRLF)
}

func (r *responseWriter) WriteBulkString(bytes []byte) {
	if bytes == nil {
		r.writeAll([]byte{'$'}, []byte(strconv.Itoa(-1)), CRLF)
		return
	}
	r.writeAll([]byte{'$'}, []byte(strconv.Itoa(len(bytes))), CRLF, bytes, CRLF)
}

func (r *responseWriter) WriteSimpleError(prefix ErrorPrefix, err string) {
	r.writeAll([]byte{'-'}, []byte(prefix), []byte{' '}, []byte(err), CRLF)
}

func (r *responseWriter) WriteBulkError(prefix ErrorPrefix, err string) {
	bytes := []byte(fmt.Sprintf("%s %s", prefix, err))
	r.writeAll([]byte{'!'}, []byte(strconv.Itoa(len(bytes))), CRLF, bytes, CRLF)
}

func (r *responseWriter) WriteInteger(i int) {
	prefix := []byte{':'}
	if i < 0 {
		prefix = append(prefix, '-')
	}
	numBytes := []byte(strconv.Itoa(i))
	r.writeAll(prefix, numBytes, CRLF)
}

func (r *responseWriter) WriteArray(size int) {
	r.writeAll([]byte{'*'}, []byte(strconv.Itoa(size)), CRLF)
}

func (r *responseWriter) WriteNullArray() {
	r.writeAll([]byte("*-1"), CRLF)
}

func (r *responseWriter) WriteMap(size int) {
	r.writeAll([]byte{'%'}, []byte(strconv.Itoa(size)), CRLF)
}

func (r *responseWriter) WriteSet(size int) {
	r.writeAll([]byte{'~'}, []byte(strconv.Itoa(size)), CRLF)
}

func (r *responseWriter) WritePush(size int) {
	r.writeAll([]byte{'>'}, []byte(strconv.Itoa(size)), CRLF)
}

func (r *responseWriter) WriteNull() {
	r.writeAll([]byte{'_'}, CRLF)
}

func (r *responseWriter) WriteBool(b bool) {
	var sign byte = 't'
	if !b {
		sign = 'f'
	}
	r.writeAll([]byte{'#', sign}, CRLF)
}

func (r *responseWriter) WriteDouble(f float64) {
	prefix := []byte{','}
	if f < 0 {
		prefix = append(prefix, '-')
	}
	data := strconv.FormatFloat(f, 'f', -1, 64)
	r.writeAll(prefix, []byte(data), CRLF)
}

func (r *responseWriter) WriteBigint(i int64) {
	prefix := []byte{'('}
	if i < 0 {
		prefix = append(prefix, '-')
	}
	r.writeAll(prefix, []byte(strconv.Itoa(int(i))), CRLF)
}

func (r *responseWriter) WriteOK() {
	r.WriteSimpleString(OK)
}
