package utils

import (
	"reflect"
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

func SendPage(ctx iris.Context, dataCount int64, perPage int, page int, data any) {
	pagesCount := CalculatePagesCount(dataCount, perPage)
	dataValue := reflect.ValueOf(data)
	if dataValue.Type().Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}
	len := dataValue.Len()

	ctx.JSON(map[string]any{
		"page":        page,
		"per_page":    perPage,
		"pages_count": pagesCount,
		"all_count":   dataCount,
		"count":       len,
		"data":        data,
	})
}

func CalculatePagesCount(dataCount int64, perPage int) int {
	pagesCount := int64(-1)
	if dataCount%int64(perPage) == 0 {
		pagesCount = dataCount / int64(perPage)
	} else {
		pagesCount = (dataCount / int64(perPage)) + 1
	}
	return int(pagesCount)
}

func Min(v1 int, v2 int) int {
	if v1 < v2 {
		return v1
	} else if v2 < v1 {
		return v2
	} else {
		return v1
	}
}
