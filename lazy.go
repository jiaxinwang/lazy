package lazy

import (
	"encoding/gob"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var (
	logMode     = false
	json        jsoniter.API
	schemaStore *sync.Map
	cacheStore  *sync.Map
)

// LogMode enable log
func LogMode(v bool) {
	// TODO:
	logMode = v
}

func init() {
	schemaStore = &sync.Map{}
	cacheStore = &sync.Map{}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	logrus.SetLevel(logrus.TraceLevel)
	gob.Register(map[string]interface{}{})
}
