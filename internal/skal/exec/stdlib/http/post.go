package http

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/argv"
	lua "github.com/yuin/gopher-lua"
)

func post(l *lua.LState) int {
	var (
		body = l.Get(-1)
		url  = l.Get(-2)
	)

	if !argv.IsString(url) {
		return 0
	}

	if !argv.IsString(body) {
		return 0
	}

	var res *http.Response
	var err error
	if res, err = client.Post(url.String(), "application/json", strings.NewReader(body.String())); err != nil {
		slog.Debug(
			"Failed to perform POST request.",
			"error", err,
		)

		return 0
	}
	defer res.Body.Close()

	var content []byte
	if content, err = io.ReadAll(res.Body); err != nil {
		slog.Debug(
			"Failed to read POST response body.",
			"error", err,
		)

		return 0
	}

	// TODO: Handle by `Content-Type` header.
	l.Push(lua.LString(string(content)))

	return 1
}
