package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"
)

func handlerGetDistritosByCanton(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//consigue el valor de la variable en el url
	//si no se encuentra retorna badrequest
	canton, err := strconv.Atoi(r.URL.Query().Get("Canton"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//utiliza el querryWrapper para devolver los datos
	distritos := querryWrapper[Distrito](
		"SELECT * FROM Distritos WHERE Cantonid = ?",
		canton,
	)
	//utiliza el jsonWrapper para enviar los datos
	jsonWrapper(distritos, w)
}
func handlerGetCantonesByProvincia(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	provincia, err := strconv.Atoi(r.URL.Query().Get("Provincia"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cantones := querryWrapper[Canton](
		"SELECT * FROM Cantones WHERE Provinciaid = ?",
		provincia,
	)
	jsonWrapper(cantones, w)
}
func handlerGetProvincias(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
	provincias := querryWrapper[Provincia]("SELECT * FROM Provincias")
	jsonWrapper(provincias, w)
}
func handlerGetDistritoById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	distritoId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var d Distrito
	row := db.QueryRow("SELECT * FROM Distritos WHERE DistritoId = ?", distritoId)
	err = row.Scan(&d.Id, &d.NombreDistrito, &d.IdCanton)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no existe un distrito con este id"))
			return
		} else {
			log.Fatal(err)
			return
		}
	}
	jsonWrapper(d, w)
}
func handlerGetCantonById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	cantonId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var c Canton
	row := db.QueryRow("SELECT * FROM Cantones WHERE CantonId = ?", cantonId)
	err = row.Scan(&c.Id, &c.NombreCanton, &c.IdProvincia)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no existe un canton con este id"))
			return
		} else {
			log.Fatal(err)
			return
		}
	}
	jsonWrapper(c, w)

}
func handlerGetProvinciaById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	provinciaId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var p Provincia
	row := db.QueryRow("SELECT * FROM Provincias WHERE ProvinciaId = ?", provinciaId)
	err = row.Scan(&p.Id, &p.NombreProvicia)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no existe una provincia con este id"))
			return
		} else {
			log.Fatal(err)
			return
		}
	}
	jsonWrapper(p, w)
}

func asociarHandlersRegiones() {
	http.HandleFunc("/getDistritosByCanton", handlerGetDistritosByCanton)
	http.HandleFunc("/getCantonesByProvincia", handlerGetCantonesByProvincia)
	http.HandleFunc("/getProvincias", handlerGetProvincias)
	http.HandleFunc("/getDistritoById", handlerGetDistritoById)
	http.HandleFunc("/getCantonById", handlerGetCantonById)
	http.HandleFunc("/getProvinciaById", handlerGetProvinciaById)
}
