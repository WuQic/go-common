// The decoder package provides support for decoding JSON documents
// into Go native types that have the appropriate tags
package decoder

import (
	"encoding/json"
	"strings"

	"github.com/ArdanStudios/go-common/helper"
	"github.com/goinggo/mapstructure"
	"github.com/goinggo/tracelog"
)

// IsArrayResponse validates if the specified JSON document
// is an array or slice
func IsArrayResponse(doc []byte) bool {
	tracelog.STARTED("utils", "IsArrayResponse")

	docString := string(doc)
	docSlice := strings.TrimLeft(docString, " ")

	if len(docSlice) > 0 && string(docSlice[0]) == "[" {
		tracelog.COMPLETEDf("utils", "IsArrayReponse", "Doc is Array")
		return true
	}

	tracelog.COMPLETEDf("utils", "IsArrayResponse", "Doc is not Array")

	return false
}

// DecodeMap decodes a JSON document into a Go native object
func DecodeMap(docMap map[string]interface{}, obj interface{}) error {
	if err := mapstructure.DecodePath(docMap, obj); err != nil {
		tracelog.ERROR(err, "utils", "Decode, Decoding Mapped Doc")
		return err
	}

	if tracelog.LogLevel() == tracelog.LEVEL_TRACE {
		helper.LogModel(obj, true)
	}

	return nil
}

// Decode takes a JSON document and decodes that into a Go native object
func Decode(doc []byte, obj interface{}) error {
	tracelog.STARTED("utils", "Decode")
	docMap := map[string]interface{}{}

	if err := json.Unmarshal(doc, &docMap); err != nil {
		tracelog.ERROR(err, "utils", "Decode, Building Mapped Doc")
		return err
	}

	return DecodeMap(docMap, obj)
}

// DecodeSlice decodes a JSON document array into a slice of Go native objects
func DecodeSlice(doc []byte, sliceObj interface{}, obj interface{}) (bool, error) {
	tracelog.STARTED("utils", "DecodeSlice")

	if IsArrayResponse(doc) == false {
		//decode as struct
		if err := Decode(doc, obj); err != nil {
			tracelog.ERROR(err, "utils", "DecodeSlice, Item Not Array, Unable to decode as struct")
			return false, err
		}
		//return false since not an array
		return false, nil
	}

	sliceMap := []map[string]interface{}{}

	if err := json.Unmarshal(doc, &sliceMap); err != nil {
		tracelog.ERROR(err, "utils", "DecodeSlice")
		return false, err
	}

	if err := mapstructure.DecodeSlicePath(sliceMap, sliceObj); err != nil {
		tracelog.ERROR(err, "utils", "DecodeSlice, Decoding Slice Object")
		return false, err
	}

	tracelog.COMPLETED("utils", "DecodeSlice")
	return true, nil
}
