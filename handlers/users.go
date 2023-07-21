package handlers

import (
	"database/sql"
	"service/dto"
	g "service/global"
	"service/models"
	"service/pkg/translator"
	"service/utils"

	"github.com/kataras/iris/v12"
)

var (
	defaultUsersParams = dto.PaginationUsers{
		OrderBy: "id",
		Search:  "",
		Sort:    "asc",
		PerPage: 10,
		Page:    1,
	}
)

func Users(ctx iris.Context) {
	// Get required data from context
	user := ctx.Values().Get(g.UserKey).(*models.User)
	db := ctx.Values().Get(g.DbInstance).(*sql.DB)
	translate := ctx.Values().Get(g.TranslateKey).(translator.TranslatorFunc)

	// Initialize params and validate them
	params := &dto.PaginationUsers{
		OrderBy: ctx.URLParamDefault("order_by", defaultUsersParams.OrderBy),
		Search:  ctx.URLParamDefault("search", defaultUsersParams.Search),
		Sort:    ctx.URLParamDefault("sort", defaultUsersParams.Sort),
		PerPage: ctx.URLParamIntDefault("per_page", defaultUsersParams.PerPage),
		Page:    ctx.URLParamIntDefault("page", defaultUsersParams.Page),
	}
	utils.Validate(params, dto.PaginationUsersValidator, translate)

	// Generate where search text
	selectParams := map[string]string{
		"display_name": params.Search,
		"phone_number": params.Search,
		"email":        params.Search,
		"first_name":   params.Search,
		"last_name":    params.Search,
	}
	wheres := user.GetLikeWheres(selectParams)

	// Get count of all users and all users in that spacific page
	users := &[]*models.User{}
	usersCount := user.SelectCount().ExecQueryCount(ctx, db)
	user.SelectWhere(wheres).OrderBy(params.OrderBy, params.Sort).Paginate(params.PerPage, params.Page).ExecQueryMulti(ctx, db, users)

	// Create and send the page
	utils.SendPage(ctx, usersCount, params.PerPage, params.Page, users)
}
