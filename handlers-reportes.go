package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"googlemaps.github.io/maps"
)

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

type LocationError struct {
	Msg                  string
	LocalizacionEsperada string
	LocalizacionRecibida string
}

func (e *LocationError) Error() string {
	return e.Msg
}

func reverseGeocodeRequest(y float64, x float64) (d int, c int, p int, r string, err error) {
	var resultadoDistrito *maps.GeocodingResult
	var resultadoCanton *maps.GeocodingResult
	var resultadoProvincia *maps.GeocodingResult
	var resultadoRoute *maps.GeocodingResult
	var resultadoCountry *maps.GeocodingResult

	req := &maps.GeocodingRequest{
		LatLng:   &maps.LatLng{Lat: y, Lng: x},
		Language: "es",
	}
	response, err := mapClient.ReverseGeocode(context.Background(), req)
	if err != nil {
		return 0, 0, 0, "", err
	}

	for i := range response {
		for _, f := range response[i].Types {
			switch f {
			case "country":
				resultadoCountry = &response[i]
				paisEsperado := "Costa Rica"
				if paisRecibido := resultadoCountry.FormattedAddress; paisRecibido != paisEsperado {
					return 0, 0, 0, "", &LocationError{
						Msg: fmt.Sprintf("Coordinadas recibidas fuera de localizacion esperada"+
							"\nlocalizacion esperada: '%s'"+
							"\nlocalizacion recibida: '%s'",
							paisEsperado, paisRecibido),
						LocalizacionEsperada: paisEsperado,
						LocalizacionRecibida: paisRecibido,
					}
				}
			case "route":
				resultadoRoute = &response[i]
			case "administrative_area_level_1":
				resultadoProvincia = &response[i]
			case "administrative_area_level_2":
				resultadoCanton = &response[i]
			case "administrative_area_level_3":
				resultadoDistrito = &response[i]

			}
		}
	}
	r = strings.Split(resultadoRoute.FormattedAddress, ",")[0]

	// log.Println("")
	// log.Println("Provincia")
	// log.Println(resultadoProvincia.FormattedAddress)
	// log.Println("")
	// log.Println("Canton")
	// log.Println(resultadoCanton.FormattedAddress)
	// log.Println("")
	// log.Println("Distrito")
	// log.Println(resultadoDistrito.FormattedAddress)
	// log.Println("")

	querryWrapper := func(q string, arg int, r *maps.GeocodingResult, splitint int) (int, error) {
		var rows *sql.Rows
		var err error
		if arg != 0 {
			rows, err = db.Query(q, arg)
		} else {
			rows, err = db.Query(q)
		}
		if err != nil {
			return 0, err
		}
		defer rows.Close()

		var nombres []string
		m := make(map[string]int)
		for rows.Next() {
			var sr string
			var intr int
			var basurero int
			err := rows.Scan(&intr, &sr)
			if err != nil {
				if err.Error() == "sql: expected 3 destination arguments in Scan, not 2" {
					err = rows.Scan(&intr, &sr, &basurero)
				} else {
					return 0, err
				}
			}

			sr, err = normalize(sr)
			if err != nil {
				return 0, err
			}

			m[sr] = intr
			nombres = append(nombres, sr)
		}

		resultSplit := strings.Split(r.FormattedAddress, ",")
		var nsplit string
		if splitint < 0 {
			nsplit, err = normalize(resultSplit[len(resultSplit)+splitint])
		} else {
			nsplit, err = normalize(resultSplit[splitint])
		}
		if err != nil {
			return 0, err
		}
		p := m[nsplit]
		if p != 0 {
			return p, nil
		}

		log.Printf("no se encontro un resultado exacto con '%s', ejecutando LCS", nsplit)
		resultadoProvinciaComparada := lcsCompare(nsplit, nombres)
		// p.NombreProvicia = resultadoProvinciaComparada
		p = m[resultadoProvinciaComparada]
		return p, nil

	}
	p, err = querryWrapper("SELECT * FROM Provincias", 0, resultadoProvincia, 0)
	if err != nil {
		return 0, 0, 0, "", err
	}
	c, err = querryWrapper("SELECT * FROM Cantones WHERE Provinciaid = ?", p, resultadoCanton, -2)
	if err != nil {
		return 0, 0, 0, "", err
	}
	d, err = querryWrapper("SELECT * FROM Distritos WHERE Cantonid = ?", c, resultadoDistrito, 0)
	if err != nil {
		return 0, 0, 0, "", err
	}
	return d, c, p, r, err

}
func normalize(s string) (string, error) {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return "", err
	}
	result = strings.ToLower(result)
	//aqui uno pone los prefijos que pueden interferir con el programa
	listaPrefijosDelMal := []string{"provincia de ", " "}
	for _, i := range listaPrefijosDelMal {
		result, _ = strings.CutPrefix(result, i)
	}
	return result, nil
}

