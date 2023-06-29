package dto

type UsersParams struct {
	OrderBy string `g:"choices=id&display_name&created_at"`
	Sort    string `g:"choices=asc&desc"`
	PerPage int    `g:"min=5"`
	Page    int    `g:"min=1"`
}

var UsersParamsValidator = g.Validator(UsersParams{})
