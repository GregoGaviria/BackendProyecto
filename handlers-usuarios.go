package main

import (
	"net/http"
)

type Usuario struct {
	UserID      int
	NivelAcceso int
	Username    string `json:"username"`
	Password    string `json:"password"`
}

func handlerEliminarUsuario(w http.ResponseWriter, r *http.Request) {

}
func handlerAsociarDistrito(w http.ResponseWriter, r *http.Request) {

}
func handlerCambiarNivelAcesso(w http.ResponseWriter, r *http.Request) {

}

func asociarHandlersUsuarios() {
	http.HandleFunc("/eliminarUsuario", handlerEliminarUsuario)
	http.HandleFunc("/asociarDistrito", handlerAsociarDistrito)
	http.HandleFunc("/cambiarNivelAcesso", handlerCambiarNivelAcesso)
}
