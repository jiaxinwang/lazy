package lazy

import (
	"fmt"
	"strings"

	"github.com/antchfx/jsonquery"
	"github.com/levigross/grequests"
	"gorm.io/gorm/schema"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// DefaultPostAction execute default post.
func DefaultPostAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	_, _, bodyParams := ContentParams(c)
	config, err := ConfigurationWithContext(c)
	if err != nil {
		return nil, err
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

// ConfigurationWithContext ...
func ConfigurationWithContext(c *gin.Context) (*Configuration, error) {
	var config *Configuration
	if v, ok := c.Get(KeyConfig); ok {
		config = v.(*Configuration)
	} else {
		return nil, ErrConfigurationMissing
	}
	return config, nil
}

// DefaultGetAction execute actions and returns response
func DefaultGetAction(c *gin.Context, actionConfig *Action, payload interface{}) (data []map[string]interface{}, err error) {
	config, err := ConfigurationWithContext(c)
	if err != nil {
		return nil, err
	}

	paramsItr, ok := c.Get(KeyParams)
	if !ok {
		logrus.WithError(ErrParamMissing).Error()
		return nil, ErrParamMissing
	}
	params := paramsItr.(map[string][]string)
	filterParams, _, _, _ := separatePage(params)
	c.Set(KeyParams, filterParams)

	var mapResults []map[string]interface{}

	relations, err := relationships(config.DB, config.Model)
	if err != nil {
		return nil, err
	}
	modelSchema, err := schema.Parse(config.Model, schemaStore, schema.NamingStrategy{})
	if err != nil {
		return nil, err
	}

	tx := config.DB.Model(config.Model)

	qParams, err := splitQueryParams(config.Model, params)
	if err != nil {
		logrus.WithError(err).Error()
	}

	for _, v := range qParams.Many2Many {
		var ids []int64
		config.DB.Table(v.JoinTable).Where(fmt.Sprintf("%s in ?", v.JoinTableForeignFieldName), v.Values).Pluck(v.JoinTableFieldName, &ids)
		if len(ids) > 0 {
			qParams.Eq["id"] = make([]interface{}, 0)
			for _, vIds := range ids {
				qParams.Eq["id"] = append(qParams.Eq["id"], vIds)
			}
		}
	}

	for k, v := range qParams.Eq {
		if len(v) == 1 {
			tx = tx.Where(fmt.Sprintf("%s = ?", k), v[0])
		} else {
			tx = tx.Where(fmt.Sprintf("%s IN ?", k), v)
		}

	}
	for k, v := range qParams.Gt {
		tx = tx.Where(fmt.Sprintf("%s > ?", k), v)
	}
	for k, v := range qParams.Lt {
		tx = tx.Where(fmt.Sprintf("%s < ?", k), v)
	}
	for k, v := range qParams.Gte {
		tx = tx.Where(fmt.Sprintf("%s >= ?", k), v)
	}
	for k, v := range qParams.Lte {
		tx = tx.Where(fmt.Sprintf("%s <= ?", k), v)
	}

	needGroup := false
	dotName := ""
	name := ""
	ptable := ""

	for _, v := range qParams.HasMany {
		tx = tx.Joins(v.Table).Where(fmt.Sprintf("%s IN ?", v.Name), v.Values)
		dotName = v.DotName
		name = v.Name
		ptable = v.PTable
		needGroup = true
	}

	count := int64(len(mapResults))
	if config.NeedCount {
		if needGroup {
			tx = tx.Group(ptable)
			type Count struct {
				cnt int64
			}
			var c Count
			tx.Select(fmt.Sprintf("%s as %s,count(%s) as cnt", dotName, name, ptable)).Scan(&c)
		} else {
			tx.Count(&count)
		}
	}

	tx.Limit(int(qParams.Limit)).Offset(int(qParams.Limit*qParams.Page + qParams.Offset)).Find(&mapResults)

	modelResults := make([]interface{}, len(mapResults))
	for k, v := range mapResults {
		if len(qParams.Ignore) != 0 {
			for k := range qParams.Ignore {
				delete(v, k)
				logrus.Print(v)
			}
		}
		if err := MapStruct(v, config.Model); err != nil {
			return nil, err
		}
		tmp := clone(config.Model)
		modelResults[k] = tmp
	}
	for _, vBelongToRelation := range relations.BelongsTo {
		ref := vBelongToRelation.References[0]
		foreignKeyName := ref.ForeignKey.Name
		primaryFieldValues := make([]interface{}, len(modelResults))
		if len(ref.PrimaryKey.Schema.PrimaryFieldDBNames) <= 0 {
			return nil, fmt.Errorf("primary fields not found")
		}

		for kModelResults, vModelResults := range modelResults {
			primaryFieldValue, err := valueOfField(vModelResults, foreignKeyName)
			if err != nil {
				return nil, err
			}
			primaryFieldValues[kModelResults] = primaryFieldValue
		}

		var belongToResults []map[string]interface{}
		config.DB.Table(ref.PrimaryKey.Schema.Table).Where(fmt.Sprintf("%s IN ?", ref.PrimaryKey.Schema.PrimaryFieldDBNames[0]), primaryFieldValues).Find(&belongToResults)

		assembleBelongTo(modelResults, belongToResults,
			ref.ForeignKey.Name, ref.PrimaryKey.Schema.PrimaryFieldDBNames[0],
			vBelongToRelation.Field.StructField.Tag.Get("json"))

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
			if value, ok := vv[vM2MRelation.JoinTable.DBNames[1]]; ok {
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

	logrus.WithField("count", count).Trace()
	c.Set(keyCount, count)
	c.Set(keyData, modelResults)
	c.Set(keyResults, map[string]interface{}{"count": count, "items": modelResults})
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
