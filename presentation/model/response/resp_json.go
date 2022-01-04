package presentation

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type NewError struct {
	StatusCode int
	ErrMsg     error
}

type Response struct {
	Message    string
	StatusCode int
	NewError
	Data interface{}
}

func NewInternalServerError(err error) NewError {
	return NewError{StatusCode: http.StatusInternalServerError,
		ErrMsg: err}
}

func NewUnauthorizedError(err error) NewError {
	return NewError{StatusCode: http.StatusUnauthorized,
		ErrMsg: err}
}

func NewBadRequestError(err error) NewError {
	return NewError{StatusCode: http.StatusBadRequest,
		ErrMsg: err}
}

func NewResponse(message string, data interface{}, errInfo NewError) gin.H {
	if errInfo.StatusCode == 0 {
		errInfo.StatusCode = http.StatusOK
	}
	return gin.H{
		"Message": message, "StatusCode": errInfo.StatusCode, "Data": data,
		"NewError": errInfo,
	}
	// return Response{Message: message, StatusCode: errInfo.StatusCode,
	// 	Data: data, NewError: errInfo}
}
