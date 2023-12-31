package middlewares

import (
	"encoding/json"
	"io"
	"reflect"
	g "service/global"
	"service/pkg/errors"
	"service/pkg/translator"

	"github.com/golodash/galidator"
	"github.com/kataras/iris/v12"
)

// Parses and validates request body
func Validate(validator galidator.Validator, inputInstance any) iris.Handler {
	return func(ctx iris.Context) {
		req := reflect.New(reflect.TypeOf(inputInstance)).Interface()
		// Parse body and check for errors
		body := ctx.Request().Body
		bytes, err1 := io.ReadAll(body)
		err2 := json.Unmarshal(bytes, req)

		if err1 != nil {
			panic(errors.New(errors.InvalidStatus, "BodyNotProvidedProperly", err1.Error()))
		} else if err2 != nil {
			panic(errors.New(errors.InvalidStatus, "BodyNotProvidedProperly", err2.Error()))
		}

		// Validate and translate error messages if errors exist
		translate := ctx.Value(g.TranslateKey).(translator.TranslatorFunc)
		errs := validator.Validate(req, galidator.Translator(translate))
		if errs != nil {
			panic(errors.New(errors.InvalidStatus, "BodyNotProvidedProperly", "", errs))
		}

		// If we come this far, data is valid, so record it in context
		ctx.Values().Set(g.RequestBody, req)

		ctx.Next()
	}
}
