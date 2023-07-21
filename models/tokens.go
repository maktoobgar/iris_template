package models

import (
	"database/sql"
	g "service/global"
	"service/pkg/repositories"
	"time"

	"github.com/kataras/iris/v12"
)

var TokenName = "tokens"

type Token struct {
	repositories.QueryGenerator `json:"-"`

	Id             int64     `json:"id" db:"id" skipInsert:"+"`
	Token          string    `json:"token" db:"token" skipUpdate:"+"`
	IsRefreshToken bool      `json:"is_refresh_token" db:"is_refresh_token" skipUpdate:"+"`
	UserId         *int64    `json:"-" db:"user_id" skipUpdate:"+" nilOnEmpty:"+"`
	User           *User     `json:"-"`
	ExpiresAt      time.Time `json:"expires_at" db:"expires_at" skipUpdate:"+"`
	CreatedAt      time.Time `json:"created_at" db:"created_at" skipUpdate:"+"`
}

func (t *Token) GetUser(ctx iris.Context, db *sql.DB) *User {
	if t.User == nil {
		user := NewUser()
		user.Id = *t.UserId
		t.User = user
	}

	t.User.GetMe().ExecQueryRow(ctx, db)
	return t.User
}

func (t *Token) InformMeToQueryProvider() *Token {
	t.QueryGenerator = repositories.NewQueryGenerator(TokenName)
	t.SetRowData(t)
	t.SetDbType(g.MainDatabaseType)
	return t
}

func NewToken(accessRefreshToken string, isRefreshToken bool, expiresAt time.Time, userId int64) *Token {
	user := NewUser()
	user.Id = userId
	token := &Token{
		QueryGenerator: repositories.NewQueryGenerator(TokenName),

		Token:          accessRefreshToken,
		IsRefreshToken: isRefreshToken,
		UserId:         &userId,
		User:           user,
		ExpiresAt:      expiresAt,
		CreatedAt:      time.Now(),
	}
	token.SetRowData(token)
	return token
}
