package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"time"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var (
	jwtKey, _       = os.LookupEnv("SECRETKEY")
	sampleSecretKey = []byte(jwtKey)
)
var userMap = map[string]string{
	"key": "key",
}

func loginHandler(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case "POST":

		var user User
		err := json.NewDecoder(request.Body).Decode(&user)
		if err != nil {
			fmt.Fprintf(writer, "invalid body")
			return
		}

		if userMap[user.Username] == "" || userMap[user.Username] != user.Password {
			fmt.Fprintf(writer, "can not authenticate this user")
			return
		}
		token, err := generateJWT(user.Username)
		if err != nil {
			fmt.Fprintf(writer, "error in generating token")
		}
		fmt.Fprintf(writer, token)
	case "GET":
		fmt.Fprintf(writer, "only POST methods is allowed.")
		return
	}
}
func generateJWT(username string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = username
	claims["exp"] = time.Now().Add(time.Minute * 1).Unix()

	tokenString, err := token.SignedString(sampleSecretKey)

	if err != nil {
		fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
func validateToken(w http.ResponseWriter, r *http.Request) (err error) {
	if r.Header["Token"] == nil {
		log.Printf("invalid token error, %v", err)
		w.WriteHeader(500) // Return 500 Internal Server Error.
		return
	}

	token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("Unauthorized access for user: %v")
			w.WriteHeader(401) // Wrong password or username, Return 401.
			return nil, nil
		}
		return sampleSecretKey, nil
	})

	if token == nil {
		fmt.Fprintf(w, "invalid token")
		return errors.New("Token error")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Fprintf(w, "couldn't parse claims")
		return errors.New("Token error")
	}

	exp := claims["exp"].(float64)
	if int64(exp) < time.Now().Local().Unix() {
		fmt.Fprintf(w, "token expired")
		return errors.New("Token error")
	}

	return nil
}
