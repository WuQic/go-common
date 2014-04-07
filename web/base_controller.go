// The Web package provides common functionality for all controllers.
package web

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ArdanStudios/go-common/appErrors"
	"github.com/ArdanStudios/go-common/helper"
	"github.com/ArdanStudios/go-common/localize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/goinggo/tracelog"
)

type (
	// BaseController provides access to common controller.
	BaseController struct {
		beego.Controller
	}

	// MessageResponse provides the document structure for sending.
	// a list of messages
	MessageResponse struct {
		Messages []string `json:"messages"`
	}
)

const (
	CACHE_CONTROL_HEADER = "Cache-control"
)

// CacheOutput outputs the cache control header for seconds passed in.
func (baseController *BaseController) CacheOutput(seconds int64) {
	baseController.Ctx.Output.Header(CACHE_CONTROL_HEADER, fmt.Sprintf("private, must-revalidate, max-age=%d", seconds))
}

// ServeBlankModel serves an empty key/value pair map as Json.
func (baseController *BaseController) ServeBlankModel() {
	baseController.Data["json"] = map[string]string{}
	baseController.ServeJson()
}

// ServeBlankModelList serves an empty slice of key/value pair maps as Json.
func (baseController *BaseController) ServeBlankModelList() {
	baseController.Data["json"] = []map[string]string{}
	baseController.ServeJson()
}

// ServeJsonModel marshals the specified object as JSON.
func (baseController *BaseController) ServeJsonModel(obj interface{}) {
	baseController.ServeJsonWithCache(obj, 0)
}

// ServeJsonWithCache marshals the specified object as JSON specifying cache time.
func (baseController *BaseController) ServeJsonWithCache(obj interface{}, secondsToCache int64) {
	if secondsToCache > 0 {
		baseController.CacheOutput(secondsToCache)
	}

	baseController.Data["json"] = obj
	baseController.ServeJson()
}

// ServeImage serves an image with the specified mime type.
func (baseController *BaseController) ServeImage(image []byte, mimeType string) {
	baseController.Ctx.Output.SetStatus(200)
	baseController.Ctx.Output.ContentType("image/" + mimeType)
	baseController.Ctx.Output.Body([]byte{})
}

// ServeUnAuthorized returns an Unauthorized error.
func (baseController *BaseController) ServeUnAuthorized() {
	tracelog.INFO("BaseController", "ServeUnAuthorized", "UnAuthorized, Exiting")

	baseController.ServeMessageWithStatus(appErrors.UNAUTHORIZED_ERROR_CODE, localize.T(appErrors.UNAUTHORIZED_ERROR_MSG))
}

// ServeValidationError returns a Validation Error's list of messages with a validation err code.
func (baseController *BaseController) ServeValidationError() {
	baseController.Ctx.Output.SetStatus(appErrors.VALIDATION_ERROR_CODE)

	msgs := MessageResponse{}
	msgs.Messages = []string{localize.T(appErrors.VALIDATION_ERROR_MSG)}
	baseController.Data["json"] = &msgs
	baseController.ServeJson()
}

// ServeValidationErrors returns a Validation Error's list of messages with a validation err code.
func (baseController *BaseController) ServeValidationErrors(validationErrors []*validation.ValidationError) {
	baseController.Ctx.Output.SetStatus(appErrors.VALIDATION_ERROR_CODE)

	response := make([]string, len(validationErrors))
	for index, validationError := range validationErrors {
		response[index] = fmt.Sprintf("%s: %s", validationError.Field, validationError.String())
	}

	msgs := MessageResponse{}
	msgs.Messages = response
	baseController.Data["json"] = &msgs
	baseController.ServeJson()
}

// ServeError serves a error interface object.
func (baseController *BaseController) ServeError(err error) {
	switch e := err.(type) {
	case *appErrors.AppError:
		if e.ErrorCode() != 0 {
			baseController.ServeMessageWithStatus(e.ErrorCode(), e.Error())
			break
		}

		baseController.ServeMessageWithStatus(appErrors.APP_ERROR_CODE, e.Error())

	default:
		// We want to always return a generic message when an application error exists
		// We don't want to give the end user any information they could use against us
		baseController.ServeMessageWithStatus(appErrors.APP_ERROR_CODE, localize.T(appErrors.APP_ERROR_MSG))
	}
}

