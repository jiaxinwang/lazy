package lazy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/antchfx/jsonquery"
	"github.com/levigross/grequests"
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
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}

	// TODO: associations

	// if !config.IgnoreAssociations {
	// 	sfs := config.DB.NewScope(config.Model).GetStructFields()
	// 	for _, v := range sfs {
	// 		if v.Relationship != nil {
	// 			r := v.Relationship
	// 			switch r.Kind {
	// 			case string(schema.HasOne), string(schema.Many2Many):
	// 				return nil, ErrUnknown
	// 			case string(schema.HasMany):
	// 				count := 0
	// 				config.DB.Table(v.DBName).Where(fmt.Sprintf("%s = ?", r.ForeignDBNames[0]), id).Count(&count)
	// 				if count > 0 {
	// 					return nil, ErrHasAssociations
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return nil, config.DB.Where(`id = ?`, id).Delete(config.Model).Error
}

// DefaultPostAction execute default post.
func DefaultPostAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	_, _, bodyParams := ContentParams(c)
	var config *Configuration
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}
	s, err := json.MarshalToString(bodyParams)
	if err != nil {
		logrus.WithError(err).Error()
		// TODO: error
		return nil, err
	}
	err = json.UnmarshalFromString(s, &config.Model)
	if err != nil {
		logrus.WithError(err).Error()
		// TODO: error
		return nil, err
	}

	err = createModel(config.DB, config.Model)
	// data = make([]map[string]interface{}, 1)
	// data[0] = make(map[string]interface{})
	// data[0][keyData] = clone(config.Model)
	// c.Set(keyResults, data)
	return data, err
}

// PostHandle executes post.
func PostHandle(c *gin.Context) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}

	if err = c.ShouldBindJSON(config.Model); err != nil {
		return nil, err
	}

	return nil, createModel(config.DB, config.Model)
}

// PutHandle executes put.
func PutHandle(c *gin.Context) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}
	if err = c.ShouldBindJSON(config.Model); err != nil {
		return nil, err
	}
	return
}

// GetHandle executes actions and returns response
func GetHandle(c *gin.Context) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}

	set := foreignOfModel((*config).Model)

	paramsItr, ok := c.Get(KeyParams)
	if !ok {
		return nil, errors.New("can't find lazy params")
	}
	params := paramsItr.(Params)
	remain, page, limit, offset := separatePage(params)
	c.Set(KeyParams, remain)
	if limit == 0 {
		limit = 10000
	}

	var merged map[string][]string
	additional, ok := c.Get("_additional_values")
	if ok {
		merged = mergeValues(c.Request.URL.Query(), additional.(map[string][]string))
	}

	eq, gt, lt, gte, lte := URLValues(config.Model, merged)

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
	// logrus.WithField("count", count).Info()
	c.Set(keyCount, count)
	c.Set(keyData, config.Results)
	// c.Set(keyResults, map[string]interface{}{"count": count, "items": config.Results})
	return
}

