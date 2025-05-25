package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (app *application) internalServerError(w http.ResponseWriter,r *http.Request, err error){
	app.logger.Errorw("internal error" , "method", r.Method,"path",r.URL.Path , "error", err)
	 writeJSONError(w,http.StatusInternalServerError,"The Server encountered a problem")
	
}

func (app *application) Conflict(w http.ResponseWriter,r *http.Request, err error){
	app.logger.Errorw("Conflict error","method", r.Method , "path" , r.URL.Path ,"error" ,err)
	
	 writeJSONError(w,http.StatusConflict,"The was duplicate ")
	
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
app.logger.Warnf("bad request","method", r.Method , "path" , r.URL.Path ,"error" ,err)
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		var combined []string
		for _, fe := range ve {
			field := strings.ToLower(fe.Field())
			msg := field + " " + validationErrorMessage(fe)
			combined = append(combined, msg)
		}
		message := strings.Join(combined, ", ")
		writeJSONError(w, http.StatusBadRequest, message)
		return
	}

	writeJSONError(w, http.StatusBadRequest, err.Error())
}


func (app *application) notFound(w http.ResponseWriter,r *http.Request, err error){
	app.logger.Errorw("Resource not found:","method", r.Method , "path" , r.URL.Path ,"error" ,err)
	writeJSONError(w,http.StatusNotFound,"Resource not found")
}

func validationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return  "is required"
	case "max":
		return "must not be longer than " + fe.Param() + " characters"
	}
	return "is invalid"
}

func (app *application) unauthorizedError(w http.ResponseWriter,r *http.Request, err error){
	app.logger.Errorw("you are not authorized:","method", r.Method , "path" , r.URL.Path ,"error" ,err)
	writeJSONError(w,http.StatusUnauthorized,"you are not authorized")
}
func (app *application) basicUnauthorizedError(w http.ResponseWriter,r *http.Request, err error){
	app.logger.Errorw("you are not authorized:","method", r.Method , "path" , r.URL.Path ,"error" ,err)
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w,http.StatusUnauthorized,"you are not authorized")
}