// ServeLocalizedError serves a validation error based on the specified key for the
// translated message.
func (baseController *BaseController) ServeLocalizedError(key string) {
	baseController.ServeMessageWithStatus(appErrors.VALIDATION_ERROR_CODE, localize.T(key))
}

// ServeAppError serves a generic application error.
func (baseController *BaseController) ServeAppError() {
	baseController.ServeMessageWithStatus(appErrors.APP_ERROR_CODE, localize.T(appErrors.APP_ERROR_MSG))
}

// ServeMessageWithStatus serves a HTTP status and message.
func (baseController *BaseController) ServeMessageWithStatus(status int, msg string) {
	baseController.ServeMessagesWithStatus(status, []string{msg})
}

// ServeMessageWithStatus serves a HTTP status and messages.
func (baseController *BaseController) ServeMessagesWithStatus(status int, msgs []string) {
	tracelog.INFO("BaseController", "ServeMessagesWithStatus", "Application Error, Exiting : %#v", msgs)

	baseController.Ctx.Output.SetStatus(status)
	response := MessageResponse{Messages: msgs}
	baseController.Data["json"] = &response
	baseController.ServeJson()
}

// ParseAndValidateJson is used to parse json into a type from the request and validate the values.
func (baseController *BaseController) ParseAndValidateJson(obj interface{}) bool {
	decoder := json.NewDecoder(baseController.Ctx.Request.Body)
	err := decoder.Decode(obj)
	if err != nil {
		baseController.ServeMessageWithStatus(appErrors.VALIDATION_ERROR_CODE, localize.T(appErrors.VALIDATION_ERROR_MSG))
		return false
	}

	return baseController.Validate(obj)
}

// ParseAndValidate is used to parse any form and query parameters from the request and validate the values.
func (baseController *BaseController) ParseAndValidate(obj interface{}) bool {
	err := baseController.ParseForm(obj)
	if err != nil {
		baseController.ServeMessageWithStatus(appErrors.VALIDATION_ERROR_CODE, localize.T(appErrors.VALIDATION_ERROR_MSG))
		return false
	}

	return baseController.Validate(obj)
}

// Validate validates a type against the valid tags in the type.
func (baseController *BaseController) Validate(params interface{}) bool {
	valid := validation.Validation{}
	ok, err := valid.Valid(params)
	if err != nil {
		baseController.ServeMessageWithStatus(appErrors.VALIDATION_ERROR_CODE, localize.T(appErrors.VALIDATION_ERROR_MSG))
		return false
	}

	if ok == false {
		// Build a map of the error messages for each field
		messages2 := map[string]string{}
		val := reflect.ValueOf(params).Elem()
		for i := 0; i < val.NumField(); i++ {
			// Look for an error tag in the field
			typeField := val.Type().Field(i)
			tag := typeField.Tag
			tagValue := tag.Get("error")

			// Was there an error tag
			if tagValue != "" {
				messages2[typeField.Name] = tagValue
			}
		}

		// Build the error response
		errors := []string{}
		for _, err := range valid.Errors {
			// Match an error from the validation framework errors
			// to a field name we have a mapping for
			message, ok := messages2[err.Field]
			if ok == true {
				// Use a localized message if one exists
				errors = append(errors, localize.T(message))
				continue
			} else {
				// No match, so use the message as is, Formats the err msg to include the key (field name).
				errors = append(errors, fmt.Sprintf("%s %s", err.Field, err.Message))
			}
		}

		baseController.ServeMessagesWithStatus(appErrors.VALIDATION_ERROR_CODE, errors)
		return false
	}

	return true
}

// CatchPanic is used to stop and process panics before they reach the Go runtime.
func (baseController *BaseController) CatchPanic(err *error, UUID string, functionName string) {
	if helper.CatchPanic(err, UUID, functionName) {
		baseController.ServeAppError()
	}
}

// CatchPanicNoErr is used to stop and process panics before they reach the Go runtime.
func (baseController *BaseController) CatchPanicNoErr(UUID string, functionName string) {
	if helper.CatchPanic(nil, UUID, functionName) {
		baseController.ServeAppError()
	}
}
