package main

import "net/http"

func handlerCrearUsuario(w http.ResponseWriter, r *http.Request) {

}
func handlerEliminarUsuario(w http.ResponseWriter, r *http.Request) {

}
func handlerAsociarDistrito(w http.ResponseWriter, r *http.Request) {

}
func handlerCambiarNivelAcesso(w http.ResponseWriter, r *http.Request) {

}

func asociarHandlersUsuarios() {
	http.HandleFunc("/crearUsuario", handlerCrearReporte)
	http.HandleFunc("/eliminarUsuario", handlerEliminarUsuario)
	http.HandleFunc("/asociarDistrito", handlerAsociarDistrito)
	http.HandleFunc("/cambiarNivelAcesso", handlerCambiarNivelAcesso)
}
