package api

import (
	"CRUD-API/models"
)

func ResponseApi(status int, data interface{}, err error) (res *models.ResponseApi) {
	if err != nil {
		res = &models.ResponseApi{
			Status:       status,
			Description:  "FAILED",
			Data:         nil,
			ErrorMessage: err.Error(),
		}
		return res
	}
	res = &models.ResponseApi{
		Status:       status,
		Description:  "SUCCESS",
		Data:         data,
		ErrorMessage: "",
	}
	return res
}
