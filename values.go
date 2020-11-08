package lazy

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
)

// StructMap converts struct to map
func StructMap(src interface{}, timeLayout string) (ret map[string]interface{}, err error) {
	if b, err := json.Marshal(src); err != nil {
		return nil, err
	} else {
		ret = make(map[string]interface{})
		if err = json.Unmarshal(b, &ret); err != nil {
			return nil, err
		}

		switch v := reflect.ValueOf(src); v.Kind() {
		case reflect.Struct:
			vofs := reflect.ValueOf(src)
			for i := 0; i < vofs.NumField(); i++ {
				switch vofs.Field(i).Interface().(type) {
				case *time.Time:
					t := vofs.Field(i).Interface().(*time.Time)
					name, err := dbNameWithFieldName(v, vofs.Field(i).Type().Name())
					if err != nil {
						logrus.WithError(fmt.Errorf("can't find schema field name: %s", name)).Warn()
						continue
					}

					if _, ok := ret[name]; ok {
						if t != nil {
							ret[name] = t.Format(timeLayout)
						}
					}
				case time.Time:
					t := vofs.Field(i).Interface().(time.Time)
					name, err := dbNameWithFieldName(v, vofs.Field(i).Type().Name())
					if err != nil {
						logrus.WithError(fmt.Errorf("can't find schema field name: %s", name)).Warn()
						continue
					}
					if _, ok := ret[name]; ok {
						ret[name] = t.Format(timeLayout)
					}
				}
			}
		default:
		}

		return ret, nil
	}
}

func toTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if t != reflect.TypeOf(time.Time{}) {
			return data, nil
		}

		switch f.Kind() {
		case reflect.String:
			return time.Parse(time.RFC3339, data.(string))
		case reflect.Float64:
			return time.Unix(0, int64(data.(float64))*int64(time.Millisecond)), nil
		case reflect.Int64:
			return time.Unix(0, data.(int64)*int64(time.Millisecond)), nil
		default:
			return data, nil
		}
	}
}

// Map2StructWithJSON ...
func Map2StructWithJSON(input map[string]interface{}, result interface{}) (err error) {
	var s string
	if s, err = json.MarshalToString(input); err != nil {
		return fmt.Errorf("can't json marshal: %w", err)
	}
	if err = json.UnmarshalFromString(s, result); err != nil {
		return fmt.Errorf("can't json unmarshal: %w", err)
	}
	return
}

// MapSlice2StructSliceWithJSON ...
func MapSlice2StructSliceWithJSON(input []map[string]interface{}, result *[]interface{}) (err error) {
	var s string
	if s, err = json.MarshalToString(input); err != nil {
		return fmt.Errorf("can't json marshal: %w", err)
	}
	if err = json.UnmarshalFromString(s, result); err != nil {
		return fmt.Errorf("can't json unmarshal: %w", err)
	}
	return
}

// MapStruct converts map to struct
func MapStruct(input map[string]interface{}, result interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata:         nil,
		WeaklyTypedInput: true,
		DecodeHook: mapstructure.ComposeDecodeHookFunc(
			toTimeHookFunc()),
		Result: result,
	})
	if err != nil {
		return err
	}

	if err := decoder.Decode(input); err != nil {
		return err
	}
	return err
}

// Parse ...
func Parse(v string, k reflect.Kind) (ret interface{}) {
	switch k {
	case reflect.Uint:
		if kv, err := strconv.ParseUint(v, 10, 64); err == nil {
			ret = uint(kv)
		}
	case reflect.Uint64:
		if kv, err := strconv.ParseUint(v, 10, 64); err == nil {
			ret = uint64(kv)
		}
	case reflect.Uint32:
		if kv, err := strconv.ParseUint(v, 10, 32); err == nil {
			ret = uint32(kv)
		}
	case reflect.Uint16:
		if kv, err := strconv.ParseUint(v, 10, 16); err == nil {
			ret = uint16(kv)
		}
	case reflect.Uint8:
		if kv, err := strconv.ParseUint(v, 10, 8); err == nil {
			ret = uint8(kv)
		}
	case reflect.Int:
		if kv, err := strconv.ParseInt(v, 10, 64); err == nil {
			ret = int(kv)
		}
	case reflect.Int64:
		if kv, err := strconv.ParseInt(v, 10, 64); err == nil {
			ret = int64(kv)
		}
	case reflect.Int32:
		if kv, err := strconv.ParseInt(v, 10, 32); err == nil {
			ret = int32(kv)
		}
	case reflect.Int16:
		if kv, err := strconv.ParseInt(v, 10, 16); err == nil {
			ret = int16(kv)
		}
	case reflect.Int8:
		if kv, err := strconv.ParseInt(v, 10, 8); err == nil {
			ret = int8(kv)
		}
	case reflect.String:
		ret = v
	case reflect.Bool:
		if kv, err := strconv.ParseBool(v); err == nil {
			ret = kv
		}
	default:
		fmt.Print("unsupported kind")
	}
	return
}

