package custom_validators

import (
	"context"
	g "service/global"
	"service/models"
	"service/pkg/errors"

	"github.com/georgysavva/scany/v2/sqlscan"
)

func PhoneIsUnique(input interface{}) bool {
	db, err := g.DB()
	if err != nil {
		panic(errors.New(errors.ServiceUnavailable, errors.Resend, "DbNotFound", err.Error(), nil))
	}
	defer db.Close()

	user := models.NewUser()
	err = user.Select(map[string]any{
		"phone_number": input.(string),
	}).QueryRowContextError(context.TODO(), db)
	if err != nil {
		if sqlscan.NotFound(err) {
			return true
		} else {
			panic(errors.New(errors.ServiceUnavailable, errors.Resend, "InternalServerError", err.Error(), nil))
		}
	}
	return false
}
