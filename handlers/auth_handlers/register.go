package auth_handlers

import (
	"database/sql"
	"net/http"
	"service/dto"
	g "service/global"
	"service/models"
	"service/pkg/copier"
	"service/pkg/translator"
	"service/utils"

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
	user.InsertInto().ExecQuery(ctx, db)
	user.Select(map[string]any{
		"phone_number": req.PhoneNumber,
	}).ExecQueryRow(ctx, db)
	ctx.StatusCode(http.StatusCreated)
	utils.SendMessage(ctx, translate, "RegisterationFinishedSuccessfully", map[string]any{
		"user": user,
	})
}
