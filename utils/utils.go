package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	g "service/global"
	"service/pkg/errors"
	"service/pkg/translator"
	"sync"

	"github.com/golodash/galidator"
	"github.com/kataras/iris/v12"
)

func sendIfCtxNotCancelled(ctx iris.Context, status int, value any) {
	if ctx.Err() == nil {
		writerLock := ctx.Values().Get(g.WriterLock).(*sync.Mutex)
		writerLock.Lock()
		defer writerLock.Unlock()

		closedWriter := ctx.Values().Get(g.ClosedWriter).(bool)
		if !closedWriter {
			if status != -1 {
				ctx.StatusCode(status)
			}
			ctx.JSON(value)
		}

		closedWriter = true
		ctx.Values().Set(g.ClosedWriter, closedWriter)
	}
}

func SendJsonMessage(ctx iris.Context, translate translator.TranslatorFunc, message string, data map[string]any, status ...int) {
	code := 200
	if len(status) > 0 {
		code = status[0]
	}
	data["message"] = translate(message)

	sendIfCtxNotCancelled(ctx, code, data)
}

func SendJson(ctx iris.Context, data any, status ...int) {
	code := 200
	if len(status) > 0 {
		code = status[0]
	}

	sendIfCtxNotCancelled(ctx, code, data)
}

func SendEmpty(ctx iris.Context, status ...int) {
	code := 200
	if len(status) > 0 {
		code = status[0]
	}

	sendIfCtxNotCancelled(ctx, code, nil)
}

func Panic500(err error) {
	panic(errors.New(errors.UnexpectedStatus, "InternalServerError", err.Error(), nil))
}

func SendMessage(ctx iris.Context, translate translator.TranslatorFunc, message string, data map[string]any) {
	data["message"] = translate(message)
	sendIfCtxNotCancelled(ctx, -1, data)
}

func SendPage(ctx iris.Context, dataCount int64, perPage int, page int, data any) {
	pagesCount := CalculatePagesCount(dataCount, perPage)
	if page > pagesCount {
		panic(errors.New(errors.NotFoundStatus, "PageNotFound", fmt.Sprintf("page %d requested but we have %d pages", page, pagesCount)))
	}
	dataValue := reflect.ValueOf(data)
	if dataValue.Type().Kind() == reflect.Ptr {
		dataValue = dataValue.Elem()
	}

	sendIfCtxNotCancelled(ctx, -1, map[string]any{
		"page":        page,
		"per_page":    perPage,
		"pages_count": pagesCount,
		"all_count":   dataCount,
		"count":       dataValue.Len(),
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

	// If there is no date, just return 1 page so that NotFound do not get returned
	if int(pagesCount) == 0 {
		return 1
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

func Validate(data any, validator galidator.Validator, translate translator.TranslatorFunc) {
	if errs := validator.Validate(data, galidator.Translator(translate)); errs != nil {
		panic(errors.New(errors.InvalidStatus, "BodyNotProvidedProperly", "", errs))
	}
}

func PrettyJsonBytes(data []byte) string {
	var prettyJSON bytes.Buffer
	if err := json.Indent(&prettyJSON, data, "", "  "); err != nil {
		return ""
	}
	return prettyJSON.String()
}
