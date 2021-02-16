package util

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func ExtractParam(param string, r *http.Request) (string, bool) {
	value, ok := mux.Vars(r)[param]
	return value, ok
}

func ExtractParamInt(param string, r *http.Request) (int, bool) {
	if value, ok := ExtractParam(param, r); ok {
		if value, err := strconv.Atoi(value); err == nil {
			return value, true
		}
	}

	return 0, false
}

func ExtractFormValue(param string, r *http.Request) (string, bool) {
	value := r.FormValue("parents-count")
	return value, value != ""
}

func ExtractFormInt(param string, r *http.Request) (int, bool) {
	if value, ok := ExtractFormValue(param, r); ok {
		if value, err := strconv.Atoi(value); err == nil {
			return value, true
		}
	}

	return 0, false
}
