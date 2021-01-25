package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jvitoroc/todo-api/resources/repo"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type UserRequestBody struct {
	Username *string `json:"username"`
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func initializeUser(r *mux.Router) {
	addUserHandlers(r)
}

func addUserHandlers(r *mux.Router) {
	sr := r.PathPrefix("/user").Subrouter()
	sr.Handle("/", appHandler(createUserHandler)).Methods("POST")
	sr.Handle("/session", appHandler(createUserSessionHandler)).Methods("POST")

	verificationRouter := sr.NewRoute().Subrouter()
	verificationRouter.Use(authenticateRequest)
	verificationRouter.Handle("/verification", appHandler(verifiyUserHandler)).Methods("POST")
	verificationRouter.Handle("/verification/resend", appHandler(resendEmailHandler)).Methods("POST")

	protectedRouter := sr.NewRoute().Subrouter()
	protectedRouter.Use(authenticateRequest)
	protectedRouter.Use(checkActivationState)
	protectedRouter.Handle("/", appHandler(getUserHandler)).Methods("GET")
}

func createUserHandler(w http.ResponseWriter, r *http.Request) *appError {
	requestBody := UserRequestBody{}

	if err := extractUser(&requestBody, r); err != nil {
		return err
	}

	if err := validateUser(false, false, &requestBody); err != nil {
		return err
	}

	var passwordHash string
	if err := generatePasswordHash(&passwordHash, *requestBody.Password, w); err != nil {
		return err
	}

	id, err := repo.InsertUser(db, *requestBody.Username, *requestBody.Email, passwordHash)

	if err != nil {
		switch err.(sqlite3.Error).Code {
		case 19:
			return createAlreadyExistsError(err)
		default:
			return unknownAppError(err)
		}
	}

	var user *repo.User

	if user, err = repo.GetUser(db, *id); err != nil {
		return unknownAppError(err)
	}

	if err := sendNewVerificationCode(*id); err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "User successfully created. An email with a verification code was just sent to your email address.",
			"data":    user,
		},
		http.StatusCreated,
		w,
	)
	return nil
}

func createUserSessionHandler(w http.ResponseWriter, r *http.Request) *appError {
	requestBody := UserRequestBody{}

	if err := extractUser(&requestBody, r); err != nil {
		return err
	}

	if err := validateUser(true, true, &requestBody); err != nil {
		return err
	}

	user, err := repo.GetUserByUsername(db, *requestBody.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			return createAppError("Username or password is incorrect.", http.StatusBadRequest)
		} else {
			return createAppError(fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), http.StatusBadRequest)
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*requestBody.Password))
	if err != nil {
		return createAppError("Username or password is incorrect.", http.StatusBadRequest)
	}

	token, err := createToken(user.ID)
	if err != nil {
		return createAppError(fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), http.StatusBadRequest)
	}

	respond(
		map[string]interface{}{
			"message": "Session successfully created.",
			"data": map[string]interface{}{
				"token": token,
				"user":  user,
			},
		},
		http.StatusCreated,
		w,
	)
	return nil
}

func verifiyUserHandler(w http.ResponseWriter, r *http.Request) *appError {
	userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)
	requestBody := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return createAppError(fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), http.StatusBadRequest)
	}

	var verificationCode string
	var ok bool
	if verificationCode, ok = requestBody["verificationCode"]; !ok {
		return createAppError(fmt.Sprintf("Verification code not provided."), http.StatusBadRequest)
	}

	var request *repo.UserActivationRequest
	var err error
	if request, err = repo.GetUserActivationRequest(db, userId); err != nil {
		return createAppError(fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), http.StatusInternalServerError)
	}

	if verificationCode != request.Code || time.Now().After(request.ExpiresAt) {
		return createAppError("Given verification code is invalid or expired.", http.StatusBadRequest)
	}

	if _, err = repo.UpdateUser(db, userId, map[string]interface{}{"active": true}); err != nil {
		return createAppError(fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), http.StatusInternalServerError)
	}

	respondWithMessage("User successfully verified.", http.StatusOK, w)
	return nil
}

func resendEmailHandler(w http.ResponseWriter, r *http.Request) *appError {
	userId, _ := strconv.ParseInt(r.Context().Value("userId").(string), 10, 64)

	var user *repo.User
	var err error

	if user, err = repo.GetUser(db, userId); err != nil {
		return unknownAppError(err)
	}

	if user.Active {
		return createAppError("Your account is already verified.", http.StatusConflict)
	}

	if err := sendNewVerificationCode(userId); err != nil {
		return err
	}

	respondWithMessage("A new verification code was just sent to your mailbox.", http.StatusOK, w)
	return nil
}

func getUserHandler(w http.ResponseWriter, r *http.Request) *appError {
	user := r.Context().Value("user").(*repo.User)

	respond(
		map[string]interface{}{
			"message": "User successfuly retrieved.",
			"data":    user,
		},
		http.StatusOK,
		w,
	)
	return nil
}

func extractUser(user *UserRequestBody, r *http.Request) *appError {
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		return createAppError(fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), http.StatusBadRequest)
	}
	return nil
}

func validateUser(onlyCheckExistence bool, ignoreEmail bool, user *UserRequestBody) *appError {
	errors := map[string]string{}

	if user.Username == nil || *user.Username == "" {
		errors["username"] = "Username field is empty or missing."
	} else if !onlyCheckExistence && len(*user.Username) < 8 {
		errors["username"] = "Username must have 8 characters or more."
	}

	if !ignoreEmail {
		if user.Email == nil || *user.Email == "" {
			errors["email"] = "Email field is empty or missing."
		} else if !isEmailValid(*user.Email) {
			errors["email"] = "Email is invalid."
		}
	}

	if user.Password == nil || *user.Password == "" {
		errors["password"] = "Password field is empty or missing."
	} else if !onlyCheckExistence && len(*user.Password) < 8 {
		errors["password"] = "Password must have 8 characters or more."
	}

	if len(errors) == 0 {
		return nil
	} else {
		return createMappedAppError(MSG_ONE_MORE_ERRORS, errors, http.StatusBadRequest)
	}
}

func createAlreadyExistsError(err error) *appError {
	msg := err.Error()
	errors := map[string]string{}

	if strings.Contains(msg, "user.username") {
		errors["username"] = "Username already exists."
	} else if strings.Contains(msg, "user.email") {
		errors["email"] = "Email already exists."
	}

	return createMappedAppError(
		MSG_ONE_MORE_ERRORS,
		errors,
		http.StatusConflict,
	)
}
