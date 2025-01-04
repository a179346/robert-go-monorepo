package roberthttp

type ResponseErrorWrapper func(statusCode int, message string, data interface{}) interface{}

type DefaultResponseError struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func defaultResponseErrorWrapper(statusCode int, message string, data interface{}) interface{} {
	return DefaultResponseError{
		Message: message,
		Data:    data,
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
