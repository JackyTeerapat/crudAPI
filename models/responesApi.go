package models

type ResponseApi struct {
	Status       int         `json:"status"`
	Description  string      `json:"description"`
	Data         interface{} `json:"data"`
	ErrorMessage string      `json:"errorMessage"`
}
