# Este es el backend del proyecto, es un REST API escrito en GOLANG con una base de datos mariadb/mysql.


## guia para correr el proyecto

・Instalar [go](https://go.dev/dl/) y [mysql](https://dev.mysql.com/downloads/workbench/)

・Descargar el proyecto y colocar los contenidos en el directorio deseado

・Correr el script dentro de el archivo de dbscripts con el nombre de elScript.sql dentro de mysql

・verificar que la base de datos este corriendo correctamente

・abrir una ventana de cmd en el direcorio donde esta main.go

・escribir el siguiente comando:

    go run . -a [La llave del api] -c [Usuario de db]:[Contraseña de db]@tcp([DireccionIP:Puerto])/[baseDeDatos] -j [llave De JWT]

・se puede usar 127.0.0.1 para la direccion ip y el puerto default es 3306
・la llave jwt puede ser cualquier string 

## endpoints de api:

# endpoints de autenticacion
/signup
    POST
    el request debe contener un json en el cuerpo que contenga `username` y `password`
/login
    POST
    el request debe contener un json en el cuerpo que contenga `username` y `password`

# endpoints de reportes
/crearReporte
    POST
    el request debe contener un json en el cuerpo que contenga `mensaje`, `tipoReporte`,  `CoordenadaX`, `CoordenadaY`
    Latitud = Y Longitud = X

