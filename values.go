package lazy

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/schema"
)

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

// HasManyQueryParam ...
type HasManyQueryParam struct {
	Name    string
	DotName string
	PTable  string
	Table   string
	Values  []interface{}
}

// Many2ManyQueryParam ...
type Many2ManyQueryParam struct {
	Name                      string
	JoinTable                 string
	JoinTableFieldName        string
	JoinTableForeignFieldName string
	Values                    []interface{}
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
	Ignore    map[string][]interface{}
	Order     map[string][]interface{}
	HasMany   map[string]HasManyQueryParam
	Many2Many map[string]Many2ManyQueryParam
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

func splitQueryParams(model interface{}, params map[string][]string) (queryParam QueryParam, err error) {
	queryParam.Eq = make(map[string][]interface{})
	queryParam.Gt = make(map[string][]interface{})
	queryParam.Lt = make(map[string][]interface{})
	queryParam.Gte = make(map[string][]interface{})
	queryParam.Lte = make(map[string][]interface{})
	queryParam.Like = make(map[string][]interface{})
	queryParam.Ignore = make(map[string][]interface{})
	queryParam.Order = make(map[string][]interface{})
	queryParam.HasMany = make(map[string]HasManyQueryParam)
	queryParam.Many2Many = make(map[string]Many2ManyQueryParam)

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
			switch v.Type {
			case schema.HasMany:
				key := fmt.Sprintf("%s", jsonKey)
				if vOfMap, ok := valueOfMap(params, key); ok {
					queryParam.HasMany[jsonKey] = HasManyQueryParam{
						Name:    fmt.Sprintf("%s__%s", vField.StructField.Name, "id"),
						DotName: fmt.Sprintf("%s.%s", vField.Name, "id"),
						Table:   vField.StructField.Name,
						PTable:  fmt.Sprintf("%s.%s", v.Schema.Table, "id"),
						Values:  toGenericArray(vOfMap),
					}
				}
			case schema.Many2Many:
				key := fmt.Sprintf("%s", jsonKey)
				if vOfMap, ok := valueOfMap(params, key); ok {
					queryParam.Many2Many[jsonKey] = Many2ManyQueryParam{
						JoinTable:                 v.JoinTable.Name,
						JoinTableFieldName:        v.JoinTable.Fields[0].DBName,
						JoinTableForeignFieldName: v.JoinTable.Fields[1].DBName,
						Values:                    toGenericArray(vOfMap),
					}
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
				key = fmt.Sprintf("%s_ignore", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Ignore[jsonKey] = toGenericArray(v)
				}
				key = fmt.Sprintf("%s_order", jsonKey)
				if _, ok := valueOfMap(params, key); ok {
					queryParam.Order[jsonKey] = []interface{}{}
				}
			case reflect.Bool:
				key := fmt.Sprintf("%s", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					switch v[0] {
					case "1", "true", "t":
						queryParam.Eq[jsonKey] = []interface{}{true}
					case "0", "false", "f":
						queryParam.Eq[jsonKey] = []interface{}{false}
					}
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
				key = fmt.Sprintf("%s_order", jsonKey)
				if v, ok := valueOfMap(params, key); ok {
					queryParam.Order[jsonKey] = []interface{}{}
				}

			}
		}
	}

	return
}
