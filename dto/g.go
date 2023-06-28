package dto

import (
	"service/dto/custom_validators"

	"github.com/golodash/galidator"
)

var g = galidator.G().CustomValidators(galidator.Validators{
	"phone_is_unique": custom_validators.PhoneIsUnique,
	"email_is_unique": custom_validators.EmailIsUnique,
}).CustomMessages(galidator.Messages{
	"phone_is_unique": "PhoneIsUnique",
	"email_is_unique": "EmailIsUnique",
})
