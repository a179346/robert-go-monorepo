package robert_router_options

import "github.com/a179346/robert-go-monorepo/pkg/roberthttp"

type JsonResponse struct {
	Data interface{} `json:"data"`
}

func ResponseJsonWrapper(_ int, data interface{}) interface{} {
	return JsonResponse{
		Data: data,
	}
}

func GetResponseOptions() *roberthttp.ResponseOptions {
	return &roberthttp.ResponseOptions{
		JsonWrapper: ResponseJsonWrapper,
	}
}
