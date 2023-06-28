package auth_handlers

import (
	"database/sql"
	"service/dto"
	g "service/global"
	"service/handlers"
	"service/models"
	"service/pkg/copier"

	"github.com/kataras/iris/v12"
)

func Register(ctx iris.Context) {
	req := ctx.Values().Get(g.RequestBody).(*dto.RegisterRequest)
	db := ctx.Values().Get(g.DbInstance).(*sql.DB)
	user := models.NewUser()
	copier.Copy(user, req)

	// Activate user
	user.IsActive = true

	// Create User
	user.InsertInto().ExecContext(ctx, db)
	user.GetMe().QueryRowContext(ctx, db)
	handlers.SendJson(ctx, user)
}
