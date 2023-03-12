package model

type ErrorResponseModel struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *ErrorResponseModel) Msg() string {
	return e.Message
}

func (e *ErrorResponseModel) Error() string {
	return e.Code
}
