package main

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB
var err error

type Usuario struct {
	ID     uint   `gorm:"primary_key"`
	Nombre string `gorm:"not null"`
	Correo string `gorm:"not null;unique"`
	Clave  string `gorm:"not null"`
	Status bool   `gorm:"default:true"`
}

func main() {
	dsn := "root:@tcp(127.0.0.1:3306)/practica?charset=utf8&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	} else {
		log.Println("Conexión a la base de datos exitosa")
	}

	err = db.AutoMigrate(&Usuario{})
	if err != nil {
		log.Fatal("Error al migrar la base de datos:", err)
	} else {
		log.Println("Migración de base de datos exitosa")
	}

	r := gin.Default()

	// Configurar sesión con cookie
	store := cookie.NewStore([]byte("mi-clave-secreta"))
	r.Use(sessions.Sessions("mysession", store))

	r.LoadHTMLGlob("templates/*")

	r.GET("/", showLoginPage)
	r.POST("/login", login)
	r.POST("/register", register)
	r.POST("/logout", logout)
	r.POST("/change-name", changeName)
	r.POST("/delete-account", deleteAccount)

	err = r.Run(":8080")
	if err != nil {
		log.Fatal("Error al iniciar el servidor:", err)
	} else {
		log.Println("Servidor iniciado en http://localhost:8080")
	}
}

func showLoginPage(c *gin.Context) {
	// Obtener la sesión del usuario
	session := sessions.Default(c)
	usuarioID := session.Get("usuario_id")

	var usuario Usuario
	if usuarioID != nil {
		if err := db.First(&usuario, usuarioID).Error; err == nil {
			c.HTML(http.StatusOK, "login.html", gin.H{
				"LoggedIn": true,
				"Usuario":  usuario.Nombre,
			})
			return
		}
	}

	// Si no está logueado
	c.HTML(http.StatusOK, "login.html", gin.H{
		"LoggedIn": false,
	})
}

func login(c *gin.Context) {
	correo := c.PostForm("correo")
	clave := c.PostForm("clave")

	var usuario Usuario
	if err := db.Where("correo = ?", correo).First(&usuario).Error; err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "Correo o contraseña incorrectos"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(usuario.Clave), []byte(clave)); err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "Correo o contraseña incorrectos"})
		return
	}

	// Guardar el usuario en la sesión
	session := sessions.Default(c)
	session.Set("usuario_id", usuario.ID)
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func register(c *gin.Context) {
	nombre := c.PostForm("nombre")
	correo := c.PostForm("correo")
	clave := c.PostForm("clave")

	var usuario Usuario
	if err := db.Where("correo = ?", correo).First(&usuario).Error; err == nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "El correo ya está registrado"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(clave), bcrypt.DefaultCost)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"Error": "Error al registrar el usuario"})
		return
	}

	newUser := Usuario{Nombre: nombre, Correo: correo, Clave: string(hashedPassword)}
	if err := db.Create(&newUser).Error; err != nil {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"Error": "Error al registrar el usuario"})
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{"Message": "Usuario registrado exitosamente"})
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("usuario_id") // Eliminar la sesión del usuario
	session.Save()

	c.Redirect(http.StatusFound, "/")
}

func changeName(c *gin.Context) {
	session := sessions.Default(c)
	usuarioID := session.Get("usuario_id")
	if usuarioID == nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "No estás logueado"})
		return
	}

	var usuario Usuario
	if err := db.First(&usuario, usuarioID).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"Error": "Error al obtener el usuario"})
		return
	}

	newName := c.PostForm("newName")
	usuario.Nombre = newName
	if err := db.Save(&usuario).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"Error": "Error al cambiar el nombre"})
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"Message":  "Nombre cambiado exitosamente",
		"LoggedIn": true,
		"Usuario":  usuario.Nombre,
	})
}

func deleteAccount(c *gin.Context) {
	session := sessions.Default(c)
	usuarioID := session.Get("usuario_id")
	if usuarioID == nil {
		c.HTML(http.StatusUnauthorized, "login.html", gin.H{"Error": "No estás logueado"})
		return
	}

	var usuario Usuario
	if err := db.First(&usuario, usuarioID).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"Error": "Error al obtener el usuario"})
		return
	}

	// Cambiar el campo Status a false (0) en lugar de eliminar la cuenta
	usuario.Status = false
	if err := db.Save(&usuario).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{"Error": "Error al desactivar la cuenta"})
		return
	}

	// Eliminar la sesión del usuario
	session.Delete("usuario_id")
	session.Save()

	c.HTML(http.StatusOK, "login.html", gin.H{"Message": "Cuenta desactivada exitosamente"})
}
