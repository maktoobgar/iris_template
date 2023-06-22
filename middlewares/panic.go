package middlewares

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"

	"service/dto"
	g "service/global"

	"service/pkg/errors"
	"service/pkg/translator"

	"github.com/kataras/iris/v12"
)

func Panic(ctx iris.Context) {
	translate := ctx.Value(g.TranslateKey).(translator.TranslatorFunc)

	defer func() {
		header := ctx.Request().Header
		w := ctx.ResponseWriter()
		errInterface := recover()
		if errInterface == nil {
			return
		}
		if err, ok := errInterface.(error); ok && errors.IsServerError(err) {
			code, action, message, _, errors := errors.HttpError(err)
			res := dto.PanicResponse{
				Message: translate(message),
				Code:    code,
				Action:  action,
				Errors:  errors,
			}
			resBytes, _ := json.Marshal(res)
			if g.CFG.Debug {
				log.Println(err)
			}
			if header.Get("timeout") == "yes" {
				return
			}
			w.WriteHeader(code)
			w.Write(resBytes)
			if code == http.StatusRequestTimeout {
				header.Set("timeout", "yes")
			}
		} else {
			stack := string(debug.Stack())
			g.Logger.Panic(errInterface, ctx.Request(), stack)
			res := dto.PanicResponse{
				Message: translate("InternalServerError"),
				Action:  int(errors.Report),
				Code:    http.StatusInternalServerError,
				Errors:  nil,
			}
			resBytes, _ := json.Marshal(res)
			w.WriteHeader(res.Code)
			w.Write(resBytes)
		}
	}()

	ctx.Next()
}
