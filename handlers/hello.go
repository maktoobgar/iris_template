package handlers

import "github.com/kataras/iris/v12"

func Hello(ctx iris.Context) {
	ctx.HTML("Hello <strong>%s</strong>!", "World")
}
