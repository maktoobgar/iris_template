package middlewares

import (
	"log"
	"net/http"
	"runtime/debug"
	"sync"

	"service/dto"
	g "service/global"

	"service/pkg/errors"
	"service/pkg/translator"

	"github.com/kataras/iris/v12"
)

func Panic(ctx iris.Context) {
	translate := ctx.Value(g.TranslateKey).(translator.TranslatorFunc)

	defer func() {
		errInterface := recover()
		if errInterface == nil {
			return
		}

		writerLock := ctx.Values().Get(g.WriterLock).(*sync.Mutex)
		writerLock.Lock()
		defer writerLock.Unlock()

		closedWriter := ctx.Values().Get(g.ClosedWriter).(bool)
		if !closedWriter {
			if err, ok := errInterface.(error); ok && errors.IsServerError(err) {
				code, action, message, _, errors := errors.HttpError(err)
				res := dto.PanicResponse{
					Message: translate(message),
					Code:    code,
					Action:  action,
					Errors:  errors,
				}
				if g.CFG.Debug {
					log.Println(err)
				}
				ctx.StopWithJSON(code, res)
			} else {
				stack := string(debug.Stack())
				g.Logger.Panic(errInterface, ctx.Request(), stack)
				res := dto.PanicResponse{
					Message: translate("InternalServerError"),
					Action:  int(errors.Report),
					Code:    http.StatusInternalServerError,
					Errors:  nil,
				}
				ctx.StopWithJSON(res.Code, res)
			}
		}
		closedWriter = true
		ctx.Values().Set(g.ClosedWriter, closedWriter)
	}()

	ctx.Next()
}
