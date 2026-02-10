package responses

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
	Data    string `json:"data,omitempty"`
}
