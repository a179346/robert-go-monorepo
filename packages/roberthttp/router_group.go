package roberthttp

type RouterGroup struct {
	Router
	prefix string
}

func newGroup(prefix string, options *RouterOptions) *RouterGroup {
	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	return &RouterGroup{
		Router: New(options),
		prefix: prefix,
	}
}
