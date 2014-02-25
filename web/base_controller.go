package web

import (
	"fmt"

	"github.com/ArdanStudios/go-common/helper"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"github.com/goinggo/tracelog"
)

type (
	BaseController struct {
		beego.Controller
	}

	MessageResponse struct {
		Messages []string
	}
)

const (
	VALIDATION_ERROR_MSG  = "Validation Error: The parameters provided were invalid."
	VALIDATION_ERROR_CODE = 409

	UNAUTHORIZED_ERROR_MSG  = "Invalid Credentials were supplied."
	UNAUTHORIZED_ERROR_CODE = 401

	APP_ERROR_MSG  = "An Application Error has occured."
	APP_ERROR_CODE = 500

	CACHE_CONTROL_HEADER = "Cache-control"
)

//Cache Output, outputs the cache control headrer for seconds passed in
func (this *BaseController) CacheOutput(seconds int64) {
	this.Ctx.Output.Header(CACHE_CONTROL_HEADER, fmt.Sprintf("private, must-revalidate, max-age=%d", seconds))
}

//Serve Empty Model {} as Json
func (this *BaseController) ServeBlankModel() {
	this.Data["json"] = map[string]string{}
	this.ServeJson()
}

//Serve Empty Array [] as Json
func (this *BaseController) ServeBlankModelList() {
	this.Data["json"] = []map[string]string{}
	this.ServeJson()
}

//Serve Model As Json
func (this *BaseController) ServeJsonModel(obj interface{}) {
	this.ServeJsonWithCache(obj, 0)
}

//Serve Model As Json
func (this *BaseController) ServeJsonWithCache(obj interface{}, secondsToCache int64) {
	if secondsToCache > 0 {
		this.CacheOutput(secondsToCache)
	}

	this.Data["json"] = obj
	this.ServeJson()
}

//ServeUnAuthorized returns an Unauthorized error
func (this *BaseController) ServeUnAuthorized() {
	tracelog.INFO("BaseController", "ServeUnAuthorized", "UnAuthorized, Exiting")

	this.ServeMessageWithStatus(UNAUTHORIZED_ERROR_CODE, UNAUTHORIZED_ERROR_MSG)
	return
}

//ServeValidationError returns a Validation Error's list of messages with a validation err code.
func (this *BaseController) ServeValidationError() {
	this.Ctx.Output.SetStatus(VALIDATION_ERROR_CODE)

	msgs := MessageResponse{}
	msgs.Messages = []string{VALIDATION_ERROR_MSG}
	this.Data["json"] = &msgs
	this.ServeJson()
}

//ServeValidationError returns a Validation Error's list of messages with a validation err code.
func (this *BaseController) ServeValidationErrors(validationErrors []*validation.ValidationError) {
	this.Ctx.Output.SetStatus(VALIDATION_ERROR_CODE)

	response := make([]string, len(validationErrors))
	for index, validationError := range validationErrors {
		response[index] = fmt.Sprintf("%s: %s", validationError.Field, validationError.String())
	}

	msgs := MessageResponse{}
	msgs.Messages = response
	this.Data["json"] = &msgs
	this.ServeJson()
}

//ServeUnAuthorized returns an Application error
func (this *BaseController) ServeAppError() {
	tracelog.INFO("BaseController", "ServeAppError", "Application Error, Exiting")

	this.ServeMessageWithStatus(APP_ERROR_CODE, APP_ERROR_MSG)
	return
}

func (this *BaseController) ServeMessageWithStatus(status int, msg string) {

	this.Ctx.Output.SetStatus(status)
	msgs := MessageResponse{}
	msgs.Messages = []string{msg}
	this.Data["json"] = &msgs
	this.ServeJson()
}

func (this *BaseController) ParseAndValidate(params interface{}) bool {
	err := this.ParseForm(params)
	if err != nil {
		this.ServeMessageWithStatus(VALIDATION_ERROR_CODE, VALIDATION_ERROR_MSG)
		return false
	}

	valid := validation.Validation{}
	ok, err := valid.Valid(params)
	if err != nil {
		this.ServeMessageWithStatus(VALIDATION_ERROR_CODE, VALIDATION_ERROR_MSG)
		return false
	}

	if ok == false {
		this.ServeValidationErrors(valid.Errors)
		return false
	}

	return true
}

func (this *BaseController) CatchPanic(err *error, UUID string, functionName string) {
	if helper.CatchPanic(err, UUID, functionName) {
		this.ServeAppError()
	}

}
