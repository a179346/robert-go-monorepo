package roberthttp

import "fmt"

type ResponseErrorWrapper func(statusCode int, message string, info interface{}) interface{}

type DefaultResponseError[T interface{}] struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Info    T      `json:"info"`
}

func (e DefaultResponseError[T]) Error() string {
	return fmt.Sprintf("HTTP Response Error (%d): %s", e.Status, e.Message)
}

func defaultResponseErrorWrapper(statusCode int, message string, info interface{}) interface{} {
	return DefaultResponseError[interface{}]{
		Status:  statusCode,
		Message: message,
		Info:    info,
	}
}

type ResponseJsonWrapper func(statusCode int, data interface{}) interface{}

func defaultResponseJsonWrapper(statusCode int, data interface{}) interface{} {
	return data
}

type ResponseOptions struct {
	ErrorWrapper ResponseErrorWrapper
	JsonWrapper  ResponseJsonWrapper
}

func defaultResponseOptions(opt *ResponseOptions) *ResponseOptions {
	if opt == nil {
		opt = &ResponseOptions{}
	}

	if opt.ErrorWrapper == nil {
		opt.ErrorWrapper = defaultResponseErrorWrapper
	}
	if opt.JsonWrapper == nil {
		opt.JsonWrapper = defaultResponseJsonWrapper
	}
	return opt
}
