package common

const (
	ECONFLICT     = 409 // action cannot be performed
	EINTERNAL     = 500 // internal error
	EINVALID      = 400 // validation failed
	EUNAUTHORIZED = 401
	EFORBIDDEN    = 403
	ENOTFOUND     = 404 // entity does not exist

	EMINTERNAL = "An internal error occurred."
	EMINVALID  = "An error ocurred while processing the request."
	EMSEVERAL  = "One or more errors ocurred while processing the request."
)

type Error struct {
	Code int `json:"-"`

	Message string `json:"message"`
	Detail  string `json:"detail"`

	Errors map[string]string `json:"errors"`
}

func CreateGenericInternalError(err error) *Error {
	return &Error{Code: EINTERNAL, Message: EMINTERNAL, Detail: err.Error()}
}

func CreateGenericBadRequestError(err error) *Error {
	return &Error{Code: EINVALID, Message: EMINVALID, Detail: err.Error()}
}

func CreateBadRequestError(message string) *Error {
	return &Error{Code: EINVALID, Message: message}
}

func CreateConflictError(message string) *Error {
	return &Error{Code: ECONFLICT, Message: message}
}

func CreateFormError(errors map[string]string) *Error {
	return &Error{Code: EINVALID, Message: EMINVALID, Errors: errors}
}

func CreateNotFoundError(message string) *Error {
	return &Error{Code: ENOTFOUND, Message: message}
}
