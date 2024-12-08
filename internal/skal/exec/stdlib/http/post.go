package http

import (
	"bytes"
	"io"
	"net/http"

	"github.com/illbjorn/skal/internal/skal/exec/stdlib/argv"
	"github.com/illbjorn/skal/internal/skal/exec/stdlib/conv"
	"github.com/illbjorn/skal/pkg/clog"
	lua "github.com/yuin/gopher-lua"
)

const (
	mimeTypeJSON = "application/json"
)

func post(l *lua.LState) int {
	var (
		lbody       = l.Get(-1)
		buf         = bytes.NewBuffer(nil)
		err         error
		contentType = l.Get(-2)
		url         = l.Get(-3)
	)

	if !argv.IsString(url) {
		return 0
	}

	if !argv.IsString(contentType) {
		return 0
	}

	if !argv.IsNil(lbody) {
		switch contentType.String() {
		case mimeTypeJSON:
			switch b := lbody.(type) {
			case *lua.LTable:
				buf.WriteString(conv.ToJSON(b))
			}
		}
	}

	var (
		res *http.Response
	)
	if res, err = client.Post(
		url.String(),
		contentType.String(),
		buf); err != nil {
		clog.Debug(
			"Failed to perform POST request.",
			"error", err,
		)

		return 0
	}
	defer res.Body.Close()

	var (
		content []byte
	)
	if content, err = io.ReadAll(res.Body); err != nil {
		clog.Debug(
			"Failed to read POST response body.",
			"error", err,
		)

		return 0
	}

	if res.StatusCode >= 400 {
		clog.Error(
			"Received >=400 status code response",
			"status code", res.StatusCode,
			"body", string(content),
		)
		return 0
	}

	// TODO: Handle by `Content-Type` header.
	l.Push(
		lua.LString(
			string(
				content,
			),
		),
	)

	return 1
}
