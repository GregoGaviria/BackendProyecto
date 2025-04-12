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

## Niveles de acceso
  * 1: usuario comun
    - puede generar reportes 
  * 2: empleado publico
    - puede tener distritos asignados y desactivar reportes
    - otorga acceso a /getReportesDistritosPropios, /getDistritosPropios, y /desactivar Reporte
  * 3: administrador
    - Puede promover un usuario comun a un empleado publico
    - Puede puede asignar distritos a usuarios
    - Puede eliminar usuarios
    - puede eliminar distritos a usuarios
    - puede ver datos de usuarios
    - otorga acceso a:
      - /eliminarUsuario
      - /asociarRegion
      - /eliminarAsociacion
      - /eliminarAsociacionTodas
      - /buscarUsuarios
      - /cambiarNivelAcceso
        - solamente puede subir un usuario a nivel 2, o bajar un usuario a nivel 1
      - /getDistritosByUsuario
      - /getReportesByUsuario
  * 4: super administrador
    - Puede promover un usuario a administrador o super administrador
  * Todos los niveles de acceso tienen acceso a las funciones de los niveles menores
## endpoints de api:

# endpoints de autenticacion
* /signup
  - POST
  - el request debe contener un json en el cuerpo que contenga `username` y `password`
  - retorna el nivel de acceso de usuario
* /login
  - POST
  - el request debe contener un json en el cuerpo que contenga `username` y `password`
  - retorna el nivel de acceso de usuario

# endpoints de reportes
* /crearReporte
  - POST
  - el request debe contener un json en el cuerpo que contenga `mensaje`, `tipoReporte`,  `CoordenadaX`, `CoordenadaY`
  - Latitud = Y Longitud = X
  - tipoReporte debe ser int
* /desactivarReporte
  - POST
  - el request debe contener un json en el cuerpo que contenga `reporteId`
* /getReportesByRegion
  - GET
  - el request debe contener el parametro de url `provinciaId`, `cantonId` o `distritoId`
* /getReporteById
  - GET
  - el request debe contener el parametro de url `id`
* /getReportesByUsuario
  - GET
  - el request debe contener el parametro de url `id`
* /getReportesPropios
  - GET
  - no se necesita ningun parametro
* /getReportesDistritosPropios
  - GET
  - no se necesita ningun parametro
  - devuelve una lista de reportes activos en los distritos que el usuario tiene asignado
  - solo para administradores

# endpoints de usuarios
* /buscarUsuarios
  - GET
  - el request debe contener el parametro de url `usuario`
* /eliminarUsuario
  - POST
  - el request debe contener un json en el cuerpo que contenga `usuarioId`
* /eliminarUsuarioPropio
  - POST
  - no requiere ningun parametro
* /asociarRegion
  - POST
  - el request debe contener un json en el cuerpo que contenga `usuarioId`, y `ProvinciaIdDeseada`, `CantonIdDeseado`,o `DistritoIdDeseado`
* /cambiarNivelAcceso
  - POST
  - el request debe contener un json en el cuerpo que contenga `usuarioId`, y `NivelAccesoDeseado`,
* /getDistritosByUsuario
  - GET
  - el request debe contener el parametro de url `usuario`
* /getDistritosPropios
  - GET
  - no necesita ningun parametro
* /eliminarAsociacion
  - elimina la asociacion entre un distrito y un usuario
  - POST
  - el request debe contener un json en el cuerpo que contenga `usuarioId`, y `ProvinciaId`, `CantonId`,o `DistritoId`
* /eliminarAsociacionTodas
  - POST
  - el request debe contener un json en el cuerpo que contenga `usuarioId`
