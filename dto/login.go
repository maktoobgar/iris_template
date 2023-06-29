package dto

type LoginRequest struct {
	PhoneNumber string `json:"phone_number" g:"phone,required"`
	Password    string `json:"password" g:"required"`
}

var LoginRequestValidator = g.Validator(LoginRequest{})
