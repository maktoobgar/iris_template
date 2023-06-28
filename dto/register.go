package dto

type RegisterRequest struct {
	PhoneNumber string `json:"phone_number" g:"phone,required,phone_is_unique"`
	DisplayName string `json:"display_name" g:"required"`
	Password    string `json:"password" g:"required"`
}

var RegisterRequestValidator = g.Validator(RegisterRequest{})