func dbNameWithFieldName(v interface{}, fieldName string) (string, error) {
	m, err := schema.Parse(v, schemaStore, schema.NamingStrategy{})
	if err != nil {
		return fieldName, err
	}
	schemaField, ok := m.FieldsByName[fieldName]
	if !ok {
		return fieldName, fmt.Errorf("can't find schema field name: %s", fieldName)
	}
	return schemaField.DBName, nil

}

// TagSlice ...
func TagSlice(v interface{}, params map[string][]string) map[string][]interface{} {
	ret := make(map[string][]interface{})
	val := reflect.ValueOf(v).Elem()

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		dbName, err := dbNameWithFieldName(v, field.Name)
		if err != nil {
			logrus.WithError(fmt.Errorf("can't find schema field name: %s", field.Name)).Warn()
			continue
		}
		if t := dbName; t != `` {
			if vv, ok := params[t]; ok {
				ret[t] = make([]interface{}, 0)
				for _, vvv := range vv {
					ret[t] = append(ret[t], Parse(vvv, field.Type.Kind()))
				}
			}
		}
	}

	return ret
}

// Tag ...
func Tag(v interface{}, m map[string]string) map[string]interface{} {
	ret := make(map[string]interface{})
	val := reflect.ValueOf(v).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		name, err := dbNameWithFieldName(v, field.Name)
		if err != nil {
			logrus.WithError(fmt.Errorf("can't find schema field name: %s", field.Name)).Warn()
			continue
		}

		if t := name; t != `` {
			if vv, ok := m[t]; ok {
				name := field.Name
				r := reflect.ValueOf(v)
				f := reflect.Indirect(r).FieldByName(name)
				fieldValue := f.Interface()
				switch vvv := fieldValue.(type) {
				case uint64:
					i, _ := strconv.ParseUint(vv, 10, 64)
					ret[t] = i
				case uint32:
					i, _ := strconv.ParseUint(vv, 10, 32)
					ret[t] = i
				case uint:
					i, _ := strconv.ParseUint(vv, 10, 64)
					ret[t] = int(i)
				case int64:
					i, _ := strconv.ParseInt(vv, 10, 64)
					ret[t] = i
				case int32:
					i, _ := strconv.ParseInt(vv, 10, 32)
					ret[t] = i
				case int:
					i, _ := strconv.ParseInt(vv, 10, 64)
					ret[t] = int(i)
				case string:
					ret[t] = vv
				case bool:
					ret[t], _ = strconv.ParseBool(vv)
				case time.Time:
					ret[t], _ = time.Parse(time.RFC3339, vv)
				default:
					_ = vvv
				}
			}
		}
	}
	return ret
}

// QueryParam ...
type QueryParam struct {
	Model     interface{}
	Eq        map[string][]interface{}
	Lt        map[string][]interface{}
	Gt        map[string][]interface{}
	Lte       map[string][]interface{}
	Gte       map[string][]interface{}
	Like      map[string][]interface{}
	HasMany   map[string][]interface{}
	Many2Many map[string][]interface{}
	Page      int
	Limit     int
	Offset    int
}

func valueOfMap(params map[string][]string, key string) (value []string, ok bool) {
	if value, ok = params[key]; !ok {
		return []string{}, ok
	}

	if len(value) == 0 {
		return []string{}, false
	}
	return
}

func toGenericArray(arr ...interface{}) []interface{} {
	return arr
}

