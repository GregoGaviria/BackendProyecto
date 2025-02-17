package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

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
		"servidor:clave@unix(/var/run/mysqld/mysqld.sock)/mydb",
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
	http.HandleFunc("/getDistritosByCanton", handlerGetDistritosByCanton)
	http.ListenAndServe(*port, nil)

}
