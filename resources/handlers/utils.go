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
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jvitoroc/todo-api/resources/common"
	"github.com/jvitoroc/todo-api/resources/repo"
	"golang.org/x/crypto/bcrypt"
)

const CHARS = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type appHandler func(http.ResponseWriter, *http.Request) *common.Error

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		respondWithError(*err, w)
	}
}

func respond(data interface{}, statusCode int, w http.ResponseWriter) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func respondWithMessage(message string, statusCode int, w http.ResponseWriter) {
	respond(map[string]string{"message": message}, statusCode, w)
}

func respondWithError(err common.Error, w http.ResponseWriter) {
	respond(err, err.Code, w)
}

func generatePasswordHash(passwordHash *string, password string, w http.ResponseWriter) *common.Error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return common.CreateGenericInternalError(err)
	}
	*passwordHash = string(hash)
	return nil
}

func createToken(userId int) (string, *common.Error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(userId),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", common.CreateGenericBadRequestError(err)
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

	lastIndex := len(CHARS) - 1
	for i := 0; i < 6; i++ {
		code = code + string(CHARS[r.Intn(lastIndex)])
	}

	return code
}

func sendNewVerificationCode(userId int) *common.Error {
	expiresAt := time.Now().Local().Add(time.Minute * time.Duration(15))
	code := generateActivationCode()

	var request *repo.UserActivationRequest
	var err *common.Error
	var user *repo.User

	if request, err = repo.UpsertUserActivationRequest(db, userId, code, expiresAt); err != nil {
		return err
	}

	if user, err = repo.GetUser(db, userId); err != nil {
		return err
	}

	if err := sendActivationEmail(user.Username, user.Email, request.Code, request.ExpiresAt); err != nil {
		return common.CreateGenericInternalError(err)
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

func extractContextValue(param string, r *http.Request) interface{} {
	return r.Context().Value(param)
}

func extractContextInt(param string, r *http.Request) (int, error) {
	return strconv.Atoi(extractContextValue(param, r).(string))
}

func extractParam(param string, r *http.Request) (string, bool) {
	value, ok := mux.Vars(r)[param]
	return value, ok
}

func extractParamInt(param string, r *http.Request) (int, bool) {
	if value, ok := extractParam(param, r); ok {
		if value, err := strconv.Atoi(value); err == nil {
			return value, true
		}
	}

	return 0, false
}
