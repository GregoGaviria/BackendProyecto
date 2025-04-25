package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
)

func handlerEliminarUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, nivelAcceso, err := authWrapper(r, w, 3)
	if err != nil {
		return
	}
	body := struct {
		UsuarioId int `json:"usuarioId"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	row := db.QueryRow(
		"SELECT NivelAcceso FROM Usuarios WHERE UsuarioId = ?",
		body.UsuarioId,
	)
	var nivelAccesoU int
	row.Scan(&nivelAccesoU)
	if nivelAccesoU >= 3 && nivelAcceso == 3 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("usuario debe tener un nivel de acceso de 4 para hacer esta accion"))
		return
	}
	_, err = db.Exec(
		"DELETE FROM Usuarios_has_Distritos WHERE Usuarios_UsuarioId = ?",
		body.UsuarioId,
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(
		"UPDATE Reporte SET UsuarioId = 1 WHERE UsuarioId = ?",
		body.UsuarioId,
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(
		"DELETE FROM Usuarios WHERE UsuarioId = ?",
		body.UsuarioId,
	)
	if err != nil {
		log.Fatal(err)
	}
}
func handlerEliminarUsuarioPropio(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioID, _, err := authWrapper(r, w, 0)
	_, err = db.Exec(
		"DELETE FROM Usuarios_has_Distritos WHERE Usuarios_UsuarioId = ?",
		usuarioID,
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(
		"UPDATE Reporte SET UsuarioId = 1 WHERE UsuarioId = ?",
		usuarioID,
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(
		"DELETE FROM Usuarios WHERE UsuarioId = ?",
		usuarioID,
	)
	if err != nil {
		log.Fatal(err)
	}
}
func handlerAsociarRegion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, _, err := authWrapper(r, w, 3)
	if err != nil {
		return
	}
	body := struct {
		UsuarioId          int `json:"usuarioId"`
		DistritoIdDeseado  int `json:"distritoIdDeseado"`
		CantonIdDeseado    int `json:"cantonIdDeseado"`
		ProvinciaIdDeseada int `json:"provinciaIdDeseada"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if body.ProvinciaIdDeseada != 0 && body.DistritoIdDeseado == 0 && body.CantonIdDeseado == 0 {
		rows, err := db.Query("SELECT DistritoId FROM ProvinciasCantonesView WHERE ProvinciaId = ?", body.ProvinciaIdDeseada)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		defer rows.Close()

		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				log.Fatal(err)
				return
			}
			_, err = db.Exec("INSERT Usuarios_has_Distritos VALUES (?,?)",
				body.UsuarioId, id,
			)
			var mysqlErr *mysql.MySQLError
			if err != nil {
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				} else {
					log.Fatal(err)
					return
				}
			}
		}

	} else if body.CantonIdDeseado != 0 && body.DistritoIdDeseado == 0 && body.ProvinciaIdDeseada == 0 {
		rows, err := db.Query("SELECT DistritoId FROM Distritos WHERE Cantonid = ?", body.CantonIdDeseado)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		defer rows.Close()

		for rows.Next() {
			var id int
			err := rows.Scan(&id)
			if err != nil {
				log.Fatal(err)
				return
			}
			_, err = db.Exec("INSERT Usuarios_has_Distritos VALUES (?,?)",
				body.UsuarioId, id,
			)
			var mysqlErr *mysql.MySQLError
			if err != nil {
				if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				} else {
					log.Fatal(err)
					return
				}
			}
		}

	} else if body.DistritoIdDeseado != 0 && body.CantonIdDeseado == 0 && body.ProvinciaIdDeseada == 0 {
		_, err = db.Exec("INSERT Usuarios_has_Distritos VALUES (?,?)",
			body.UsuarioId, body.DistritoIdDeseado,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Por favor solamente llenar canton, distrito o provincia"))
		return
	}

}
func handlerEliminarAsociacion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, _, err := authWrapper(r, w, 3)
	body := struct {
		UsuarioId   int `json:"usuarioId"`
		DistritoId  int `json:"distritoId"`
		CantonId    int `json:"cantonId"`
		ProvinciaId int `json:"provinciaId"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var res sql.Result
	if body.ProvinciaId != 0 && body.DistritoId == 0 && body.CantonId == 0 {
		res, err = db.Exec(
			"DELETE FROM Usuarios_has_Distritos "+
				"WHERE Usuarios_UsuarioId = ? AND Distritos_DistritoId LIKE ?",
			body.UsuarioId, fmt.Sprintf("%d____", body.ProvinciaId),
		)
	} else if body.CantonId != 0 && body.DistritoId == 0 && body.ProvinciaId == 0 {
		res, err = db.Exec(
			"DELETE FROM Usuarios_has_Distritos "+
				"WHERE Usuarios_UsuarioId = ? AND Distritos_DistritoId LIKE ?",
			body.UsuarioId, fmt.Sprintf("%d__", body.CantonId),
		)
	} else if body.DistritoId != 0 && body.CantonId == 0 && body.ProvinciaId == 0 {
		res, err = db.Exec(
			"DELETE FROM Usuarios_has_Distritos "+
				"WHERE Usuarios_UsuarioId = ? AND Distritos_DistritoId = ?",
			body.UsuarioId, body.DistritoId,
		)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Por favor solamente llenar canton, distrito o provincia"))
		return
	}
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	rowsaffected, err := res.RowsAffected()
	if err != nil {
		log.Fatal(err)
		return
	}
	if rowsaffected == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no se encontro ninguna asociacion con el id regional proveido, por favor revisar el id enviado e intentar otra vez"))
	}
}
func handlerEliminarAsociacionTodas(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, _, err := authWrapper(r, w, 3)
	body := struct {
		UsuarioId int `json:"usuarioId"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	_, err = db.Exec(
		"DELETE FROM Usuarios_has_Distritos WHERE Usuarios_UsuarioId = ?",
		body.UsuarioId,
	)
}
func handlerCambiarNivelAcceso(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, nivelAcceso, err := authWrapper(r, w, 3)
	if err != nil {
		return
	}
	body := struct {
		UsuarioId          int `json:"usuarioId"`
		NivelAccesoDeseado int `json:"nivelAccesoDeseado"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	if body.NivelAccesoDeseado > 4 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("no se puede asignar un nivel de acceso mayor a 4"))
		return
	} else if body.NivelAccesoDeseado >= 3 && nivelAcceso == 3 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("el usuario asignador debe tener un nivel de acceso de 4 para hacer esta accion"))
		return
	}
	var nivelAccesoUsuario int
	row := db.QueryRow("SELECT NivelAcceso FROM Usuarios WHERE UsuarioId = ?", body.UsuarioId)
	err = row.Scan(&nivelAccesoUsuario)
	if err != nil {
		log.Fatal(err)
		return
	}
	if nivelAccesoUsuario >= 3 && nivelAcceso == 3 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(
			"usuario con nivel de acceso 3 no puede cambiar" +
				" el nivel de un usuario con nivel de acceso mayor a 2",
		))
	}
	_, err = db.Exec(
		"UPDATE Usuarios SET NivelAcceso = ? WHERE UsuarioId = ?",
		body.NivelAccesoDeseado, body.UsuarioId,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

}
func handlerGetDistritosByUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	_, _, err := authWrapper(r, w, 3)
	if err != nil {
		return
	}
	usuarioID, err := strconv.Atoi(r.URL.Query().Get("usuario"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	distritos := querryWrapper[UsuarioDistrito](
		"SELECT * FROM UsuariosDistritosView WHERE UsuarioId = ?",
		usuarioID,
	)
	jsonWrapper(distritos, w)
}
func handlerGetDistritosPropios(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioID, _, err := authWrapper(r, w, 0)
	if err != nil {
		return
	}
	distritos := querryWrapper[UsuarioDistrito](
		"SELECT * FROM UsuariosDistritosView WHERE UsuarioId = ?",
		usuarioID,
	)
	jsonWrapper(distritos, w)
}
func handlerBuscarUsuarios(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioID := r.URL.Query().Get("usuario")
	_, _, err := authWrapper(r, w, 3)
	if err != nil {
		return
	}
	usuarios := querryWrapper[Usuario](
		"SELECT * FROM Usuarios WHERE Username LIKE ?",
		"%"+usuarioID+"%",
	)
	jsonWrapper(usuarios[1:], w)
}

func asociarHandlersUsuarios() {
	http.HandleFunc("/buscarUsuarios", handlerBuscarUsuarios)
	http.HandleFunc("/eliminarUsuario", handlerEliminarUsuario)
	http.HandleFunc("/eliminarUsuarioPropio", handlerEliminarUsuarioPropio)
	http.HandleFunc("/asociarRegion", handlerAsociarRegion)
	http.HandleFunc("/cambiarNivelAcceso", handlerCambiarNivelAcceso)
	http.HandleFunc("/getDistritosByUsuario", handlerGetDistritosByUsuario)
	http.HandleFunc("/getDistritosPropios", handlerGetDistritosPropios)
	http.HandleFunc("/eliminarAsociacion", handlerEliminarAsociacion)
	http.HandleFunc("/eliminarAsociacionTodas", handlerEliminarAsociacionTodas)
}