// DefaultGetAction execute actions and returns response
func DefaultGetAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	var config *Configuration
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}

	paramsItr, ok := c.Get(KeyParams)
	if !ok {
		logrus.WithError(ErrParamMissing).Error()
		return nil, ErrParamMissing
	}
	params := paramsItr.(Params)
	filterParams, page, limit, offset := separatePage(params)
	c.Set(KeyParams, filterParams)
	if limit == 0 {
		limit = 10000
	}

	var mapResults []map[string]interface{}

	relations, err := relationships(config.DB, config.Model)
	if err != nil {
		return nil, err
	}
	modelSchema, err := schema.Parse(config.Model, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, err
	}

	eq, gt, lt, gte, lte := URLValues(config.Model, params)

	tx := config.DB.Model(config.Model)

	for k, v := range eq {
		tx = tx.Where(fmt.Sprintf("%s IN ?", k), v)
	}
	for k, v := range gt {
		tx = tx.Where(fmt.Sprintf("%s > ?", k), v)
	}
	for k, v := range lt {
		tx = tx.Where(fmt.Sprintf("%s < ?", k), v)
	}
	for k, v := range gte {
		tx = tx.Where(fmt.Sprintf("%s >= ?", k), v)
	}
	for k, v := range lte {
		tx = tx.Where(fmt.Sprintf("%s <= ?", k), v)
	}

	tx.Limit(int(limit)).Offset(int(limit*page + offset)).Find(&mapResults)

	count := int64(len(mapResults))
	if config.NeedCount {
		tx.Count(&count)
	}

	modelResults := make([]interface{}, len(mapResults))
	for k, v := range mapResults {
		if err := MapStruct(v, config.Model); err != nil {
			return nil, err
		}
		tmp := clone(config.Model)
		modelResults[k] = tmp
	}

	for _, vM2MRelation := range relations.Many2Many {
		primaryFieldValues := make([]interface{}, len(modelResults))
		if len(modelSchema.PrimaryFields) <= 0 {
			return nil, fmt.Errorf("primary fields not found")
		}
		primaryFieldName := modelSchema.PrimaryFields[0].Name
		for kModelResults, vModelResults := range modelResults {
			primaryFieldValue, err := valueOfField(vModelResults, primaryFieldName)
			if err != nil {
				return nil, err
			}
			primaryFieldValues[kModelResults] = primaryFieldValue
		}

		var joinTableResults []map[string]interface{}
		config.DB.Table(vM2MRelation.JoinTable.Table).Where(fmt.Sprintf("%s IN ?", vM2MRelation.JoinTable.DBNames[0]), primaryFieldValues).Find(&joinTableResults)

		mapForeignFieldValues := make(map[interface{}]bool)
		for _, vv := range joinTableResults {
			if value, ok := vv[vM2MRelation.JoinTable.DBNames[0]]; ok {
				mapForeignFieldValues[value] = true
			}
		}
		foreignFieldValues := make([]interface{}, 0)
		for k := range mapForeignFieldValues {
			foreignFieldValues = append(foreignFieldValues, k)
		}

		var foreignTableResults []map[string]interface{}
		config.DB.Table(vM2MRelation.FieldSchema.Table).Where(fmt.Sprintf("%s IN ?", vM2MRelation.FieldSchema.PrimaryFields[0].DBName), foreignFieldValues).Find(&foreignTableResults)

		assembleMany2Many(modelResults, joinTableResults, foreignTableResults,
			vM2MRelation.Schema.PrimaryFields[0].Name, vM2MRelation.FieldSchema.PrimaryFields[0].DBName,
			vM2MRelation.JoinTable.DBNames[0], vM2MRelation.JoinTable.DBNames[1],
			vM2MRelation.Field.StructField.Tag.Get("json"))

	}

	for _, vHasManyRelation := range relations.HasMany {
		primaryFieldValues := make([]interface{}, len(modelResults))

		if len(modelSchema.PrimaryFields) <= 0 {
			return nil, fmt.Errorf("primary fields not found")
		}
		primaryFieldName := modelSchema.PrimaryFields[0].Name

		for k, v := range modelResults {
			primaryFieldValue, err := valueOfField(v, primaryFieldName)
			if err != nil {
				return nil, err
			}
			primaryFieldValues[k] = primaryFieldValue
		}
		var hasManyResults []map[string]interface{}
		config.DB.Table(vHasManyRelation.References[0].ForeignKey.Schema.Table).Where(fmt.Sprintf("%s IN ?", vHasManyRelation.References[0].ForeignKey.DBName), primaryFieldValues).Find(&hasManyResults)

		assembleHasMany(modelResults, hasManyResults,
			vHasManyRelation.FieldSchema.PrimaryFields[0].Name, vHasManyRelation.References[0].ForeignKey.DBName,
			vHasManyRelation.Field.StructField.Tag.Get("json"))
	}

	logrus.WithField("count", count).Info()
	config.Results = modelResults
	c.Set(keyCount, count)
	c.Set(keyData, config.Results)
	c.Set(keyResults, map[string]interface{}{"count": count, "items": config.Results})
	return
}

// DefaultHTTPRequestAction ...
func DefaultHTTPRequestAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	requestConfig := actionConfig.Payload.(*HTTPRequest)
	logrus.Printf("%+v", requestConfig)
	var resp *grequests.Response
	ro := &grequests.RequestOptions{
		JSON: requestConfig.RequestBody,
	}
	if resp, err = grequests.Get(requestConfig.RequestURL, ro); err != nil {
		return nil, err
	}
	doc, err := jsonquery.Parse(strings.NewReader(resp.String()))
	node, err := jsonquery.Query(doc, "data")
	logrus.Print(err)
	logrus.Print(node.InnerText())

	return
}
