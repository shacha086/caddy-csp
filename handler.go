package csp

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"net/http"
)

const LiteralAll = "!all"

func init() {
	caddy.RegisterModule(&CaddyCSPHandler{})
	httpcaddyfile.RegisterHandlerDirective("csp", parseCaddyFile)
}

func (h *CaddyCSPHandler) Provision(ctx caddy.Context) error {
	h.logger = ctx.Logger(h)
	return nil
}

func parseCaddyFile(helper httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	csp := &CaddyCSPHandler{
		Enabled: false,
		Cmd:     nil,
	}
	for helper.Next() {
		for helper.NextBlock(0) {
			cmd := helper.Val()
			args := helper.RemainingArgs()
			if len(args) < 2 {
				continue
			}
			k := args[0]
			v := args[1:]

			if cmd == "add" || cmd == "remove" || cmd == "set" {
				csp.Enabled = true
				csp.Cmd = append(csp.Cmd, CmdType{Op: cmd, Key: k, Values: v})
			}
		}
	}
	return csp, nil
}

type CmdType struct {
	Op     string // "add", "remove", "set"
	Key    string
	Values []string
}

type CaddyCSPHandler struct {
	Enabled bool
	logger  *zap.Logger
	Cmd     []CmdType
}
type headerInterceptor struct {
	http.ResponseWriter
	wrote  bool
	logger *zap.Logger
	Cmd    []CmdType
}

func (h *headerInterceptor) WriteHeader(statusCode int) {
	if !h.wrote {
		cspHeader := h.Header().Get("Content-Security-Policy")
		csp, warn := NewCSPFromHeader(cspHeader)
		if warn != nil {
			h.logger.Warn(warn.Error())
		}
		if csp != nil {
			h.ApplyAll(csp)
			encoded := csp.Encoded()
			if encoded != "" {
				h.Header().Set("Content-Security-Policy", encoded)
			}
		}
		h.wrote = true
	}
	h.ResponseWriter.WriteHeader(statusCode)
}

func (h *headerInterceptor) Write(b []byte) (int, error) {
	if !h.wrote {
		h.WriteHeader(http.StatusOK)
	}
	return h.ResponseWriter.Write(b)
}

func (h *CaddyCSPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, next caddyhttp.Handler) error {
	if !h.Enabled {
		return next.ServeHTTP(writer, request)
	}
	interceptor := &headerInterceptor{
		ResponseWriter: writer,
		Cmd:            h.Cmd,
		logger:         h.logger,
	}
	return next.ServeHTTP(interceptor, request)
}

func (h *headerInterceptor) ApplyAll(csp *CSP) {
	for _, cmd := range h.Cmd {
		k, v := cmd.Key, cmd.Values
		switch cmd.Op {
		case "add":
			if k != LiteralAll {
				csp.Add(k, v...)
			} else {
				csp.AddAllDirective(v...)

			}
		case "remove":
			if k != LiteralAll {
				csp.Remove(k, v...)
			} else {
				csp.RemoveAllDirective(v...)
			}
		case "set":
			if k != LiteralAll {
				csp.Set(k, v...)
			} else {
				csp.SetAllDirective(v...)
			}
		}
	}
}

func (*CaddyCSPHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.csp",
		New: func() caddy.Module { return new(CaddyCSPHandler) },
	}
}
