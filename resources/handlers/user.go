package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jvitoroc/todo-api/resources/common"
	"github.com/jvitoroc/todo-api/resources/repo"
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

func createUserHandler(w http.ResponseWriter, r *http.Request) *common.Error {
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

	var user *repo.User
	var err *common.Error

	if user, err = repo.InsertUser(db, *requestBody.Username, *requestBody.Email, passwordHash); err != nil {
		return err
	}

	if err := sendNewVerificationCode(user.ID); err != nil {
		return err
	}

	respond(
		map[string]interface{}{
			"message": "User successfully created.",
			"detail":  "An email with a verification code was just sent to your email address.",
			"data":    user,
		},
		http.StatusCreated,
		w,
	)
	return nil
}

func createUserSessionHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	requestBody := UserRequestBody{}

	if err := extractUser(&requestBody, r); err != nil {
		return err
	}

	if err := validateUser(true, true, &requestBody); err != nil {
		return err
	}

	var user *repo.User
	var err *common.Error

	if user, err = repo.GetUserByUsername(db, *requestBody.Username); err != nil {
		if err.Code == common.ENOTFOUND {
			err.Message = "Username or password is incorrect."
			err.Code = common.EINVALID
		}
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(*requestBody.Password)); err != nil {
		return common.CreateBadRequestError("Username or password is incorrect.")
	}

	token, err := createToken(user.ID)
	if err != nil {
		return err
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

func verifiyUserHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))
	requestBody := map[string]string{}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return common.CreateGenericBadRequestError(err)
	}

	var verificationCode string
	var ok bool

	if verificationCode, ok = requestBody["verificationCode"]; !ok {
		return common.CreateBadRequestError("Verification code not provided.")
	}

	var request *repo.UserActivationRequest
	var err *common.Error

	if request, err = repo.GetUserActivationRequest(db, userId); err != nil {
		return err
	}

	if verificationCode != request.Code || time.Now().After(request.ExpiresAt) {
		return common.CreateBadRequestError("Given verification code is invalid or expired.")
	}

	if err = repo.UpdateUser(db, &repo.User{ID: userId, Active: true}); err != nil {
		return err
	}

	respondWithMessage("User successfully verified.", http.StatusOK, w)
	return nil
}

func resendEmailHandler(w http.ResponseWriter, r *http.Request) *common.Error {
	userId, _ := strconv.Atoi(r.Context().Value("userId").(string))

	var user *repo.User
	var err *common.Error

	if user, err = repo.GetUser(db, userId); err != nil {
		return err
	}

	if user.Active {
		return common.CreateConflictError("Your account is already verified.")
	}

	if err := sendNewVerificationCode(userId); err != nil {
		return err
	}

	respondWithMessage("A new verification code was just sent to your mailbox.", http.StatusOK, w)
	return nil
}

func getUserHandler(w http.ResponseWriter, r *http.Request) *common.Error {
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

func extractUser(user *UserRequestBody, r *http.Request) *common.Error {
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		return common.CreateGenericBadRequestError(err)
	}
	return nil
}

func validateUser(checkWhetherExistsOnly bool, ignoreEmail bool, user *UserRequestBody) *common.Error {
	errors := map[string]string{}

	if user.Username == nil || *user.Username == "" {
		errors["username"] = "Username field is empty or missing."
	} else if !checkWhetherExistsOnly && len(*user.Username) < 8 {
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
	} else if !checkWhetherExistsOnly && len(*user.Password) < 8 {
		errors["password"] = "Password must have 8 characters or more."
	}

	if len(errors) == 0 {
		return nil
	} else {
		return common.CreateFormError(errors)
	}
}
