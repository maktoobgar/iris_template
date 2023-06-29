package auth_handlers

import (
	"database/sql"
	"service/dto"
	g "service/global"
	"service/handlers"
	"service/models"
	"service/pkg/copier"
	"service/pkg/translator"

	"github.com/kataras/iris/v12"
)

func Register(ctx iris.Context) {
	translate := ctx.Value(g.TranslateKey).(translator.TranslatorFunc)
	req := ctx.Values().Get(g.RequestBody).(*dto.RegisterRequest)
	db := ctx.Values().Get(g.DbInstance).(*sql.DB)
	user := models.NewUser()
	copier.Copy(user, req)

	// Activate user & hash the password
	user.HashMyPassword()
	user.IsActive = true

	// Create User
	user.InsertInto().ExecContext(ctx, db)
	user.GetMe().QueryRowContext(ctx, db)
	handlers.SendMessage(ctx, translate, "RegisterationFinishedSuccessfully", map[string]any{
		"user": user,
	})
}
