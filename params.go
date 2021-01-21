package lazy

import (
	"bytes"
	"encoding/gob"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Params maps a string key to a list of values.
// type Params map[string][]string

func params(c *gin.Context) (map[string][]string, error) {
	paramsItr, ok := c.Get(KeyParams)
	if !ok {
		return nil, ErrParamMissing
	}
	params, ok := paramsItr.(map[string][]string)
	if !ok {
		return nil, ErrUnknown
	}
	return params, nil
}

func valueSliceWithParamKey(params map[string][]string, key string) []string {
	if v, exist := params[key]; exist {
		return v
	}
	return nil
}

func valueOfSingleParam(params map[string][]string, key string) (value string) {
	if v, exist := params[key]; exist {
		if len(v) == 1 {
			return v[0]
		}
	}
	return ``
}

func separatePage(params map[string][]string) (filterParams map[string][]string, page, limit, offset uint64) {
	var s map[string][]string
	s, filterParams = separateParams(params, "offset", "page", "limit")
	str := valueOfSingleParam(s, `offset`)
	offset, _ = strconv.ParseUint(str, 10, 64)
	str = valueOfSingleParam(s, `page`)
	page, _ = strconv.ParseUint(str, 10, 64)
	str = valueOfSingleParam(s, `limit`)
	limit, _ = strconv.ParseUint(str, 10, 64)
	return
}

func separateParams(maps map[string][]string, keys ...string) (separated, remain map[string][]string) {
	separated = make(map[string][]string)
	remain = make(map[string][]string)
	for k, v := range maps {
		remain[k] = v
	}
	for _, v := range keys {
		if vItem, ok := maps[v]; ok {
			separated[v] = vItem
		}
	}
	for k := range separated {
		if _, ok := remain[k]; ok {
			delete(remain, k)
		}
	}
	return
}

func separatePrefixParams(whole map[string][]string, prefix string) (separated, remain map[string][]string) {
	keys := make([]string, 0)
	for k := range whole {
		if strings.HasPrefix(k, prefix) {
			keys = append(keys, k)
		}
	}
	return separateParams(whole, keys...)
}

// ContentParams return params in content
func ContentParams(c *gin.Context) (union, query, body map[string]interface{}) {
	if v, ok := c.Get(KeyParamsUnion); ok {
		union = v.(map[string]interface{})
	} else {
		union = make(map[string]interface{})
	}
	if v, ok := c.Get(KeyParams); ok {
		p := v.(map[string][]string)
		query = make(map[string]interface{})
		for kk, vv := range p {
			query[kk] = vv
		}
	} else {
		query = make(map[string]interface{})
	}
	if v, ok := c.Get(KeyBody); ok {
		body = v.(map[string]interface{})
	} else {
		body = make(map[string]interface{})
	}
	return
}

// ModifyContentParamWithJSONPath ...
func ModifyContentParamWithJSONPath(c *gin.Context, jsonPath string, value interface{}) {
	_, _, body := ContentParams(c)
	var srcStr string
	var err error
	if srcStr, err = json.MarshalToString(body); err != nil {
		return
	}
	if srcStr, err = sjson.Set(srcStr, jsonPath, value); err != nil {
		return
	}
	m := make(map[string]interface{})
	if err := json.UnmarshalFromString(srcStr, &m); err != nil {
		return
	}
	c.Set(KeyParamsUnion, m)
}

// ContentParamWithJSONPath ...
func ContentParamWithJSONPath(c *gin.Context, jsonPath string) (param interface{}) {
	_, _, body := ContentParams(c)
	var srcStr string
	var err error
	if srcStr, err = json.MarshalToString(body); err != nil {
		return
	}
	ret := gjson.Get(srcStr, jsonPath)
	return ret.Value()
}

// SetJSON ...
func SetJSON(src map[string]interface{}, maps []JSONPathMap) (dest map[string]interface{}) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	if err := enc.Encode(src); err != nil {
		logrus.WithError(err).Trace()
		return
	}
	if err := dec.Decode(&dest); err != nil {
		logrus.WithError(err).Trace()
		return
	}
	var err error
	for _, v := range maps {
		src, dest, err = SetJSONSingle(src, dest, v)
		if err != nil {
			return
		}
	}

	return
}

// SetJSONSingle ...
func SetJSONSingle(src, dest map[string]interface{}, m JSONPathMap) (convertSrc, convertDesc map[string]interface{}, err error) {
	srcStr, err := json.MarshalToString(src)
	if err != nil {
		return src, dest, err
	}
	destStr, err := json.MarshalToString(dest)
	if err != nil {
		return src, dest, err
	}
	value := gjson.Get(srcStr, m.Src).Value()
	if value == nil {
		value = m.Default
	}
	destStr, err = sjson.Set(destStr, m.Dest, value)
	if err != nil {
		return src, dest, err
	}
	if m.Remove {
		if srcStr, err = sjson.Delete(srcStr, m.Src); err != nil {
			return src, dest, err
		}
	}
	if err = json.UnmarshalFromString(srcStr, &convertSrc); err != nil {
		return src, dest, err
	}
	if err = json.UnmarshalFromString(destStr, &convertDesc); err != nil {
		return src, dest, err
	}
	return
}
