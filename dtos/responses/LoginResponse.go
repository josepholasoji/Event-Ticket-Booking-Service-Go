package responses

type LoginResponse struct {
	Response
}

func NewLoginResponse(code int, message string, token string) *LoginResponse {
	return &LoginResponse{
		Response: Response{
			Code:    code,
			Message: message,
			Data:    token,
		},
	}
}
