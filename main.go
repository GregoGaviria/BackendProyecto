package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type Ejemplo struct {
	Dato1 string `json:"dato1"`
	Dato2 int    `json:"dato2"`
}

var datosEjemplo = []Ejemplo{
	{Dato1: "aaa", Dato2: 111},
	{Dato1: "bbb", Dato2: 222},
	{Dato1: "ccc", Dato2: 333},
}


func handlerEjemplo(w http.ResponseWriter, r *http.Request) {
	elJson, err := json.Marshal(datosEjemplo)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(elJson)

}

func initDB(connString *string) error {
	var err error
	db, err = sql.Open("mysql", *connString)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
func main() {
	fmt.Println("a")
	var connString *string = flag.String("c",
		"servidor:clave@unix(/var/run/mysqld/mysqld.sock)/webdev",
		`El string de conección para la conección de la base de datos
        revisar https://github.com/go-sql-driver/mysql?tab=readme-ov-file#dsn-data-source-name
        para mas detalles`,
	)
	var port *string = flag.String("p", ":8000", "el puerto donde va a correr el servidor")
	flag.Parse()
	err := initDB(connString)
	if err != nil {
		log.Fatal(err)
	}
	http.HandleFunc("/ejemplo", handlerEjemplo)
	http.ListenAndServe(*port, nil)

}