// algoritmo de subsecuencia comun mas larga utilizado para lidear con las diferencias
// en los nombres de google maps y los datos de el TSE. Por ejemplo hay casos donde google maps retorna
// "provincia de alajuela" y en la base de datos es solamente "alajuela"
// esta funcion acepta un string comparador, y una lista de strings con la que se compara, y al utilizar
// el algorirmo de subsecuencia comun mas larga para comparar el comparador con cada string en la lista,
// se retrona el string de la lista con la que el comparador tenga la mayor similaritud.
// adicionalmente utiliza hilos para acelerar la ejecucuión del programa
func lcsCompare(comparador string, lista []string) string {
	var lcsRecursive func(s1 *string, s2 *string, l1 int, l2 int, result *int)
	lcsRecursive = func(s1 *string, s2 *string, l1 int, l2 int, result *int) {
		// return caso base
		if l1 == 0 || l2 == 0 {
			// log.Printf("終了 '%s' '%s' '%d' '%d'", *s1, *s2, l1, l2)
			*result = 0
			return
		}

		var r int
		if (*s1)[l1-1] == (*s2)[l2-1] {
			// log.Printf("正解 '%s' '%s' '%d' '%d'", *s1, *s2, l1, l2)
			lcsRecursive(s1, s2, l1-1, l2-1, &r)
			*result = 1 + r
		} else {
			// log.Printf("違い '%s' '%s' '%d' '%d'", *s1, *s2, l1, l2)
			var result1 int
			var result2 int
			var wg sync.WaitGroup
			wg.Add(2)
			go func() {
				defer wg.Done()
				lcsRecursive(s1, s2, l1-1, l2, &result1)
			}()
			go func() {
				defer wg.Done()
				lcsRecursive(s1, s2, l1, l2-1, &result2)
			}()
			wg.Wait()
			// lcsRecursive(s1, s2, l1-1, l2, &result1)
			// lcsRecursive(s1, s2, l1, l2-1, &result2)
			*result = max(result1, result2)
		}
	}

	var wg sync.WaitGroup
	type stringPuntaje struct {
		str     string
		puntaje int
	}
	listaPuntajes := []stringPuntaje{}
	for _, i := range lista {
		wg.Add(1)
		go func() {
			var resultado int
			defer wg.Done()
			lcsRecursive(&comparador, &i, len(comparador), len(i), &resultado)
			listaPuntajes = append(listaPuntajes, stringPuntaje{str: i, puntaje: resultado})
		}()
		wg.Wait()
	}
	var puntajeGanador int
	var stringGanador string
	for _, f := range listaPuntajes {
		println(f.puntaje, f.str)
		if f.puntaje > puntajeGanador {
			puntajeGanador = f.puntaje
			stringGanador = f.str
		} else if f.puntaje == puntajeGanador {
			lg := len(stringGanador)
			lf := len(f.str)
			if lg > lf {
				stringGanador = f.str
			} else if lg == lf {
				log.Printf("no c como paso esto!!! sg==%s, pg==%d, sf==%s, pf==%d",
					stringGanador, puntajeGanador, f.str, f.puntaje,
				)
			}
		}
	}

	return stringGanador
}

