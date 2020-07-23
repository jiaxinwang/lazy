package lazy

import (
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/gin-gonic/gin"
	"github.com/jiaxinwang/common/db"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// DefaultBeforeAction ...
func DefaultBeforeAction(c *gin.Context, gormDB *gorm.DB, config Configuration, payload interface{}) (result interface{}, reduce map[string][]string, err error) {
	v, _ := c.Get(keyParams)
	params := v.(Params)

	// eq, gt, lt, gte, lte, r := BeforeLazy(c.Request.URL.Query())
	eq, gt, lt, gte, lte, r := BeforeLazy(params)
	eqm := LazyTagSlice(config.Before.Model, eq)
	gtm := LazyTag(config.Before.Model, gt)
	ltm := LazyTag(config.Before.Model, lt)
	gtem := LazyTag(config.Before.Model, gte)
	ltem := LazyTag(config.Before.Model, lte)
	gormDB.LogMode(true)

	cols := make([]string, 0)

	for k := range config.Before.ResultMap {
		cols = append(cols, k)
	}
	colStr := strings.Join(cols, `,`)

	sel := db.SelectBuilder(sq.Select(colStr).From(config.Before.Table), eqm, gtm, ltm, gtem, ltem)
	result, err = db.Query(gormDB, sel)
	conv := result.([]map[string]interface{})
	logrus.Printf("%+v", conv)
	// if queryName, ok := config.Before.ResultMap[config.Before.Columms]; ok {
	// 	m := make(map[string][]string)
	// 	for _, v := range conv {
	// 		for mk, mv := range v {
	// 			if strings.EqualFold(mk, config.Before.Columms) {
	// 				if _, ok := m[queryName]; !ok {
	// 					m[queryName] = make([]string, 0)
	// 				}
	// 				m[queryName] = append(m[queryName], fmt.Sprintf("%+v", mv))
	// 			}
	// 		}
	// 	}
	// 	additionValues(c, m)
	// }

	return result, r, err
}
