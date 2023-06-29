package handlers

import (
	g "service/global"
	"service/models"

	"github.com/kataras/iris/v12"
)

func Me(ctx iris.Context) {
	user := ctx.Values().Get(g.UserKey).(*models.User)

	ctx.JSON(user)
}
