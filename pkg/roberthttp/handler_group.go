package roberthttp

type handlerFuncWithPrefix struct {
	f      HandlerFunc
	prefix string
}

func newHandlerFuncWithPrefix(f HandlerFunc, prefix string) handlerFuncWithPrefix {
	return handlerFuncWithPrefix{
		prefix: prefix,
		f:      f,
	}
}

type handlerGroup struct {
	funcWithPrefixs []handlerFuncWithPrefix
	pattern         string
	all             bool
}

func newHandlerGroup(pattern string, all bool) handlerGroup {
	return handlerGroup{
		pattern: pattern,
		all:     all,
	}
}

func (hg *handlerGroup) addHandlerFuncWithPrefix(funcWithPrefix handlerFuncWithPrefix) {
	hg.funcWithPrefixs = append(hg.funcWithPrefixs, funcWithPrefix)
}

func (hg *handlerGroup) getHandlerFuncWithPrefix(idx int) handlerFuncWithPrefix {
	return hg.funcWithPrefixs[idx]
}

func (hg *handlerGroup) len() int {
	return len(hg.funcWithPrefixs)
}

func (dst *handlerGroup) copyHandlerFuncWithPrefixsFrom(src *handlerGroup) {
	for _, funcWithPrefix := range src.funcWithPrefixs {
		dst.addHandlerFuncWithPrefix(funcWithPrefix)
	}
}
