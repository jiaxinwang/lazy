package lazy

import (
	"strconv"
	"strings"
)

// Params maps a string key to a list of values.
type Params map[string][]string

func mergeParams(a, b Params) (ret Params) {
	return
}

func valueOfParams(params Params, key string) (value string) {
	if v, exist := params[key]; exist {
		if len(v) == 1 {
			return v[0]
		}
	}
	return ``
}

func separatePage(params Params) (remain Params, page, limit, offset uint64) {
	var s Params
	s, remain = separateParams(params, "offset", "page", "limit")
	str := valueOfParams(s, `offset`)
	offset, _ = strconv.ParseUint(str, 10, 64)
	str = valueOfParams(s, `page`)
	page, _ = strconv.ParseUint(str, 10, 64)
	str = valueOfParams(s, `limit`)
	limit, _ = strconv.ParseUint(str, 10, 64)
	return
}

func separateParams(whole Params, keys ...string) (separated, remain Params) {
	separated = make(Params)
	remain = make(Params)
	for k, v := range whole {
		remain[k] = v
	}
	for _, v := range keys {
		if vInWhole, ok := whole[v]; ok {
			separated[v] = vInWhole
		}
	}
	for k := range separated {
		if _, ok := remain[k]; ok {
			delete(remain, k)
		}
	}
	return
}

func separatePrefixParams(whole Params, prefix string) (separated, remain Params) {
	keys := make([]string, 0)
	for k := range whole {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return separateParams(whole, keys...)
}
