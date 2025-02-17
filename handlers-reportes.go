package main

import (
	"log"
	"net/http"
	"strconv"
)

type Reporte struct {
	Id          int     `json:"id"`
	Mensaje     string  `json:"mensaje"`
	CoordenadaX float32 `json:"coordenadaX"`
	CoordenadaY float32 `json:"coordenadaY"`
	Activo      bool    `json:"activo"`
	TipoReporte int     `json:"tipoReporte"`
	UsuarioId   int     `json:"usuarioId"`
	CalleId     int     `json:"calleId"`
	DistritoId  int     `json:"distritoId"`
}

func handlerCrearReporte(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
func handlerGetReportesByDistrito(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	distrito, err := strconv.Atoi(r.URL.Query().Get("Distrito"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	log.Println(distrito)
}
func handlerGetReporteById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	reporteId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	var re Reporte
	row := db.QueryRow("SELECT * FROM Reporte WHERE ReporteId = ?", reporteId)
	err = row.Scan(
		&re.Id,
		&re.Mensaje,
		&re.CoordenadaX,
		&re.CoordenadaY,
		&re.Activo,
		&re.TipoReporte,
		&re.UsuarioId,
		&re.CalleId,
		&re.DistritoId,
	)
	if err != nil {
		log.Fatal(err)
	}

}

func asociarHandlersReportes() {
	http.HandleFunc("/crearReporte", handlerCrearReporte)
	http.HandleFunc("/getReportesByDistrito", handlerGetReportesByDistrito)
	http.HandleFunc("/getReporteById", handlerGetReporteById)
}
