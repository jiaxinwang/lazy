package lazy

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
)

var (
	logMode = false
	json    jsoniter.API
)

// LogMode enable log
func LogMode(v bool) {
	logMode = v
}

func init() {
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	logrus.SetLevel(logrus.TraceLevel)
}
