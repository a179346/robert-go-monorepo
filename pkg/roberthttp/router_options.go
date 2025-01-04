package roberthttp

type RouterOption func(*RouterOptions)

type RouterOptions struct {
	Response *ResponseOptions
}

func defaultRouterOptions(opt *RouterOptions) *RouterOptions {
	if opt == nil {
		opt = &RouterOptions{}
	}
	opt.Response = defaultResponseOptions(opt.Response)
	return opt
}
