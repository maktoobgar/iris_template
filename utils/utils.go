package utils

import (
	g "service/global"
	"service/pkg/errors"
	"service/pkg/translator"
	"sync"

	"github.com/kataras/iris/v12"
)

func SendJson(ctx iris.Context, v any, status ...int) {
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

func Panic500(err error) {
	panic(errors.New(errors.UnexpectedStatus, errors.Resend, "InternalServerError", err.Error(), nil))
}

func SendMessage(ctx iris.Context, translate translator.TranslatorFunc, message string, data map[string]any) {
	data["message"] = translate(message)
	ctx.JSON(data)
}
