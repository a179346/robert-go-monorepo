package roberthttp_extended

import (
	"github.com/a179346/robert-go-monorepo/pkg/roberthttp"
)

func ResponseJsonWrapper(_ int, data interface{}) interface{} {
	return JsonResponse[interface{}]{
		Data: data,
	}
}

func GetResponseOptions() *roberthttp.ResponseOptions {
	return &roberthttp.ResponseOptions{
		JsonWrapper: ResponseJsonWrapper,
	}
}
