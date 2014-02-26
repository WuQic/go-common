package helper

import (
	"encoding/json"
	"fmt"
	"net/url"
	"runtime"

	"github.com/goinggo/tracelog"
)

//Logs a model in formatted JSON
func LogModel(obj interface{}, useTrace bool) {
	bArray, _ := json.MarshalIndent(obj, "", "    ")

	if useTrace {
		tracelog.TRACE("utils", "LogModel", "Obj => \n\n%s\n\n", string(bArray))
	} else {
		tracelog.INFO("utils", "LogModel", "Obj => \n\n%s\n\n", string(bArray))
	}
}

//Merges url values from one map to another
func MergeUrlValues(target *url.Values, source url.Values) {
	for key, _ := range source {
		target.Add(key, source.Get(key))
	}
}

//Catches Panic
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
