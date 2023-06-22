package dto

type PanicResponse struct {
	Message string `json:"message"`
	Action  int    `json:"action"`
	Code    int    `json:"code"`
	Errors  any    `json:"errors"`
}
