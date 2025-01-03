package roberthttp

type RouterGroup struct {
	Router
	prefix string
}

func NewGroup(prefix string) RouterGroup {
	if prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}

	return RouterGroup{
		Router: New(),
		prefix: prefix,
	}
}
