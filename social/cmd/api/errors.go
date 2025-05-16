package main

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (app *application) internalServerError(w http.ResponseWriter,r *http.Request, err error){
	log.Panicf("internal server error: %s path: %s error: %s", r.Method,r.URL.Path,err.Error() )
	writeJSONError(w,http.StatusInternalServerError,"The Server encountered a problem")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Invalid request: %s %s: %s", r.Method, r.URL.Path, err.Error())

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
	log.Panicf("Resource not found: %s path: %s error: %s", r.Method,r.URL.Path,err.Error() )
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