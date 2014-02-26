// The helper package provides generic or common functions required by
// the application
package helper

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"runtime"

	"github.com/goinggo/tracelog"
)

// LogModel writes a specified model object into the logs
// as a formatted JSON document
func LogModel(obj interface{}, useTrace bool) {
	bArray, _ := json.MarshalIndent(obj, "", "    ")

	if useTrace {
		tracelog.TRACE("utils", "LogModel", "Obj => \n\n%s\n\n", string(bArray))
	} else {
		tracelog.INFO("utils", "LogModel", "Obj => \n\n%s\n\n", string(bArray))
	}
}

// MergeUrlValues merges url values from one map to another
func MergeUrlValues(target *url.Values, source url.Values) {
	for key, _ := range source {
		target.Add(key, source.Get(key))
	}
}

// CatchPanic is used to catch any Panic and log exceptions to Stdout.
// It will also write the stack trace
func CatchPanic(err *error, goRoutine string, functionName string) bool {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		err2 := fmt.Errorf("PANIC Defered [%v] : Stack Trace : %v", r, string(buf))
		tracelog.ALERT("Unhandled Exception", goRoutine, functionName, err2.Error())

		if err != nil {
			*err = err2
		}
		return true
	}

	return false
}

// Round can round up 64 it floats to any specified precision
func Round(x float64, prec int) float64 {
	var rounder float64

	pow := math.Pow(10, float64(prec))
	intermed := x * pow

	_, frac := math.Modf(intermed)
	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow
}
