package service

type ResponseOk struct {
	Status  string `json:"status" example:"ok"`
	Message string `json:"message" example:"success"`
	Body    any    `json:",omitempty"`
}

type ResponseError struct {
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"fail"`
}
