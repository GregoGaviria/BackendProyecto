package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// se crea una interfaz para las estructuras que contienen scan
type scanInterface interface {
	Scan(dest ...any) error
}

// se crea una interfaz para las estructuras que contengan populate
// debido a la manera que los metodos y los punteros estan integrados en golang
// se vuelve necesario integrar genericos para asegurarse que sirva
type populateInterface[T any] interface {
	populate(s scanInterface)
	*T
}

type Usuario struct {
	UserID      int
	NivelAcceso int
	Username    string `json:"username"`
	Password    string `json:"password"`
}

// funcion para popular los datos de la funcion al recibir una estructura que implemente Scan(dest ...any) error
func (u *Usuario) populate(s scanInterface) {
	var basurero string
	err := s.Scan(
		&u.UserID,
		&u.NivelAcceso,
		&u.Username,
		&basurero,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type UsuarioDistrito struct {
	DistritoID int    `json:"distritoID"`
	Distrito   string `json:"distrito"`
	Canton     string `json:"canton"`
	Provincia  string `json:"provincia"`
	UsuarioId  int    `json:"usuarioId"`
}

func (d *UsuarioDistrito) populate(s scanInterface) {
	err := s.Scan(
		&d.DistritoID,
		&d.Distrito,
		&d.Canton,
		&d.Provincia,
		&d.UsuarioId,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}

type Reporte struct {
	Id          int     `json:"id"`
	Mensaje     string  `json:"mensaje"`
	CoordenadaX float32 `json:"coordenadaX"`
	CoordenadaY float32 `json:"coordenadaY"`
	Activo      bool    `json:"activo"`
	TipoReporte int     `json:"tipoReporte"`
	UsuarioId   int     `json:"usuarioId"`
	DistritoId  int     `json:"distritoId"`
	CalleId     int     `json:"calleId"`
}

func (re *Reporte) populate(s scanInterface) {
	err := s.Scan(
		&re.Id,
		&re.Mensaje,
		&re.CoordenadaX,
		&re.CoordenadaY,
		&re.Activo,
		&re.TipoReporte,
		&re.UsuarioId,
		&re.DistritoId,
		&re.CalleId,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}

type Distrito struct {
	Id             int    `json:"id"`
	IdCanton       int    `json:"idCanton"`
	NombreDistrito string `json:"nombreDistrito"`
}

func (d *Distrito) populate(s scanInterface) {
	err := s.Scan(
		&d.Id,
		&d.NombreDistrito,
		&d.IdCanton,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
}

type Canton struct {
	Id           int    `json:"id"`
	IdProvincia  int    `json:"idProvincia"`
	NombreCanton string `json:"nombreCanton"`
}

func (c *Canton) populate(s scanInterface) {
	err := s.Scan(&c.Id, &c.NombreCanton, &c.IdProvincia)
	if err != nil {
		log.Fatal()
		return
	}
}

type Provincia struct {
	Id             int    `json:"id"`
	NombreProvicia string `json:"nombreProvincia"`
}

func (p *Provincia) populate(s scanInterface) {
	err := s.Scan(&p.Id, &p.NombreProvicia)
	if err != nil {
		log.Fatal()
		return
	}
}

// Funcion Generica para reducir codigo reutilizado al ejecutar queries de sql donde se retornan varias rows
func querryWrapper[T any, t populateInterface[T]](querry string, args ...any) []t {
	//se ejecuta el querry con todos los argumentos
	rows, err := db.Query(querry, args...)
	if err != nil {
		log.Fatal(err)
	}
	//se cierran los rows al terminar de correr la funcion
	defer rows.Close()
	//se genera una lista de el tipo generico donde se contienen los resultados
	var l []t
	for rows.Next() {
		//magia negra de punteros para crear una variable i de el tipo generico
		var j T
		i := t(&j)
		//se utliza el metodo de populate para popular la variable con los datos del row
		i.populate(rows)
		l = append(l, i)
	}
	return l

}

// funcion para reducir codigo repetido en los requests donde se necesita enviar un json
func jsonWrapper(v any, w http.ResponseWriter) {
	elJson, err := json.Marshal(v)
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(elJson)
}
