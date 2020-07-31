package lazy

import (
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm/schema"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"

	"github.com/tidwall/sjson"
)

// DeleteHandle executes delete.
func DeleteHandle(c *gin.Context) (data []map[string]interface{}, err error) {
	id := c.Param("id")
	if err = validator.New().Var(id, "required,number"); err != nil {
		return nil, err
	}
	var config *Configuration
	if v, ok := c.Get(keyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrNoConfiguration
	}

	if !config.IgnoreAssociations {
		sfs := config.DB.NewScope(config.Model).GetStructFields()
		for _, v := range sfs {
			if v.Relationship != nil {
				r := v.Relationship
				switch r.Kind {
				case string(schema.HasOne), string(schema.Many2Many):
					return nil, ErrUnknown
				case string(schema.HasMany):
					count := 0
					config.DB.Table(v.DBName).Where(fmt.Sprintf("%s = ?", r.ForeignDBNames[0]), id).Count(&count)
					if count > 0 {
						return nil, ErrHasAssociations
					}
				}
			}
		}
	}

	return nil, config.DB.Where(`id = ?`, id).Delete(config.Model).Error
}

// PostHandle executes post.
func PostHandle(c *gin.Context) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(keyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrNoConfiguration
	}

	if err = c.ShouldBindJSON(config.Model); err != nil {
		return nil, err
	}

	return nil, createModel(config.DB, config.Model)
}

// PutHandle executes put.
func PutHandle(c *gin.Context) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(keyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrNoConfiguration
	}
	if err = c.ShouldBindJSON(config.Model); err != nil {
		return nil, err
	}
	return
}

// GetHandle executes actions and returns response
func GetHandle(c *gin.Context) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(keyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrNoConfiguration
	}

	set := foreignOfModel((*config).Model)

	if config.Before != nil {
		_, _, errBefore := config.Before.Action(c, config.DB, *config, nil)
		if errBefore != nil {
			return nil, errBefore
		}
	}

	paramsItr, ok := c.Get(keyParams)
	if !ok {
		return nil, errors.New("can't find lazy params")
	}
	params := paramsItr.(Params)
	remain, page, limit, offset := separatePage(params)
	c.Set(keyParams, remain)
	if limit == 0 {
		limit = 10000
	}

	var merged map[string][]string
	additional, ok := c.Get("_additional_values")
	if ok {
		merged = mergeValues(c.Request.URL.Query(), additional.(map[string][]string))
	}

	eq, gt, lt, gte, lte := LazyURLValues(config.Model, merged)

	sel := sq.Select(config.Columms).From(config.Table).Limit(limit).Offset(limit*page + offset)
	sel = SelectBuilder(sel, eq, gt, lt, gte, lte)
	data, err = ExecSelect(config.DB, sel)
	if err != nil {
		return
	}

	for _, v := range data {
		if err := MapStruct(v, config.Model); err != nil {
			return nil, err
		}
		tmp := clone(config.Model)

		for _, v := range set {
			value := valueOfTag(tmp, v[ForeignOfModelID])
			eq := map[string][]interface{}{v[ForeignOfModelForeignID]: []interface{}{value}}
			data, err := SelectEq(config.DB, v[ForeignOfModelForeignTable], "*", eq)
			if err != nil {
				return nil, err
			}
			if len(data) == 1 {
				jbyte, _ := json.Marshal(tmp)
				assemble, _ := sjson.Set(string(jbyte), v[ForeignOfModelName], data[0])
				json.Unmarshal([]byte(assemble), tmp)
			}
		}

		// TODO: batch

		config.Results = append(config.Results, tmp)
	}

	count := int64(len(data))

	if config.NeedCount {
		sel := sq.Select(`count(1) as c`).From(config.Table)
		sel = SelectBuilder(sel, eq, gt, lt, gte, lte)
		data, err = ExecSelect(config.DB, sel)
		if err != nil {
			return
		}
		if len(data) == 1 {
			iter, _ := data[0][`c`]
			count, err = strconv.ParseInt(fmt.Sprintf("%v", iter), 10, 64)
			if err != nil {
				return
			}
		}
	}
	logrus.WithField("count", count).Info()
	c.Set(keyCount, count)
	c.Set(keyData, config.Results)
	c.Set(keyResults, map[string]interface{}{"count": count, "items": config.Results})
	return
}