func handlerDesactivarReporte(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioId, nivelAcceso, _, s, err := validarCookie(r)
	if s != http.StatusOK {
		w.WriteHeader(s)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	} else if err != nil {
		log.Println("No se como paso esto (yupi!!!!)")
		log.Fatal(err)
		return
	}
	if nivelAcceso < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("usuario debe tener nivel de acceso 2 o mayor para accesar"))
		return
	}
	body := struct {
		ReporteId int `json:"reporteId"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	var distritoId int
	var exists int
	row := db.QueryRow("SELECT DistritoId from Reporte WHERE ReporteId = ?", body.ReporteId)
	err = row.Scan(&distritoId)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("no existe un reporte con este id"))
			return
		} else {
			log.Fatal(err)
			return
		}
	}
	row = db.QueryRow(
		"SELECT EXISTS(SELECT * FROM Usuarios_has_Distritos WHERE Distritos_DistritoId = ? AND Usuarios_UsuarioId = ?)",
		distritoId, usuarioId,
	)
	err = row.Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}
	if exists != 1 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("El usuario no tiene permiso para modificar reportes en este distrito"))
		return
	}
	_, err = db.Exec(
		"UPDATE Reporte SET Activo = 0 WHERE ReporteId = ?",
		body.ReporteId,
	)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

}
func handlerCrearReporte(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioID, nivelAcceso, _, s, err := validarCookie(r)
	if usuarioID == 0 {
		w.WriteHeader(s)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	}
	if nivelAcceso == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body := struct {
		Mensaje     string  `json:"mensaje"`
		TipoReporte int     `json:"tipoReporte"`
		CoordenadaX float64 `json:"CoordenadaX"`
		CoordenadaY float64 `json:"CoordenadaY"`
	}{}
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	d, _, _, route, err := reverseGeocodeRequest(body.CoordenadaY, body.CoordenadaX)
	if err != nil {
		locErr := &LocationError{}
		if errors.As(err, &locErr) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(locErr.Msg))
			return
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}
	var calleID int
	row := db.QueryRow("SELECT CalleId FROM Calles WHERE NombreCalle=? AND DistritoId=?",
		route, d)
	if err := row.Scan(&calleID); err != nil {
		if err == sql.ErrNoRows {
			a, err := db.Exec("INSERT INTO Calles (NombreCalle, DistritoId) VALUES (?,?)",
				route, d)
			if err != nil {
				log.Fatal(err)
			}
			lastinsert, err := a.LastInsertId()
			calleID = int(lastinsert)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
	}

	// log.Printf("INSERT INTO Reporte "+
	// 	"(Mensaje,CoordenadaX,CoordenadaY,TipoReporteID,UsuarioId,DistritoId,CalleId) "+
	// 	"VALUES ('%s','%f','%f','%d','%d','%d','%d')",
	// 	body.Mensaje, body.CoordenadaX, body.CoordenadaY, body.TipoReporte, usuarioID, d, calleID,
	// )
	_, err = db.Exec("INSERT INTO Reporte "+
		"(Mensaje,CoordenadaX,CoordenadaY,TipoReporteID,UsuarioId,DistritoId,CalleId) "+
		"VALUES (?,?,?,?,?,?,?)",
		body.Mensaje, body.CoordenadaX, body.CoordenadaY, body.TipoReporte, usuarioID, d, calleID,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return

	}
	return

}
func handlerGetReportesByRegion(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	distritoId := r.URL.Query().Get("distritoId")
	cantonId := r.URL.Query().Get("cantonId")
	provinciaId := r.URL.Query().Get("provinciaId")
	var res []Reporte

	if provinciaId != "" && distritoId == "" && cantonId == "" {
		id, err := strconv.Atoi(provinciaId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		rows, err := db.Query(
			"SELECT DistritoId FROM ProvinciasCantonesView WHERE ProvinciaId = ?",
			id,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()
		for rows.Next() {
			var idD int
			err := rows.Scan(&idD)
			if err != nil {
				log.Fatal(err)
				return
			}
			rowsi, err := db.Query(
				"SELECT * FROM Reporte WHERE DistritoId = ? AND Activo = 1",
				idD,
			)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer rowsi.Close()
			for rowsi.Next() {
				var re Reporte
				err = rowsi.Scan(
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
				res = append(res, re)
			}

		}
	} else if cantonId != "" && provinciaId == "" && distritoId == "" {
		id, err := strconv.Atoi(cantonId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		rows, err := db.Query(
			"SELECT DistritoId FROM Distritos WHERE CantonId = ?",
			id,
		)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		defer rows.Close()
		for rows.Next() {
			var idD int
			err := rows.Scan(&idD)
			if err != nil {
				log.Fatal(err)
				return
			}
			rowsi, err := db.Query(
				"SELECT * FROM Reporte WHERE DistritoId = ? AND Activo = 1",
				idD,
			)
			if err != nil {
				log.Fatal(err)
				return
			}
			defer rowsi.Close()
			for rowsi.Next() {
				var re Reporte
				err = rowsi.Scan(
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
				res = append(res, re)
			}

		}
	} else if distritoId != "" && provinciaId == "" && cantonId == "" {
		id, err := strconv.Atoi(distritoId)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		rows, err := db.Query(
			"SELECT * FROM Reporte WHERE DistritoId = ? AND Activo = 1",
			id,
		)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var re Reporte
			err = rows.Scan(
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
			res = append(res, re)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Por favor solamente llenar canton, distrito o provincia"))
		return

	}

	elJson, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write(elJson)
}
func handlerGetReporteById(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	reporteId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
		&re.DistritoId,
		&re.CalleId,
	)
	if err != nil {
		log.Fatal(err)
	}
	elJson, err := json.Marshal(re)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write(elJson)

}
func handlerGetReportesByUsuario(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	_, nivelAcceso, _, status, err := validarCookie(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	} else if err != nil {
		log.Println("No se como paso esto (yupi!!!!)")
		log.Fatal(err)
		return
	}
	if nivelAcceso < 2 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("usuario debe tener nivel de acceso 2 o mayor para accesar"))
		return
	}
	var res []Reporte
	rows, err := db.Query("SELECT * FROM Reporte WHERE UsuarioId = ?", usuarioId)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var re Reporte
		err = rows.Scan(
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
		res = append(res, re)
	}
	elJson, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write(elJson)

}
func handlerGetReportesPropios(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioId, _, _, status, err := validarCookie(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	} else if err != nil {
		log.Println("No se como paso esto (yupi!!!!)")
		log.Fatal(err)
		return
	}
	var res []Reporte
	rows, err := db.Query("SELECT * FROM Reporte WHERE UsuarioId = ?", usuarioId)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var re Reporte
		err = rows.Scan(
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
		res = append(res, re)
	}
	elJson, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write(elJson)
}
func handlerGetReportesDistritosPropios(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	usuarioID, _, _, status, err := validarCookie(r)
	if status != http.StatusOK {
		w.WriteHeader(status)
		if err != nil {
			w.Write([]byte(err.Error()))
			return
		}
	} else if err != nil {
		log.Println("No se como paso esto (yupi!!!!)")
		log.Fatal(err)
		return
	}
	var res []Reporte
	rows, err := db.Query("SELECT * FROM UsuariosDistritosView WHERE UsuarioId = ?", usuarioID)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var d UsuarioDistrito
		err = rows.Scan(&d.DistritoID, &d.Distrito, &d.Canton, &d.Provincia, &d.UsuarioId)
		if err != nil {
			log.Fatal(err)
			return
		}
		rowsi, err := db.Query(
			"SELECT * FROM Reporte WHERE DistritoId = ? AND Activo = 1",
			d.DistritoID,
		)
		if err != nil {
			log.Fatal(err)
			return
		}
		defer rowsi.Close()
		for rowsi.Next() {
			var re Reporte
			err = rowsi.Scan(
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
			res = append(res, re)
		}
	}

	elJson, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Write(elJson)

}

func asociarHandlersReportes() {
	http.HandleFunc("/desactivarReporte", handlerDesactivarReporte)
	http.HandleFunc("/crearReporte", handlerCrearReporte)
	http.HandleFunc("/getReportesByRegion", handlerGetReportesByRegion)
	http.HandleFunc("/getReporteById", handlerGetReporteById)
	http.HandleFunc("/getReportesByUsuario", handlerGetReportesByUsuario)
	http.HandleFunc("/getReportesPropios", handlerGetReportesPropios)
	http.HandleFunc("/getReportesDistritosPropios", handlerGetReportesDistritosPropios)
}
