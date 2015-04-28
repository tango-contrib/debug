// Copyright 2015 The Tango Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package debug

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/lunny/tango"
)

type bufferWriter struct {
	tango.ResponseWriter
	content []byte
}

func (b *bufferWriter) Write(bs []byte) (int, error) {
	b.content = append(b.content, bs...)
	return b.ResponseWriter.Write(bs)
}

type Options struct {
	HideRequest      bool
	HideRequestHead  bool
	HideRequestBody  bool
	HideResponse     bool
	HideResponseHead bool
	HideResponseBody bool
	IgnorePrefix     string
}

func prepareOptions(opts []Options) Options {
	var opt Options
	if len(opts) > 0 {
		opt = opts[0]
	}
	return opt
}

func Debug(options ...Options) tango.HandlerFunc {
	return func(ctx *tango.Context) {
		opt := prepareOptions(options)

		if opt.HideRequest && opt.HideResponse {
			ctx.Next()
			return
		}
		if len(opt.IgnorePrefix) > 0 && strings.HasPrefix(ctx.Req().URL.Path, opt.IgnorePrefix) {
			ctx.Next()
			return
		}

		if !opt.HideRequest {
			ctx.Debug("[debug] request:", ctx.Req().Method, ctx.Req().URL, ctx.Req().RemoteAddr)
			if !opt.HideRequestHead {
				ctx.Debug("[debug] head:", ctx.Req().Header)
			}
			if !opt.HideRequestBody {
				if ctx.Req().Body != nil {
					requestbody, _ := ioutil.ReadAll(ctx.Req().Body)
					ctx.Req().Body.Close()
					bf := bytes.NewBuffer(requestbody)
					ctx.Req().Body = ioutil.NopCloser(bf)
					ctx.Debug("[debug] body:", string(requestbody))
				}
			}
			ctx.Debug("[debug] ----------------------- end request")
		}

		buf := &bufferWriter{
			ctx.ResponseWriter,
			make([]byte, 0),
		}
		ctx.ResponseWriter = buf

		ctx.Next()

		if !opt.HideResponse {
			ctx.Debug("[debug] response ------------------", ctx.Status())
			if !opt.HideRequestHead {
				ctx.Debug("[debug] head:", buf.ResponseWriter.Header())
			}
			if !opt.HideResponseBody {
				ctx.Debug("[debug] body:", string(buf.content))
			}
			ctx.Debug("[debug] ----------------------- end response")
		}

		ctx.ResponseWriter = buf.ResponseWriter
	}
}
