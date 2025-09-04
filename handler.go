package csp

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"net/http"
)

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
		Add:     make(map[string][]string),
		Remove:  make(map[string][]string),
		Set:     make(map[string][]string),
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

			switch cmd {
			case "add":
				csp.Enabled = true
				SetOrAppend(csp.Add, k, v...)
			case "remove":
				csp.Enabled = true
				SetOrAppend(csp.Remove, k, v...)
			case "set":
				csp.Enabled = true
				SetOrAppend(csp.Set, k, v...)
			}
		}
	}

	return csp, nil
}

type CaddyCSPHandler struct {
	Enabled bool
	logger  *zap.Logger
	Add     map[string][]string `json:"add,omitempty"`
	Remove  map[string][]string `json:"remove,omitempty"`
	Set     map[string][]string `json:"set,omitempty"`
}

func (h *CaddyCSPHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request, next caddyhttp.Handler) error {
	if !h.Enabled {
		return next.ServeHTTP(writer, request)
	}
	err := next.ServeHTTP(writer, request)
	if err != nil {
		return err
	}
	cspHeader := writer.Header().Get("Content-Security-Policy")
	csp, warn := NewCSPFromHeader(cspHeader)
	if warn != nil {
		h.logger.Warn(warn.Error())
	}
	if csp == nil {
		return nil
	}

	h.ApplyAdd(csp)
	h.ApplyRemove(csp)
	h.ApplySet(csp)

	writer.Header().Set("Content-Security-Policy", csp.Encoded())
	return nil
}

func (h *CaddyCSPHandler) ApplyAdd(csp *CSP) {
	for k, v := range h.Add {
		if k != "all" {
			csp.Add(k, v...)
		} else {
			csp.AddAllDirective(v...)
		}
	}
}

func (h *CaddyCSPHandler) ApplyRemove(csp *CSP) {
	for k, v := range h.Remove {
		if k != "all" {
			csp.Remove(k, v...)
		} else {
			csp.RemoveAllDirective(v...)
		}
	}
}

func (h *CaddyCSPHandler) ApplySet(csp *CSP) {
	for k, v := range h.Set {
		if k != "all" {
			csp.Set(k, v...)
		} else {
			csp.SetAllDirective(v...)
		}
	}
}

func (*CaddyCSPHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.csp",
		New: func() caddy.Module { return new(CaddyCSPHandler) },
	}
}
