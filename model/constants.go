package model

const (
	USERNAME_MINIMUM_LENGTH = 8
	PASSWORD_MINIMUM_LENGTH = 8

	VERIFICATION_CODE_LENGTH     = 6
	VERIFICATION_CODE_EXPIRATION = 15 // maximum verification code valid time in minutes since its creation
	VERIFICATION_CODE_CHARS      = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

	SESSION_EXPIRATION = 24 // maximum user session valid time in hours since its creation
)

const (
	MSG_TODO_CREATED   = "The todo was successfully created."
	MSG_TODO_RETRIEVED = "The todo was successfully retrieved."
	MSG_TODO_UPDATED   = "The todo was successfully updated."
	MSG_TODO_DELETED   = "The todo was successfully deleted."
	MSG_TODOS_DELETED  = "The todos were successfully deleted."

	MSG_TODO_NOT_FOUND           = "Todo not found under given id (%d)."
	MSG_TODO_DESCRIPTION_MISSING = "Description field is empty or missing."
	MSG_TODO_IDS_NOT_PROVIDED    = "List of todo ids not provided."

	MSG_SESSION_CREATED    = "The session was successfully created."
	MSG_INVALID_CREDENTIAL = "Invalid username or password."

	MSG_USER_CREATED   = "The user was successfully created."
	MSG_USER_RETRIEVED = "The user was successfully retrieved."

	MSG_USER_NOT_AUTHORIZED   = "The user is not authorized."
	MSG_USER_VERIFIED         = "The user was successfully verified."
	MSG_USER_ALREADY_VERIFIED = "The user is already verified."

	MSG_USER_USERNAME_MISSING = "Username field is empty or missing."
	MSG_USER_USERNAME_LENGTH  = "Username must have %d characters or more."
	MSG_USER_EMAIL_MISSING    = "Email field is empty or missing."
	MSG_USER_EMAIL_INVALID    = "Email is invalid"
	MSG_USER_PASSWORD_MISSING = "Password field is empty or missing."
	MSG_USER_PASSWORD_LENGTH  = "Password must have %d characters or more."

	MSG_USER_NOT_FOUND           = "User not found under given id (%d)."
	MSG_USERNAME_NOT_FOUND       = "User not found under given username (%s)."
	MSG_USER_GOOGLE_ID_NOT_FOUND = "User not found under given Google id (%s)."

	MSG_USER_USERNAME_ALREADY_EXISTS = "Username already exists."
	MSG_USER_EMAIL_ALREADY_EXISTS    = "Email already exists."

	MSG_INVALID_CODE           = "Given verification code is invalid or expired."
	MSG_VERIFICATION_SENT      = "A new verification code was just sent to your mailbox."
	MSG_VERIFICATION_NOT_FOUND = "User verification request not found under given user id (%d)."

	MSG_GOOGLE_UNAVAILABLE = "Google authentication is not available at the moment."
	MSG_GOOGLE_TOKEN_ERROR = "An error occurred while trying to validate the token."

	MSG_TOKEN_NOT_PROVIDED = "Token not provided."

	MSG_ERR_INTERNAL = "An internal error occurred."
	MSG_ERR_INVALID  = "An error ocurred while processing the request."
	MSG_ERR_SEVERAL  = "One or more errors ocurred while processing the request."
)
