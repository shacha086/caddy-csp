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
		return nil
	}
	//cspHeader := writer.Header().Get("Content-Security-Policy")
	return nil
}

func (*CaddyCSPHandler) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.csp",
		New: func() caddy.Module { return new(CaddyCSPHandler) },
	}
}
