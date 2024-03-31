package resp_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/treeverse/lakefs/pkg/resp"
)

func TestParseCommands(t *testing.T) {
	cases := []struct {
		Name             string
		Data             []byte
		ExpectedError    error
		ExpectedCommands [][]byte
	}{
		{
			Name:             "regular_command",
			Data:             []byte("*2\r\n$4\r\nLLEN\r\n$6\r\nmylist\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{[]byte("LLEN"), []byte("mylist")},
		},
		{
			Name:             "command_with_crlf_string",
			Data:             []byte("*2\r\n$4\r\nLLEN\r\n$8\r\nmylis\r\nt\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{[]byte("LLEN"), []byte("mylis\r\nt")},
		},
		{
			Name:             "command_with_crlf_String_and_more",
			Data:             []byte("*3\r\n$4\r\nLLEN\r\n$8\r\nmylis\r\nt\r\n$14\r\nabcdefghijklmn\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{[]byte("LLEN"), []byte("mylis\r\nt"), []byte("abcdefghijklmn")},
		},
		{
			Name:             "empty_commands",
			Data:             []byte("*0\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{},
		},
		{
			Name:             "too_many_commands",
			Data:             []byte("*0\r\n$4\r\nLLEN\r\n$8\r\nmylis\r\nt\r\n$14\r\nabcdefghijklmn\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{},
		},
		{
			Name:          "too_few_commands",
			Data:          []byte("*3\r\n$4\r\nLLEN\r\n$8\r\nmylis\r\nt\r\n$14\r\n"),
			ExpectedError: resp.ErrProtocol,
		},
		{
			Name:             "inline_simple",
			Data:             []byte("GET foo bar\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{[]byte("GET"), []byte("foo"), []byte("bar")},
		},
		{
			Name:             "inline_single_quote",
			Data:             []byte("GET \"foo bar\" bazbaz\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{[]byte("GET"), []byte("foo bar"), []byte("bazbaz")},
		},
		{
			Name:             "inline_double_quote",
			Data:             []byte("GET foo bar 'baz     baz'\r\n"),
			ExpectedError:    nil,
			ExpectedCommands: [][]byte{[]byte("GET"), []byte("foo"), []byte("bar"), []byte("baz     baz")},
		},
		{
			Name:          "inline_unclosed_quotes",
			Data:          []byte("GET \"foo bar\r\n"),
			ExpectedError: resp.ErrProtocol,
		},
	}
	for _, cas := range cases {
		t.Run(cas.Name, func(t *testing.T) {
			r := bytes.NewReader(cas.Data)
			ctx := context.Background()
			req, err := resp.ParseRequest(ctx, r)
			if !errors.Is(err, cas.ExpectedError) {
				t.Errorf("unexpected error parsing commands: %v\n", err)
				return
			}
			if cas.ExpectedError != nil {
				return
			}
			// compare output
			cmds := req.Commands()
			if len(cmds) != len(cas.ExpectedCommands) {
				t.Errorf("expected %d cmds got %d", len(cas.ExpectedCommands), len(cmds))
				return
			}
			for i := 0; i < len(cmds); i++ {
				if !bytes.Equal(cmds[i], cas.ExpectedCommands[i]) {
					t.Errorf("expected cmds[%d] = '%s' got '%s'", i, cas.ExpectedCommands[i], cmds[i])
				}
			}

		})
	}
}
