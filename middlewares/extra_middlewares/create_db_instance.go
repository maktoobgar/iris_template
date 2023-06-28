package extra_middlewares

import (
	g "service/global"
	"service/pkg/errors"

	"github.com/kataras/iris/v12"
)

func CreateDbInstance(ctx iris.Context) {
	db, err := g.DB()
	if err != nil {
		panic(errors.New(errors.ServiceUnavailable, errors.Resend, "DbNotFound", err.Error(), nil))
	}
	defer db.Close()

	ctx.Values().Set(g.DbInstance, db)

	ctx.Next()
}
