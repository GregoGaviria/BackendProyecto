package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID      int    `json:"userID"`
	NivelAcceso int    `json:"nivelAcesso"`
	Username    string `json:"username"`
	jwt.StandardClaims
}

func validarLogin(username string, password string) (cookie *http.Cookie, nivelAcceso int, status int, err error) {

	if username == "" || password == "" {
		return nil, 0, http.StatusBadRequest, nil
	}
	if password == "*" {
		return nil, 0, http.StatusBadRequest, nil
	}

	// expectedPassword, err := getPasswordById(cedula)
	// if err != nil || expectedPassword != clave {
	// 	log.Println(err)
	// 	return nil, http.StatusUnauthorized
	// }

	var userID int
	var expectedPassword string

	row := db.QueryRow("SELECT UsuarioId, NivelAcceso, Password FROM Usuarios WHERE Username=?", username)
	err = row.Scan(&userID, &nivelAcceso, &expectedPassword)

	if err != nil || password != expectedPassword {
		return nil, 0, http.StatusUnauthorized, err
	}

	expirationTime := time.Now().Add(time.Minute * 60)
	var claims *Claims = &Claims{
		UserID:      userID,
		NivelAcceso: nivelAcceso,
		Username:    username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	var token *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}
	cookie = &http.Cookie{
		Name:    "token",
		Value:   tokenstring,
		Expires: expirationTime,
	}
	return cookie, nivelAcceso, http.StatusOK, err

}
func validarCookie(r *http.Request) (usuarioID int, nivelAcceso int, Username string, status int, err error) {

	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			return 0, 0, "", http.StatusUnauthorized, err
		} else {
			return 0, 0, "", http.StatusBadRequest, err
		}
	}

	tokenStr := cookie.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return 0, 0, "", http.StatusUnauthorized, err
		} else {
			return 0, 0, "", http.StatusBadRequest, err
		}
	}
	if !tkn.Valid {
		return 0, 0, "", http.StatusUnauthorized, err
	}

	// cedulaInt, err := strconv.Atoi(claims.UserID)
	// if err != nil {
	// 	return 0, 0, "", http.StatusBadRequest, err

	return claims.UserID, claims.NivelAcceso, claims.Username, http.StatusOK, err
}

func handlerSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var u Usuario
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	_, err = db.Exec("INSERT INTO Usuarios (NivelAcceso, Username, Password) VALUES (1,?,?)",
		u.Username, u.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	cookie, nivelAcceso, status, err := validarLogin(u.Username, u.Password)
	if cookie == nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	http.SetCookie(w, cookie)
	w.Write(fmt.Appendf(nil, "{\"nivelAcceso\":%d}", nivelAcceso))
}
func handlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var u Usuario
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	cookie, nivelAcceso, status, err := validarLogin(u.Username, u.Password)
	if cookie == nil {
		w.WriteHeader(status)
		w.Write([]byte(err.Error()))
		return
	}

	http.SetCookie(w, cookie)
	w.Write(fmt.Appendf(nil, "{\"nivelAcceso\":%d}", nivelAcceso))

}

func asociarHandlersAuth() {
	http.HandleFunc("/signup", handlerSignup)
	http.HandleFunc("/login", handlerLogin)
}
