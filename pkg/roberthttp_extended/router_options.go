package roberthttp_extended

import "github.com/a179346/robert-go-monorepo/pkg/roberthttp"

func GetRouterOptions() *roberthttp.RouterOptions {
	return &roberthttp.RouterOptions{
		Response: GetResponseOptions(),
	}
}
