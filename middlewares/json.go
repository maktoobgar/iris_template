package middlewares

import "github.com/kataras/iris/v12"

func Json(ctx iris.Context) {
	ctx.Header("Content-Type", "application/json; charset=utf-8")

	ctx.Next()
}
