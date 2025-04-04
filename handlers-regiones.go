package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type Distrito struct {
	Id             int    `json:"id"`
	IdCanton       int    `json:"idCanton"`
	NombreDistrito string `json:"nombreDistrito"`
}
type Canton struct {
	Id           int    `json:"id"`
	IdProvincia  int    `json:"idProvincia"`
	NombreCanton string `json:"nombreCanton"`
}
type Provincia struct {
	Id             int    `json:"id"`
	NombreProvicia string `json:"nombreProvincia"`
}

func handlerGetDistritosByCanton(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var canton int
	canton, err := strconv.Atoi(r.URL.Query().Get("Canton"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var distritos []Distrito
	rows, err := db.Query("SELECT * FROM Distritos WHERE Cantonid = ?", canton)
	if err != nil {
		log.Fatal(err)
	}

	defer rows.Close()

	for rows.Next() {
		var d Distrito
		err := rows.Scan(&d.Id, &d.NombreDistrito, &d.IdCanton)
		if err != nil {
			log.Fatal(err)
		}
		distritos = append(distritos, d)
	}
	elJson, err := json.Marshal(distritos)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(elJson)

}
func handlerGetCantonesByProvincia(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	provincia, err := strconv.Atoi(r.URL.Query().Get("Provincia"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var cantones []Canton
	rows, err := db.Query("SELECT * FROM Cantones WHERE Provinciaid = ?", provincia)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var c Canton
		err := rows.Scan(&c.Id, &c.NombreCanton, &c.IdProvincia)
		if err != nil {
			log.Fatal(err)
		}
		cantones = append(cantones, c)
	}
	elJson, err := json.Marshal(cantones)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(elJson)

}
func handlerGetProvincias(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	var provincias []Provincia
	rows, err := db.Query("SELECT * FROM Provincias")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var p Provincia
		err := rows.Scan(&p.Id, &p.NombreProvicia)
		if err != nil {
			log.Fatal(err)
		}
		provincias = append(provincias, p)
	}
	elJson, err := json.Marshal(provincias)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(elJson)

}
func handlerGetDistritoById(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	distritoId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var d Distrito
	row := db.QueryRow("SELECT * FROM Distritos WHERE DistritoId = ?", distritoId)
	err = row.Scan(&d.Id, &d.NombreDistrito, &d.IdCanton)
	if err != nil {
		log.Fatal(err)
	}
	elJson, err := json.Marshal(d)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(elJson)
}
func handlerGetCantonById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cantonId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var c Canton
	row := db.QueryRow("SELECT * FROM Cantones WHERE CantonId = ?", cantonId)
	err = row.Scan(&c.Id, &c.NombreCanton, &c.IdProvincia)
	if err != nil {
		log.Fatal(err)
	}
	elJson, err := json.Marshal(c)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(elJson)

}
func handlerGetProvinciaById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	provinciaId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var p Provincia
	row := db.QueryRow("SELECT * FROM Provincias WHERE ProvinciaId = ?", provinciaId)
	err = row.Scan(&p.Id, &p.NombreProvicia)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			log.Fatal(err)
		}
	}
	elJson, err := json.Marshal(p)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(elJson)

}
func handlerGetCantonesByUsuario(w http.ResponseWriter, r *http.Request) {

}

func asociarHandlersRegiones() {
	http.HandleFunc("/getDistritosByCanton", handlerGetDistritosByCanton)
	http.HandleFunc("/getCantonesByProvincia", handlerGetCantonesByProvincia)
	http.HandleFunc("/getProvincias", handlerGetProvincias)
	http.HandleFunc("/getDistritoById", handlerGetDistritoById)
	http.HandleFunc("/getCantonById", handlerGetCantonById)
	http.HandleFunc("/getProvinciaById", handlerGetProvinciaById)
	http.HandleFunc("/getCantonesByUsuario", handlerGetCantonesByUsuario)
}
