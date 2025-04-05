package main

import (
	"database/sql"
	"flag"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"googlemaps.github.io/maps"
)

var db *sql.DB
var mapClient *maps.Client
var jwtKey []byte

func initDB(connString *string) error {
	var err error
	db, err = sql.Open("mysql", *connString)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func initMaps(apiKey *string) error {
	var err error
	mapClient, err = maps.NewClient(maps.WithAPIKey(*apiKey))
	return err

}

func main() {
	var connString *string = flag.String("c", "",
		`El string de conección para la conección de la base de datos
        revisar https://github.com/go-sql-driver/mysql?tab=readme-ov-file#dsn-data-source-name
        para mas detalles`,
	)
	var apiKey *string = flag.String("a", "", "el api key para google maps")
	var port *string = flag.String("p", ":8000", "el puerto donde va a correr el servidor")
	var jwt *string = flag.String("j", "", "la llave de autenticacion jwt")
	flag.Parse()
	err := initDB(connString)
	if err != nil {
		log.Fatal(err)
	}
	err = initMaps(apiKey)
	if err != nil {
		log.Fatal(err)
	}
	jwtKey = []byte(*jwt)
	asociarHandlersRegiones()
	asociarHandlersReportes()
	asociarHandlersUsuarios()
	asociarHandlersAuth()
	http.ListenAndServe(*port, nil)

}
