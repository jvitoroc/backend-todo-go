package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jvitoroc/todo-api/resources/repo"
	"golang.org/x/crypto/bcrypt"
)

const CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

/* << Pre-defined messages */
var MSG_UNKNOWN_ERROR = "An unknown error ocurred while processing the request: %s"
var MSG_ONE_MORE_ERRORS = "One or more errors ocurred while processing the request."
var MSG_NOT_FOUND_ERROR = "%s not found under the given id (%d)."

/* Pre-defined messages >> */

type appError struct {
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
	Code    int               `json:"-"`
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		respond(e, e.Code, w)
	}
}

func unknownAppError(err error) *appError {
	return &appError{fmt.Sprintf(MSG_UNKNOWN_ERROR, err.Error()), nil, http.StatusInternalServerError}
}

func createAppError(message string, code int) *appError {
	return &appError{message, nil, code}
}

func createMappedAppError(message string, errors map[string]string, code int) *appError {
	return &appError{message, errors, code}
}

func respond(data interface{}, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondWithMessage(message string, statusCode int, w http.ResponseWriter) {
	respond(map[string]string{"message": message}, statusCode, w)
}

func generatePasswordHash(passwordHash *string, password string, w http.ResponseWriter) *appError {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return unknownAppError(err)
	}
	*passwordHash = string(hash)
	return nil
}

func createToken(userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(userId),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func parseToken(token string) (jwt.MapClaims, error) {
	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func generateActivationCode() string {
	code := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var mu sync.Mutex

	mu.Lock()
	lastIndex := len(CHARS) - 1
	for i := 0; i < 6; i++ {
		code = code + string(CHARS[r.Intn(lastIndex)])
	}
	mu.Unlock()

	return code
}

func sendNewVerificationCode(userId int64) *appError {
	expiresAt := time.Now().Local().Add(time.Minute * time.Duration(15))
	code := generateActivationCode()

	var user *repo.User
	var err error

	if err := repo.UpsertUserActivationRequest(db, userId, code, expiresAt); err != nil {
		return unknownAppError(err)
	}

	if user, err = repo.GetUser(db, userId); err != nil {
		return unknownAppError(err)
	}

	if err := sendActivationEmail(user.Username, user.Email, code, expiresAt); err != nil {
		return unknownAppError(err)
	}

	return nil
}

func sendActivationEmail(username string, to string, code string, expirationDate time.Time) error {
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpAddr := os.Getenv("SMTP_ADDR")

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := []byte("To: " + to + "\r\n" +
		"Subject: Todo app: activate your account!\r\n" +
		"\r\n" +
		"Hey " + username + ", here's the code needed to activate your account: " + code + "\r\n" +
		"It will expire on " + expirationDate.Format(time.RFC1123))

	err := smtp.SendMail(smtpAddr, auth, smtpUser, []string{to}, msg)

	return err
}

func isEmailValid(e string) bool {
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	if !emailRegex.MatchString(e) {
		return false
	}

	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}

	return true
}
