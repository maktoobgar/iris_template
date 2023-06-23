package handlers

import (
	"github.com/kataras/iris/v12"
)

func Hello(ctx iris.Context) {
	sendJson(ctx, map[string]string{
		"msg": "Hello World 🥳",
	})
}
