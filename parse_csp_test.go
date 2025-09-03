package csp

import (
	"github.com/zeebo/assert"
	"testing"
)

func TestParseCsp(t *testing.T) {
	header := "default-src 'self'; img-src 'self' https: data: blob:; media-src 'self'; script-src 'self' 'unsafe-eval' blob: 'nonce-1nline-m4p' 'sha256-HuyBNEnumn/Bw3njx2R0EXAv9HicWHLQQd9NJ9ruyrk=' 'sha256-R0uX6VU/LU5M8MQi65kWm7KfGnLcUpk2wfubUV6oOsc=' 'sha256-yxwQ9j8YGPsfU554CNGiSCW08z5yqDVvuQmssjoPsm8=' 'sha256-2Q+j4hfT09+1+imS46J2YlkCtHWQt0/BE79PXjJ0ZJ8=' 'sha256-/r7rqQ+yrxt57sxLuQ6AMYcy/lUpvAIzHjIJt/OeLWU='; child-src 'self'; frame-src 'self' https://*.vscode-cdn.net data:; worker-src 'self' data: blob:; style-src 'self' 'unsafe-inline'; connect-src 'self' ws: wss: https:; font-src 'self' blob:; manifest-src 'self';"
	csp, _ := NewCSPFromHeader(header)
	encoded := csp.Encoded()
	assert.Equal(t, encoded, header)
}
