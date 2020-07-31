package lazy

import "github.com/sirupsen/logrus"

var (
	logMode = false
)

// LogMode enable log
func LogMode(v bool) {
	logMode = v
}

func init() {
	logrus.SetLevel(logrus.TraceLevel)
}
