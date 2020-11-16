package lazy

import (
	"encoding/gob"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var (
	json        jsoniter.API
	schemaStore *sync.Map
	cacheStore  *sync.Map
)

// LogLevel ...
func LogLevel(level int) {
	logrus.SetLevel(logrus.Level(level))
}

func init() {
	schemaStore = &sync.Map{}
	cacheStore = &sync.Map{}
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	logrus.SetLevel(logrus.TraceLevel)
	gob.Register(map[string]interface{}{})
}
