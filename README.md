# Proyecto Go con Gin y MySQL

Este es un proyecto básico en Go que utiliza el framework Gin para el manejo de rutas y sesiones, junto con MySQL para la base de datos. A continuación se detallan los pasos necesarios para instalar y ejecutar este proyecto en tu entorno local.
no se si descargué unas librerias de mas porque no me funcionaba TT

## Requisitos

Antes de ejecutar el proyecto, asegúrate de tener lo siguiente instalado:

1. **Go**: Debes tener Go instalado en tu sistema. Puedes descargarlo e instalarlo desde el siguiente enlace:
   - [Descargar Go](https://golang.org/dl/)

2. **Gorm**: La biblioteca de ORM para Go que se utiliza para interactuar con MySQL.
   
3. **Gin**: El framework web que utilizamos para las rutas y sesiones.

### Instalar dependencias

Una vez tengas Go y MySQL instalados, puedes instalar las dependencias necesarias ejecutando los siguientes comandos en tu terminal dentro del directorio de tu proyecto.

```bash
go mod init <nombre_del_proyecto>
go get github.com/gin-contrib/sessions
go get github.com/gin-contrib/sessions/cookie
go get github.com/gin-gonic/gin
go get golang.org/x/crypto/bcrypt
go get gorm.io/driver/mysql
go get gorm.io/gorm
