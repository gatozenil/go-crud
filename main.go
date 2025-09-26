package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/gatozenil/go-crud/db"

	handlers "github.com/gatozenil/go-crud/handlers"
)

func main() {
	// Cargar htmls
	http.Handle("/templates/", http.StripPrefix("/templates/", http.FileServer(http.Dir("templates"))))
	// Cargar CSS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	// Cargar templates
	t := template.Must(template.ParseGlob("templates/*.html"))
	handlers.SetTemplates(t)
	// Conecta a la base de datos
	db.Init()

	// Rutas
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/templates/login.html", http.StatusSeeOther)
	})
	http.HandleFunc("/menu", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session")
		if err != nil || cookie.Value != "authenticated" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		// Llamar al handler que muestra el menú con datos dinámicos
		handlers.MenuHandler(w, r)
	})

	// Rutas
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/genero/crear", handlers.CrearGenero)
	http.HandleFunc("/logout", loginHandler)
	http.HandleFunc("/generos", handlers.ListarGeneros)
	http.HandleFunc("/generos/crear", handlers.CrearGenero)
	http.HandleFunc("/genero/eliminar", handlers.EliminarGenero)
	http.HandleFunc("/generos/html", handlers.GeneroListaHTML)
	http.HandleFunc("/editores/crear", handlers.CrearEditor)
	http.HandleFunc("/editores", handlers.ListarEditor)
	http.HandleFunc("/editores/eliminar", handlers.EliminarEditor)
	http.HandleFunc("/editor/html", handlers.EditorListaHTML)
	http.HandleFunc("/Desarrollador/crear", handlers.CrearDesarrollador)
	http.HandleFunc("/Desarrollador", handlers.ListarDesarrollador)
	http.HandleFunc("/Desarrollador/eliminar", handlers.EliminarDesarrollador)
	http.HandleFunc("/desarrollador/html", handlers.DesarrolladorListaHTML)
	http.HandleFunc("/Plataforma/crear", handlers.CrearPlataforma)
	http.HandleFunc("/Plataforma", handlers.ListarPlataforma)
	http.HandleFunc("/Plataforma/eliminar", handlers.EliminarPlataforma)
	http.HandleFunc("/plataforma/html", handlers.PlataformaListaHTML)
	http.HandleFunc("/videojuegos/genero", handlers.VideojuegosPorGenero)
	http.HandleFunc("/videojuego", handlers.ObtenerVideojuego)
	http.HandleFunc("/videojuegos/editar", handlers.EditarVideojuego)
	http.HandleFunc("/videojuegos", handlers.ListarVideojuegos)
	http.HandleFunc("/videojuegos/crear", handlers.CrearVideojuego)
	http.HandleFunc("/videojuegos/eliminar", handlers.EliminarVideojuego)
	http.HandleFunc("/videojuegos/html", handlers.VideojuegosListaHTML)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println("Servidor iniciado en http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

// Función para manejar el login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "login.html", http.StatusSeeOther)
		return
	}

	user := r.FormValue("username")
	pass := r.FormValue("password")

	var storedPassword string
	err := db.GetDB().QueryRow("SELECT password FROM users WHERE username=$1", user).Scan(&storedPassword)

	if err == sql.ErrNoRows {
		http.Redirect(w, r, "/templates/login.html?msg=Usuario+no+encontrado", http.StatusSeeOther)
		return
	} else if err != nil {
		http.Redirect(w, r, "/templates/login.html?msg=Error+interno", http.StatusSeeOther)
		return
	}

	if pass != storedPassword {
		http.Redirect(w, r, "/templates/login.html?msg=Contraseña+incorrecta", http.StatusSeeOther)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "authenticated",
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
	})

	// Redirigir al menú
	http.Redirect(w, r, "/menu", http.StatusSeeOther)
}
