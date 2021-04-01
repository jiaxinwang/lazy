package lazy

import (
	"encoding/gob"
	"sync"

	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
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
	extra.RegisterFuzzyDecoders()
	extra.SetNamingStrategy(extra.LowerCaseWithUnderscores)
	logrus.SetLevel(logrus.TraceLevel)
	gob.Register(map[string]interface{}{})
}
