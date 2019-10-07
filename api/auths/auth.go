package auths

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("TopSecret")

//GenerateJWT will generate JSON Web Tokens
func GenerateJWT(name string, password string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"name":     name,
		"password": password,
	})
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		log.Fatal(err)
		fmt.Println("Error:", err)
		return "", err
	}

	return tokenString, nil
}

//ValidateMiddleware func
func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("X-API-Token")
		if authorizationHeader != "" {
			splitted := strings.Split(authorizationHeader, " ")

			if len(splitted) == 1 {
				tokenPart := splitted[0]
				token, error := jwt.Parse(tokenPart, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return mySigningKey, nil
				})
				if error != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				if token.Valid {
					ctx := context.WithValue(req.Context(), "player", token.Claims)
					req = req.WithContext(ctx)
					next(w, req)
				} else {
					w.WriteHeader(http.StatusNetworkAuthenticationRequired)
					return
				}
			}
		} else {
			w.WriteHeader(http.StatusNetworkAuthenticationRequired)
			return
		}
	})
}
