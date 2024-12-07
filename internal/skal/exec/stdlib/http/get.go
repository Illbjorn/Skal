package http

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/argv"
	lua "github.com/yuin/gopher-lua"
)

func get(l *lua.LState) int {
	var (
		url = l.Get(-1)
	)

	if !argv.IsString(url) {
		return 0
	}

	if !argv.IsURL(url) {
		return 0
	}

	var (
		req *http.Request
		err error
	)

	if req, err = http.NewRequest(http.MethodGet, url.String(), nil); err != nil {
		slog.Debug(
			"Failed to create request.",
			"error", err,
		)

		return 0
	}

	var (
		res *http.Response
	)

	if res, err = client.Do(req); err != nil {
		slog.Debug(
			"Failed to perform GET request.",
			"error", err,
		)
	}
	defer res.Body.Close()

	var content []byte
	if content, err = io.ReadAll(res.Body); err != nil {
		slog.Debug(
			"Failed to read GET response body.",
			"error", err,
		)

		return 0
	}

	l.Push(
		lua.LString(string(content)),
	)

	slog.Debug("Response successfully pushed to stack.")

	return 1
}
