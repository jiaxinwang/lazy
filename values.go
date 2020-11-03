package lazy

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/adam-hanna/arrayOperations"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

// BeforeLazy get before-action's param
func BeforeLazy(params map[string][]string) (eq map[string][]string, gt, lt, gte, lte map[string]string, reduced map[string][]string) {
	var buf bytes.Buffer
	reduced = make(map[string][]string)
	if err := gob.NewEncoder(&buf).Encode(params); err != nil {
		return
	}
	gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(&reduced)

	eq = make(map[string][]string)
	gt = make(map[string]string)
	lt = make(map[string]string)
	gte = make(map[string]string)
	lte = make(map[string]string)

	for k, vv := range params {
		if vv == nil || !strings.HasPrefix(k, "before_") {
			continue
		}
		for _, v := range vv {
			name := k
			destM := &eq
			destS := &gt
			switch {
			case strings.HasSuffix(k, `_gt`):
				name = strings.TrimPrefix(strings.TrimSuffix(k, `_gt`), "before_")
				destM = nil
				destS = &gt
			case strings.HasSuffix(k, `_lt`):
				name = strings.TrimPrefix(strings.TrimSuffix(k, `_lt`), "before_")
				destM = nil
				destS = &lt
			case strings.HasSuffix(k, `_gte`):
				name = strings.TrimPrefix(strings.TrimSuffix(k, `_gte`), "before_")
				destM = nil
				destS = &gte
			case strings.HasSuffix(k, `_lte`):
				name = strings.TrimPrefix(strings.TrimSuffix(k, `_lte`), "before_")
				destM = nil
				destS = &lte
			default:
				name = strings.TrimPrefix(k, "before_")
				destM = &eq
				destS = nil
			}
			if destS != nil {
				(*destS)[name] = v
				delete(reduced, k)
			}
			if destM != nil {
				if (*destM)[name] == nil {
					(*destM)[name] = make([]string, 0)
				}
				(*destM)[name] = append((*destM)[name], v)
				delete(reduced, k)
			}
		}
	}

	return
}

// Lazy ...
func Lazy(params map[string][]string) (eq map[string][]string, gt, lt, gte, lte map[string]string) {
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
			case strings.EqualFold(k, "size"):
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
	eq, gt, lt, gte, lte := Lazy(q)
	eqm = TagSlice(s, eq)
	gtm = Tag(s, gt)
	ltm = Tag(s, lt)
	gtem = Tag(s, gte)
	ltem = Tag(s, lte)
	return
}

// SelectBuilder ...
func SelectBuilder(s sq.SelectBuilder, eq map[string][]interface{}, gt, lt, gte, lte map[string]interface{}) sq.SelectBuilder {
	if eq != nil {
		for k, v := range eq {
			switch {
			case len(v) == 1:
				eqs := sq.Eq{k: v[0]}
				s = s.Where(eqs)
			case len(v) > 1:
				eqs := sq.Eq{k: v}
				s = s.Where(eqs)
			}
		}
	}
	if gt != nil && len(gt) > 0 {
		m := sq.Gt(gt)
		s = s.Where(m)
	}
	if lt != nil && len(lt) > 0 {
		m := sq.Lt(lt)
		s = s.Where(m)
	}
	if gte != nil && len(gte) > 0 {
		m := sq.GtOrEq(gte)
		s = s.Where(m)
	}
	if lte != nil && len(lte) > 0 {
		m := sq.LtOrEq(lte)
		s = s.Where(m)
	}
	return s
}

// SelectEq ...
func SelectEq(db *gorm.DB, table, columms string, eq map[string][]interface{}) (ret []map[string]interface{}, err error) {
	sel := sq.Select(columms).From(table)
	sel = SelectBuilder(sel, eq, nil, nil, nil, nil)
	return ExecSelect(db, sel)
}

// ExecSelect ...
func ExecSelect(db *gorm.DB, active sq.SelectBuilder) (ret []map[string]interface{}, err error) {
	ret = make([]map[string]interface{}, 0)
	sql, args, err := active.ToSql()
	if err != nil {
		return ret, err
	}

	rows, sqlErr := db.Raw(sql, args...).Rows()

	defer rows.Close()
	if sqlErr != nil {
		return ret, sqlErr
	}

	columns, err := rows.Columns()
	if err != nil {
		return ret, err
	}
	length := len(columns)
	for rows.Next() {
		current := makeResultReceiver(length)
		if err := rows.Scan(current...); err != nil {
			return ret, err
		}
		value := make(map[string]interface{})
		for i := 0; i < length; i++ {
			k := columns[i]
			val := *(current[i]).(*interface{})
			if val == nil {
				value[k] = nil
				continue
			}
			vType := reflect.TypeOf(val)
			switch vType.String() {
			case "uint8":
				value[k] = val.(int8)
			case "uint16":
				value[k] = val.(int16)
			case "uint32":
				value[k] = val.(int32)
			case "uint64":
				value[k] = val.(int64)
			case "int8":
				value[k] = val.(int8)
			case "int16":
				value[k] = val.(int16)
			case "int32":
				value[k] = val.(int32)
			case "int64":
				value[k] = val.(int64)
			case "bool":
				value[k] = val.(bool)
			case "string":
				value[k] = val.(string)
			case "time.Time":
				value[k] = val.(time.Time)
			case "[]uint8":
				value[k] = string(val.([]uint8))
			default:
				// logrus.Warnf("unsupport data type '%s' now\n", vType)
			}
		}
		ret = append(ret, value)
	}

	return ret, nil

}

func makeResultReceiver(length int) []interface{} {
	result := make([]interface{}, 0, length)
	for i := 0; i < length; i++ {
		var current interface{}
		current = struct{}{}
		result = append(result, &current)
	}
	return result
}

func ignoreValues(c *gin.Context) (ret map[string][]string) {
	if v, ok := c.Get(`_ignore_values`); ok {
		ret = v.(map[string][]string)
	} else {
		ret = make(map[string][]string, 0)
	}
	params := map[string][]string(c.Request.URL.Query())
	ret = mergeValues(params, ret)
	if len(ret) > 0 {
		c.Set("_ignore_values", ret)
	}
	return
}

func additionValues(c *gin.Context, add map[string][]string) (ret map[string][]string) {
	if v, ok := c.Get(`_additional_values`); ok {
		ret = v.(map[string][]string)
	} else {
		ret = make(map[string][]string, 0)
	}
	ret = mergeValues(ret, add)
	if len(ret) > 0 {
		c.Set("_additional_values", ret)
	}
	return
}

func mergeValues(a, b map[string][]string) (ret map[string][]string) {
	ret = make(map[string][]string, len(a))
	for k, v := range a {
		tmp := make([]string, len(v))
		copy(tmp, v)
		ret[k] = tmp
	}
	for k, v := range b {
		if inRetV, ok := ret[k]; ok {
			z, ok := arrayOperations.Union(inRetV, v)
			if !ok {
				fmt.Println("Cannot find difference")
			}
			ret[k] = z.Interface().([]string)
		} else {
			ret[k] = v
		}
	}
	return
}
