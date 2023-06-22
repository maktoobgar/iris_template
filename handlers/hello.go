package handlers

import "github.com/kataras/iris/v12"

func Hello(ctx iris.Context) {
	ctx.JSON(map[string]string{
		"msg": "Hello World 🥳",
	})
}