func splitParams1(model interface{}, params map[string][]string) (queryParam QueryParam, err error) {
	queryParam.Eq = make(map[string][]interface{})
	queryParam.Gt = make(map[string][]interface{})
	queryParam.Lt = make(map[string][]interface{})
	queryParam.Gte = make(map[string][]interface{})
	queryParam.Lte = make(map[string][]interface{})
	queryParam.Like = make(map[string][]interface{})
	queryParam.HasMany = make(map[string][]interface{})
	queryParam.Many2Many = make(map[string][]interface{})

	if v, ok := valueOfMap(params, "offset"); ok {
		if offset, err := strconv.Atoi(v[0]); err == nil {
			queryParam.Offset = offset
		} else {
			queryParam.Offset = 0
		}
	}

	if v, ok := valueOfMap(params, "limit"); ok {
		if limit, err := strconv.Atoi(v[0]); err == nil {
			queryParam.Limit = limit
		} else {
			queryParam.Limit = 1000
		}
	}

	if v, ok := valueOfMap(params, "page"); ok {
		if page, err := strconv.Atoi(v[0]); err == nil {
			queryParam.Page = page
		} else {
			queryParam.Page = 0
		}
	}

	m, err := schema.Parse(model, schemaStore, schema.NamingStrategy{})
	if err != nil {
		logrus.WithError(err).Error()
		return queryParam, fmt.Errorf("can't get schema : %w", err)
	}

	for _, vField := range m.Fields {
		jsonKey := vField.StructField.Tag.Get("json")
		if v, ok := vField.Schema.Relationships.Relations[vField.Name]; ok {
			logrus.Print("R: ", vField.Name)
			switch v.Type {
			case schema.HasMany:
				key := fmt.Sprintf("%s", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.HasMany[jsonKey] = toGenericArray(v)
				}
			case schema.Many2Many:
				key := fmt.Sprintf("%s", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Many2Many[jsonKey] = toGenericArray(v)
				}
			}

		} else {
			switch vField.FieldType.Kind() {
			case reflect.String:
				key := fmt.Sprintf("%s", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Eq[jsonKey] = toGenericArray(v)
				}
				key = fmt.Sprintf("%s_like", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Like[jsonKey] = toGenericArray(v)
				}
			case reflect.Bool:
				key := fmt.Sprintf("%s", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Eq[jsonKey] = toGenericArray(v)
				}
			case reflect.Struct:
				t := time.Now()
				switch vField.FieldType {
				case reflect.TypeOf(t), reflect.TypeOf(&t):
					// TODO: time format
					key := fmt.Sprintf("%s", jsonKey)
					if v, ok := valueOfMap(params, key); ok {
						queryParam.Eq[jsonKey] = toGenericArray(v)
					}
					key = fmt.Sprintf("%s_lt", jsonKey)
					if v, ok := valueOfMap(params, key); ok {
						queryParam.Lt[jsonKey] = toGenericArray(v)
					}
					key = fmt.Sprintf("%s_gt", jsonKey)
					if v, ok := valueOfMap(params, key); ok {
						queryParam.Gt[jsonKey] = toGenericArray(v)
					}
					key = fmt.Sprintf("%s_lte", jsonKey)
					if v, ok := valueOfMap(params, key); ok {
						queryParam.Lte[jsonKey] = toGenericArray(v)
					}
					key = fmt.Sprintf("%s_gte", jsonKey)
					if v, ok := valueOfMap(params, key); ok {
						queryParam.Gte[jsonKey] = toGenericArray(v)
					}
				}
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				fallthrough
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				fallthrough
			case reflect.Float32, reflect.Float64:
				key := fmt.Sprintf("%s", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Eq[jsonKey] = toGenericArray(v)
				}
				key = fmt.Sprintf("%s_lt", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Lt[jsonKey] = toGenericArray(v)
				}
				key = fmt.Sprintf("%s_gt", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Gt[jsonKey] = toGenericArray(v)
				}
				key = fmt.Sprintf("%s_lte", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Lte[jsonKey] = toGenericArray(v)
				}
				key = fmt.Sprintf("%s_gte", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Gte[jsonKey] = toGenericArray(v)
				}

			}
			logrus.WithFields(
				logrus.Fields{
					"DBName": vField.DBName,
					"Name":   vField.Name,
					"Kind":   vField.FieldType.Kind(),
				},
			).Info()

		}
	}

	return
}

// splitParams ...
func splitParams(params map[string][]string) (eq map[string][]string, gt, lt, gte, lte map[string]string) {
	eq = make(map[string][]string)
	gt = make(map[string]string)
	lt = make(map[string]string)
	gte = make(map[string]string)
	lte = make(map[string]string)

	for k, vv := range params {
		if vv == nil {
			continue
		}
		for _, v := range vv {
			name := k
			destM := &eq
			destS := &gt
			switch {
			case strings.EqualFold(k, "limit"):
				fallthrough
			case strings.EqualFold(k, "offset"):
				fallthrough
			case strings.EqualFold(k, "page"):
				destM = nil
				destS = nil
			case strings.HasSuffix(k, `_gt`):
				name = strings.TrimSuffix(k, `_gt`)
				destM = nil
				destS = &gt
			case strings.HasSuffix(k, `_lt`):
				name = strings.TrimSuffix(k, `_lt`)
				destM = nil
				destS = &lt
			case strings.HasSuffix(k, `_gte`):
				name = strings.TrimSuffix(k, `_gte`)
				destM = nil
				destS = &gte
			case strings.HasSuffix(k, `_lte`):
				name = strings.TrimSuffix(k, `_lte`)
				destM = nil
				destS = &lte
			default:
				destM = &eq
				destS = nil
			}
			if destS != nil {
				(*destS)[name] = v
			}
			if destM != nil {
				if (*destM)[name] == nil {
					(*destM)[name] = make([]string, 0)
				}
				(*destM)[name] = append((*destM)[name], v)
			}
		}
	}
	return
}

// URLValues ...
func URLValues(s interface{}, q map[string][]string) (eqm map[string][]interface{}, gtm, ltm, gtem, ltem map[string]interface{}) {
	eq, gt, lt, gte, lte := splitParams(q)
	eqm = TagSlice(s, eq)
	gtm = Tag(s, gt)
	ltm = Tag(s, lt)
	gtem = Tag(s, gte)
	ltem = Tag(s, lte)
	return
}
