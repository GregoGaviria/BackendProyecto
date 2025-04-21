package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	UserID      int    `json:"userID"`
	NivelAcceso int    `json:"nivelAcesso"`
	Username    string `json:"username"`
	jwt.StandardClaims
}

// funcion para validar el login al recibir un nombre de usuario y una clave
// retorna un cookie que contiene el token de autorizacion
func validarLogin(username string, password string) (cookie *http.Cookie, nivelAcceso int, status int, err error) {

	if username == "" || password == "" {
		return nil, 0, http.StatusBadRequest, nil
	}
	//la contraseña para el el usuario eliminado es * entonces esto deniega el login para
	//los usuarios intentando entrar a el usuario eliminado
	if password == "*" {
		return nil, 0, http.StatusBadRequest, nil
	}

	var userID int
	var expectedPassword string
	row := db.QueryRow("SELECT UsuarioId, NivelAcceso, Password FROM Usuarios WHERE Username=?", username)
	err = row.Scan(&userID, &nivelAcceso, &expectedPassword)
	//retorna si la contraseña no es la esperada o si el nombre de usuario no existe
	if err != nil || password != expectedPassword {
		return nil, 0, http.StatusUnauthorized, err
	}
	//genera el tiempo de expitacion que se le asigna al token
	//puede ser cambiado dependiendo de las necesidades de la aplicacion
	expirationTime := time.Now().Add(time.Minute * 60)
	//genera un puntero a una estructura que contiene los claims contenidos en el token
	var claims *Claims = &Claims{
		UserID:      userID,
		NivelAcceso: nivelAcceso,
		Username:    username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	//encodifica el token a un string con el jwtKey especificado en el -j utilizado al ejecutar el programa
	var token *jwt.Token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenstring, err := token.SignedString(jwtKey)
	if err != nil {
		return nil, 0, http.StatusInternalServerError, err
	}
	//genera un cookie con el token
	cookie = &http.Cookie{
		Name:    "token",
		Value:   tokenstring,
		Expires: expirationTime,
	}
	return cookie, nivelAcceso, http.StatusOK, err

}

// funcion que puede ser llamada en un handler para las siguientes funciones
// ・asegurar que el usuario este autenticado
// ・verificar que el nivel de acceso sea el correcto para la funcion
// ・recibir el id del usuario para integrarlo en los queries de sql
func authWrapper(r *http.Request, w http.ResponseWriter, nivelAuthDeseado int) (usuarioID int, nivelAcceso int, err error) {
	//codigo para ingresar el valor de el token en una variable y resoponde con un error si no la encuentra
	cookie, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return 0, 0, errors.New("unauthorized")
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return 0, 0, errors.New("bad request")
		}
	}

	//asigna el valor del cookie en la variable
	tokenStr := cookie.Value

	//se genera un puntero a un objeto de Claims nuevo
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tokenStr, claims,
		func(t *jwt.Token) (any, error) {
			return jwtKey, nil
		})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return 0, 0, err
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return 0, 0, err
		}
	}
	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		err = errors.New("invalid token")
		w.Write([]byte(err.Error()))
		return 0, 0, err
	}
	if claims.NivelAcceso < nivelAuthDeseado {
		w.WriteHeader(http.StatusUnauthorized)
		err = errors.New("usuario debe tener nivel de acceso " + strconv.Itoa(nivelAuthDeseado) +
			" o mayor para accesar\n Valor actual : " + strconv.Itoa(claims.NivelAcceso),
		)
		w.Write([]byte(err.Error()))
		return 0, 0, err
	}

	return claims.UserID, claims.NivelAcceso, err
}

func handlerSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var u Usuario
	//llena el usuario u con los datos del cuerpo
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

	//aplica el cookie en el request
	http.SetCookie(w, cookie)
	//escribe en el cuerpo un json que contiene el nivel de acceso
	w.Write(fmt.Appendf(nil, "{\"nivelAcceso\":%d}", nivelAcceso))
}
func handlerLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var u Usuario
	//llena el usuario u con los datos del cuerpo
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

	//aplica el cookie en el result
	http.SetCookie(w, cookie)
	//escribe en el cuerpo un json que contiene el nivel de acceso
	w.Write(fmt.Appendf(nil, "{\"nivelAcceso\":%d}", nivelAcceso))
}

func asociarHandlersAuth() {
	http.HandleFunc("/signup", handlerSignup)
	http.HandleFunc("/login", handlerLogin)
}
