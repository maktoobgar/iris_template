package models

import (
	"database/sql"
	"fmt"
	g "service/global"
	"service/pkg/repositories"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"golang.org/x/crypto/bcrypt"
)

var UserName = "users"

type UserInternal struct {
	repositories.QueryGenerator `json:"-"`

	Id          int64     `json:"id" db:"id" skipInsert:"+"`
	DisplayName string    `json:"display_name" db:"display_name"`
	CreatedAt   time.Time `json:"created_at" db:"created_at" skipUpdate:"+"`
}

func NewUserInternal() *UserInternal {
	user := &UserInternal{
		QueryGenerator: repositories.NewQueryGenerator(UserName),
		CreatedAt:      time.Now(),
	}
	user.SetRowData(user)
	user.SetDbType(g.MainDatabaseType)
	return user
}

type User struct {
	UserInternal

	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"-" db:"password"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	IsActive    bool   `json:"-" db:"is_active"`
	IsAdmin     bool   `json:"-" db:"is_admin"`
	IsSuperuser bool   `json:"-" db:"is_superuser"`
}

func (u *User) CreateAccessToken(ctx iris.Context, db *sql.DB) *Token {
	expirationTime := time.Now().Add(time.Duration(g.CFG.AccessTokenLifePeriod) * (time.Hour * 24))

	claims := &Claims{
		UserId: u.Id,
		Type:   AccessTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := tkn.SignedString(g.SecretKeyBytes)
	token := NewToken(tokenString, false, expirationTime, u.Id)
	if token.InsertInto().ExecQuery(ctx, db) == 0 {
		token.Select(map[string]any{"token": token.Token, "user_id": u.Id}).ExecQueryRow(ctx, db)
	}
	token.Token = fmt.Sprintf("%d|%s", token.Id, token.Token)
	token.User = u
	return token
}

func (u *User) CreateRefreshToken(ctx iris.Context, db *sql.DB) *Token {
	expirationTime := time.Now().Add(time.Duration(g.CFG.RefreshTokenLifePeriod) * (time.Hour * 24 * 30))

	claims := &Claims{
		UserId: u.Id,
		Type:   RefreshTokenType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := tkn.SignedString(g.SecretKeyBytes)
	token := NewToken(tokenString, true, expirationTime, u.Id)
	if token.InsertInto().ExecQuery(ctx, db) == 0 {
		token.Select(map[string]any{"token": token.Token, "user_id": u.Id}).ExecQueryRow(ctx, db)
	}
	token.Token = fmt.Sprintf("%d|%s", token.Id, token.Token)
	token.User = u
	return token
}

func (u *User) HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 16)
	return string(bytes)
}

func (u *User) HashMyPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 16)
	u.Password = string(bytes)
}

func (u *User) IsPasswordEqualToHash(password string, hashed string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)); err != nil {
		return false
	}

	return true
}

func (u *User) IsPasswordEqualToMyHash(password string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return false
	}

	return true
}

func (u *User) InformMeToQueryProvider() *User {
	u.QueryGenerator = repositories.NewQueryGenerator(UserName)
	u.SetRowData(u)
	u.SetDbType(g.MainDatabaseType)
	return u
}

func NewUser() *User {
	user := &User{
		UserInternal: UserInternal{
			QueryGenerator: repositories.NewQueryGenerator(UserName),
			CreatedAt:      time.Now(),
		},
	}
	user.SetRowData(user)
	user.SetDbType(g.MainDatabaseType)
	return user
}
