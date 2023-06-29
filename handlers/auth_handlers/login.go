package auth_handlers

import (
	"database/sql"
	"service/dto"
	g "service/global"
	"service/handlers"
	"service/models"
	"service/pkg/copier"
	"service/pkg/errors"
	"service/pkg/translator"

	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/kataras/iris/v12"
)

func Login(ctx iris.Context) {
	translate := ctx.Value(g.TranslateKey).(translator.TranslatorFunc)
	req := ctx.Values().Get(g.RequestBody).(*dto.LoginRequest)
	db := ctx.Values().Get(g.DbInstance).(*sql.DB)
	user := models.NewUser()
	copier.Copy(user, req)

	err := user.Select(map[string]any{
		"phone_number": req.PhoneNumber,
	}).QueryRowContextError(ctx, db)
	if err != nil {
		if sqlscan.NotFound(err) {
			panic(errors.New(errors.InvalidStatus, errors.ReSignIn, "UserWithPhoneNumberNotFound", err.Error()))
		} else {
			handlers.Panic500(err)
		}
	}

	if !user.IsPasswordEqualToMyHash(req.Password) {
		panic(errors.New(errors.InvalidStatus, errors.ReSignIn, "PasswordOrPhoneNumberDoNotMatch", "password didn't match"))
	}

	accessToken := user.CreateAccessToken(ctx, db)
	refreshToken := user.CreateRefreshToken(ctx, db)

	handlers.SendMessage(ctx, translate, "Welcome", map[string]any{
		"access_token":  accessToken.Token,
		"refresh_token": refreshToken.Token,
	})
}
