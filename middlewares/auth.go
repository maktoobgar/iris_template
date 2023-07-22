package middlewares

import (
	"database/sql"
	"regexp"
	g "service/global"
	"service/models"
	"service/pkg/errors"
	"service/utils"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/kataras/iris/v12"
)

var tokenPattern, _ = regexp.Compile(`^\d+\|.*$`)

func Auth(ctx iris.Context) {
	db := ctx.Values().Get(g.DbInstance).(*sql.DB)

	// Check if a token has sent
	tokenString := ctx.GetHeader(g.AccessToken)
	if tokenString == "" {
		tokenString = ctx.GetCookie(g.AccessToken)
	}
	if !tokenPattern.MatchString(tokenString) {
		panic(errors.New(errors.UnauthorizedStatus, "LoginPlease", "sent token is not valid"))
	}

	// Get the actual token
	tokenSplit := strings.Split(tokenString, "|")
	tokenId, _ := strconv.ParseInt(tokenSplit[0], 10, 64)
	tokenString = strings.Replace(tokenString, tokenSplit[0]+"|", "", 1)

	// Check if token is valid and decrypt if so
	claims := &models.Claims{}
	tkn, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return g.SecretKeyBytes, nil
	})
	if err != nil {
		panic(errors.New(errors.UnauthorizedStatus, "LoginPlease", err.Error()))
	}
	if !tkn.Valid {
		panic(errors.New(errors.UnauthorizedStatus, "LoginPlease", "token is invalid"))
	}
	if claims.Type != models.AccessTokenType {
		panic(errors.New(errors.UnauthorizedStatus, "LoginPlease", "token is not access token"))
	}
	if claims.ExpiresAt < time.Now().Unix() {
		panic(errors.New(errors.UnauthorizedStatus, "LoginPlease", "token is expired"))
	}

	// Check that token inside database too
	token := &models.Token{
		Id:    tokenId,
		Token: tokenString,
	}
	token.InformMeToQueryProvider()
	err = token.GetMe().ExecQueryRowErr(ctx, db)
	if err != nil {
		if sqlscan.NotFound(err) {
			panic(errors.New(errors.UnauthorizedStatus, "LoginPlease", err.Error()))
		} else {
			utils.Panic500(err)
		}
	}
	if token.Token != tokenString {
		panic(errors.New(errors.UnauthorizedStatus, "LoginPlease", err.Error()))
	}

	// Now that everything is fine, get user instance
	user := models.NewUser()
	user.Id = claims.UserId
	user.GetMe().ExecQueryRow(ctx, db)

	// Set user instance and token into context
	ctx.Values().Set(g.UserKey, user)
	ctx.Values().Set(g.AccessToken, token)

	ctx.Next()
}
