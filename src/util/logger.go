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
	return log.WithField("component", component).
		WithField("color", getColor(component))
}

// Log format:
// 	All information before the double colon is metadata. This includes
//	Date, Time, Log Level, and Component all separated by spaces.
//	After the double colon is the actual log message.
func (f *logFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	ts := entry.Time.Format("2006-01-02 15:04:05.000")
	out := fmt.Sprintf("%v \x1b[%dm%v\x1b[0m %v :: %v\n",
		ts,
		entry.Data["color"].(int),
		strings.ToUpper(entry.Data["component"].(string)),
		strings.ToUpper(entry.Level.String()),
		entry.Message,
	)
	return []byte(out), nil
}

var colorCounter int
var colorMap = map[string]int{}

func getColor(c string) int {
	// Different packages can use the same component name.
	// These should keep consistent coloring.
	color, ok := colorMap[c]
	if ok {
		return color
	}

	switch colorCounter % 4 {
	case 0:
		color = 32 // green
	case 1:
		color = 33 // yellow
	case 2:
		color = 35 // magenta
	case 3:
		color = 36 // blue
	}

	colorMap[c] = color
	colorCounter++
	return color
}
