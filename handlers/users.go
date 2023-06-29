package handlers

import (
	"database/sql"
	"fmt"
	"service/dto"
	g "service/global"
	"service/models"
	"service/pkg/errors"
	"service/utils"

	"github.com/kataras/iris/v12"
)

var (
	defaultUsersParams = dto.UsersParams{
		OrderBy: "id",
		Sort:    "asc",
		PerPage: 10,
		Page:    1,
	}
)

func Users(ctx iris.Context) {
	db := ctx.Values().Get(g.DbInstance).(*sql.DB)
	userInternal := models.NewUserInternal()

	orderBy := ctx.URLParamDefault("order_by", defaultUsersParams.OrderBy)
	sort := ctx.URLParamDefault("sort", defaultUsersParams.Sort)
	perPage := ctx.URLParamIntDefault("per_page", defaultUsersParams.PerPage)
	page := ctx.URLParamIntDefault("page", defaultUsersParams.Page)
	queryParams := &dto.UsersParams{
		OrderBy: orderBy,
		Sort:    sort,
		PerPage: perPage,
		Page:    page,
	}
	errs := dto.UsersParamsValidator.Validate(queryParams)
	if errs != nil {
		panic(errors.New(errors.InvalidStatus, errors.DoNothing, "InvalidPageParameters", "query parameters validation failed", errs))
	}

	usersCount := userInternal.SelectCount(nil).QueryCountContext(ctx, db)
	pagesCount := utils.CalculatePagesCount(usersCount, perPage)
	if page > pagesCount {
		panic(errors.New(errors.NotFoundStatus, errors.DoNothing, "PageNotFound", fmt.Sprintf("page %d requested but we have %d pages", page, pagesCount)))
	}

	desc := false
	if queryParams.Sort == "desc" {
		desc = true
	}

	users := &[]*models.UserInternal{}
	userInternal.Select(nil).OrderBy(orderBy, desc).Paginate(perPage, page).QueryContext(ctx, db, users)

	utils.SendPage(ctx, usersCount, perPage, page, users)
}
