package helper

import (
	"encoding/json"
	"net/url"

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
