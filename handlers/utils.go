package handlers

import (
	g "service/global"
	"sync"

	"github.com/kataras/iris/v12"
)

func sendJson(ctx iris.Context, v any, status ...int) {
	code := 200
	if len(status) > 0 {
		code = status[0]
	}

	if ctx.Err() == nil {
		writerLock := ctx.Values().Get(g.WriterLock).(*sync.Mutex)
		writerLock.Lock()
		defer writerLock.Unlock()

		closedWriter := ctx.Values().Get(g.ClosedWriter).(bool)
		if !closedWriter {
			ctx.StatusCode(code)
			ctx.JSON(v)
		}

		closedWriter = true
		ctx.Values().Set(g.ClosedWriter, closedWriter)
	}
}
