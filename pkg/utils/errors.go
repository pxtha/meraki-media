package utils

import (
	"fmt"
	"net/http"
)

// Generate error response

type ErrorCustomMessage int

const (
	Custom ErrorCustomMessage = iota
	WrongOTP
	CheckOTPLimit
)

const (
	ErrorEmailAlreadyUse ErrorCustomMessage = iota + 1
	ErrorLockTime
)

func ErrorMsg(err ErrorCustomMessage, a ...interface{}) string {
	msg := "Error | no message"
	switch err {
	case Custom:
		msg = fmt.Sprint(a...)
	case WrongOTP:
		msg = fmt.Sprint("Invalid OTP. Please try again.")
	case CheckOTPLimit:
		msg = fmt.Sprint("You’ve reached the maximum attempts. Try again later")
	}
	return msg
}

var messageError map[int]string

func LoadMessageError() {
	messageError = make(map[int]string)
	messageError[http.StatusOK] = "Successfully"
	messageError[http.StatusForbidden] = "Something when wrong, Your request has been rejected"
	messageError[http.StatusInternalServerError] = "Internal server error"
	messageError[http.StatusBadRequest] = "Something when wrong with your request"
	messageError[http.StatusUnauthorized] = "Unauthorized - Permission denied"
	messageError[http.StatusNotFound] = "Request not found - Check your input"
	messageError[http.StatusCreated] = "Created successfully"
	messageError[http.StatusGatewayTimeout] = "Gateway time out"
	messageError[http.StatusConflict] = "Your input has been conflict with another data"
	messageError[http.StatusTooManyRequests] = "Too many request"
}

func MessageError() map[int]string {
	return messageError
}
