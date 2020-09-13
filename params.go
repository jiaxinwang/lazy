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
type Params map[string][]string

func params(c *gin.Context) (*Params, error) {
	paramsItr, ok := c.Get(KeyParams)
	if !ok {
		return nil, ErrParamMissing
	}
	params, ok := paramsItr.(Params)
	if !ok {
		return nil, ErrUnknown
	}
	return &params, nil
}

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

func separatePage(params Params) (filterParams Params, page, limit, offset uint64) {
	var s Params
	s, filterParams = separateParams(params, "offset", "page", "limit")
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
