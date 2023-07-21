package dto

type PaginationUsers struct {
	OrderBy string `g:"choices=id&display_name&created_at"`
	Search  string `g:""`
	Sort    string `g:"choices=asc&desc"`
	PerPage int    `g:"min=5"`
	Page    int    `g:"min=1"`
}

var PaginationUsersValidator = g.Validator(PaginationUsers{})
