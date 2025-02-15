package main

import (
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

func handlerGetDistritosByCanton(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// if r.Body==""{
	//     w.WriteHeader(http.StatusBadRequest)
	// }
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

// type Ejemplo struct {
// 	Dato1 string `json:"dato1"`
// 	Dato2 int    `json:"dato2"`
// }
//
// var datosEjemplo = []Ejemplo{
// 	{Dato1: "aaa", Dato2: 111},
// 	{Dato1: "bbb", Dato2: 222},
// 	{Dato1: "ccc", Dato2: 333},
// }
// func getNombresConsultorios() ([]string, error) {
// 	rows, err := db.Query("SELECT nombre_consultorio FROM consultorio")
// 	var resultado []string
// 	if err != nil {
// 		log.Println(err)
// 		return resultado, err
// 	}
// 	defer rows.Close()
// 	for rows.Next() {
// 		var s string
// 		rows.Scan(&s)
// 		resultado = append(resultado, s)
// 	}
// 	return resultado, nil
// }
//
// func handlerEjemplo(w http.ResponseWriter, r *http.Request) {
// 	elJson, err := json.Marshal(datosEjemplo)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	w.Write(elJson)
//
// }
//
//
