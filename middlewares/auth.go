package middlewares

import (
	"database/sql"
	"fmt"
	"regexp"
	g "service/global"
	"service/models"
	"service/pkg/errors"
	"service/utils"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/georgysavva/scany/v2/sqlscan"
	"github.com/kataras/iris/v12"
)

var tokenPattern, _ = regexp.Compile(`^\d+\|.*$`)

func Auth(ctx iris.Context) {
	db := ctx.Values().Get(g.DbInstance).(*sql.DB)

	// Check if a token has sent
	tokenString := ctx.GetHeader(g.AccessToken)
	if !tokenPattern.MatchString(tokenString) {
		panic(errors.New(errors.UnauthorizedStatus, errors.ReSignIn, "LoginPlease", "sent token is not valid"))
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
		fmt.Println()
		panic(errors.New(errors.UnauthorizedStatus, errors.ReSignIn, "LoginPlease", err.Error()))
	}
	if !tkn.Valid {
		panic(errors.New(errors.UnauthorizedStatus, errors.ReSignIn, "LoginPlease", "token is invalid"))
	}
	if claims.Type != models.AccessTokenType {
		panic(errors.New(errors.UnauthorizedStatus, errors.ReSignIn, "LoginPlease", "token is not access token"))
	}

	// Check that token inside database too
	token := models.Token{
		Id: tokenId,
	}
	token.InformMeToQueryProvider()
	err = token.GetMe().QueryRowContextError(ctx, db)
	if err != nil {
		if sqlscan.NotFound(err) {
			panic(errors.New(errors.UnauthorizedStatus, errors.ReSignIn, "LoginPlease", err.Error()))
		} else {
			utils.Panic500(err)
		}
	}

	// Now that everything is fine, get user instance
	user := models.NewUser()
	user.Id = claims.UserId
	user.GetMe().QueryRowContext(ctx, db)
	ctx.Values().Set(g.UserKey, user)

	ctx.Next()
}
