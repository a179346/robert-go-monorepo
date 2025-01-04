package robert_router_options

import "github.com/a179346/robert-go-monorepo/packages/roberthttp"

func GetRouterOptions() *roberthttp.RouterOptions {
	return &roberthttp.RouterOptions{
		Response: GetResponseOptions(),
	}
}
