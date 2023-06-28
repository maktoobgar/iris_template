package handlers

import (
	"github.com/kataras/iris/v12"
)

func Hello(ctx iris.Context) {
	SendJson(ctx, map[string]string{
		"msg": "Hello World 🥳",
	})
}
