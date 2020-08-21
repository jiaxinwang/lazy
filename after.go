package lazy

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// DefaultAfterAction ...
func DefaultAfterAction(c *gin.Context, gormDB *gorm.DB, config Configuration, payload interface{}) (result interface{}, reduce map[string][]string, err error) {
	return
}
