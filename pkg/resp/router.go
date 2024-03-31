package resp

import (
	"bytes"
	"fmt"
	"strings"
)

type ValidateFn func(args [][]byte) error
type HandlerFn func(request Request, args [][]byte, w ResponseWriter)

type CommandHandler interface {
	Handle(request Request, args [][]byte, w ResponseWriter)
}

type Command struct {
	Name      string
	Validate  ValidateFn
	HandlerFn HandlerFn
	Handler   CommandHandler
}

type route struct {
	cmd     [][]byte
	handler *Command
}

type Router struct {
	routes []*route
}

func split(cmd string) [][]byte {
	data := []byte(strings.ToLower(cmd))
	return bytes.Split(data, []byte{' '})
}

func (r *Router) Add(command *Command) {
	r.routes = append(r.routes, &route{
		cmd:     split(command.Name),
		handler: command,
	})
}

func (r *Router) Handle(req Request, w ResponseWriter) {
	cmds := req.Commands()
	for _, route := range r.routes {
		cmdParts := len(route.cmd)
		if len(cmds) < cmdParts {
			continue
		}
		match := true
		for i := 0; i < cmdParts; i++ {
			if !bytes.EqualFold(route.cmd[i], cmds[i]) {
				match = false
				break
			}
		}
		if !match {
			continue
		}
		// found a match!
		args := cmds[cmdParts:]
		if route.handler.Validate != nil {
			if err := route.handler.Validate(args); err != nil {
				w.WriteSimpleError(ErrorPrefixGeneric, err.Error())
				return
			}
		}

		// run the thing!
		if route.handler.HandlerFn != nil {
			route.handler.HandlerFn(req, args, w)
		} else {
			route.handler.Handler.Handle(req, args, w)
		}

		return
	}
	// no match
	w.WriteSimpleError(ErrorPrefixGeneric, fmt.Sprintf("command not found"))
}

func NewRouter() *Router {
	return &Router{routes: make([]*route, 0)}
}
