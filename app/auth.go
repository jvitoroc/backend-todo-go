package app

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/jvitoroc/todo-go/model"
)

func (app *App) VerifyToken(token string) (int, *model.AppError) {
	tk, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(app.Config.Auth.JwtSecret), nil
	})

	if err != nil {
		return 0, model.NewUnauthorizedError(err.Error())
	}

	if claims, ok := tk.Claims.(jwt.MapClaims); ok && tk.Valid {
		userId, _ := strconv.Atoi(claims["userId"].(string))
		return userId, nil
	} else {
		return 0, model.NewUnauthorizedError(model.MSG_USER_NOT_AUTHORIZED)
	}
}

func (app *App) SignToken(userId int) (string, *model.AppError) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(userId),
		"exp":    time.Now().Add(time.Hour * model.SESSION_EXPIRATION).Unix(),
	})

	tokenString, err := token.SignedString([]byte(app.Config.Auth.JwtSecret))
	if err != nil {
		return "", model.NewGenericBadRequestError(err)
	}

	return tokenString, nil
}
