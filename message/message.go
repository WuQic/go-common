// The message package provides localization support to abstract
// end user messages for different languages
package message

import ()

type (
	// messageOp provides a data control structure for accessing
	// the message map
	messageOp struct {
		Operation string
		Response  chan string
		Data      interface{}
	}
)

var (
	// messageMap contains key/value pairs for all associated
	// end user messages that can be provided used as responses
	messageMap map[string]string

	// operationChannel is a channel used to perform operations on
	// the message map
	operationChannel chan *messageOp
)

// init is called once the package is used. Initialized the channel
// access system
func init() {
	// Channel to handle synchronous communication requests
	operationChannel = make(chan *messageOp)

	// Launch the goroutine to handle message map requests
	go messageAccess()
}

// Close will close the operational channel and release the goroutine
// handling requests
func Close() {
	operation := &messageOp{
		Operation: "Close",
		Response:  make(chan string),
	}

	operationChannel <- operation
	<-operation.Response
}

// LoadMessageMap takes the provided map and uses it for the
// message mappings
func LoadMessageMap(messageMap map[string]string) {
	operation := &messageOp{
		Operation: "Load",
		Data:      messageMap,
		Response:  make(chan string),
	}

	operationChannel <- operation
	<-operation.Response
}

// Message returns the string mapped to the specified key or
// it will return the key
func Message(key string) string {
	operation := &messageOp{
		Operation: "Find",
		Data:      key,
		Response:  make(chan string),
	}

	operationChannel <- operation
	return <-operation.Response
}

// messageAccess is the goroutine that handles all of the access
// to the message map based on communication using the operation
// channel
func messageAccess() {
	for {
		// Wait for an operation
		operation := <-operationChannel

		switch operation.Operation {
		case "Close":
			operation.Response <- "Closed"
			return

		case "Load":
			messageMap = operation.Data.(map[string]string)
			operation.Response <- "Loaded"

		case "Find":
			key := operation.Data.(string)
			value, ok := messageMap[key]
			if ok == false {
				value = key
			}
			operation.Response <- value
		}
	}
}
