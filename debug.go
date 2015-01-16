package debug

import (
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

func Debug() tango.HandlerFunc {
	return func(ctx *tango.Context) {
		buf := &bufferWriter{
			ctx.ResponseWriter,
			make([]byte, 0),
		}
		ctx.ResponseWriter = buf

		ctx.Debug("[debug] request:", ctx.Req().Method, ctx.Req().URL, ctx.Req().RemoteAddr)

		ctx.Next()

		ctx.Debug("[debug] result:", ctx.Status(), string(buf.content))

		ctx.ResponseWriter = buf.ResponseWriter
	}
}
