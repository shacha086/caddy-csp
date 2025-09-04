package csp

import (
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type fakeHandler struct{}

func (fakeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Security-Policy", "default-src 'self';")
	w.WriteHeader(http.StatusOK)
	return nil
}

func TestCaddyCSPHandler(t *testing.T) {
	// 模拟下游 handler
	next := fakeHandler{}

	// 构造 CSP 中间件
	h := &CaddyCSPHandler{
		Enabled: true,
		Add: map[string][]string{
			"default-src": {src},
		},
		Remove: map[string][]string{
			"default-src": {"'self'"},
		},
		Set: map[string][]string{
			"script-src": {"'none'"},
		},
		logger: zap.NewNop(), // 不输出日志
	}

	// 创建 httptest recorder
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	// 调用中间件
	err := h.ServeHTTP(rec, req, next)
	if err != nil {
		t.Fatal(err)
	}

	// 获取最终 CSP
	csp := rec.Header().Get("Content-Security-Policy")
	t.Logf("Final CSP: %s", csp)

	// 简单断言
	if !strings.Contains(csp, "https://example.com") {
		t.Errorf("expected CSP to contain https://example.com")
	}
	if !strings.Contains(csp, "'none'") {
		t.Errorf("expected CSP to contain 'none'")
	}
	if strings.Contains(csp, "'self'") {
		t.Errorf("expected 'self' to be removed")
	}
}
