package util

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type logFormatter struct{}

// Get custom logger
func LoggerFor(component string) *logrus.Entry {
	log := logrus.New()
	log.SetFormatter(new(logFormatter))
	return log.WithField("component", component)
}

func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	ts := entry.Time.Format("2006-01-02 15:04:05")
	out := fmt.Sprintf("%v %v %v :: %v\n",
		ts,
		strings.ToUpper(entry.Level.String()),
		strings.ToUpper(entry.Data["component"].(string)),
		entry.Message,
	)
	return []byte(out), nil
}
