package handlers

import (
	models "CRUD/Models"
	"CRUD/db"
	"encoding/json"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

var tmpl *template.Template

func SetTemplates(t *template.Template) {
	tmpl = t
}
func Tmpl() *template.Template {
	return tmpl
}

// Estructuras para manejar datos
type Genero struct {
	ID     int    `json:"ID"`
	Nombre string `json:"Nombre"`
}

type Editor struct {
	ID     int    `json:"ID"`
	Nombre string `json:"Nombre"`
}

type Desarrollador struct {
	ID     int    `json:"ID"`
	Nombre string `json:"Nombre"`
}

type Plataforma struct {
	ID     int    `json:"ID"`
	Nombre string `json:"Nombre"`
}
type ConteoDashboard struct {
	Videojuegos     int
	Desarrolladores int
	Editores        int
	Generos         int
	Plataformas     int
}

// /////////////////////
// GÉNEROS
// /////////////////////
func CrearGenero(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nombre := r.FormValue("nombre")
		if nombre == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "El nombre es obligatorio"})
			return
		}
		// Validar si ya existe el género
		var count int
		err := db.GetDB().QueryRow("SELECT COUNT(*) FROM generos WHERE LOWER(nombre) = LOWER($1)", nombre).Scan(&count)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error validando género existente"})
			return
		}
		if count > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]string{"error": "El género ya existe"})
			return
		}
		// Insertar género
		_, err = db.GetDB().Exec("INSERT INTO generos (nombre) VALUES ($1)", nombre)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "Error al insertar género"})
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"mensaje": "Género creado exitosamente"})
		return
	}
	tmpl.ExecuteTemplate(w, "genero_crear.html", nil)
}

func ListarGeneros(w http.ResponseWriter, r *http.Request) {
	generos, err := db.ObtenerGeneros()
	if err != nil {
		http.Error(w, "Error al obtener géneros: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(generos)
}

func EliminarGenero(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "Método no permitido"})
		return
	}

	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID no proporcionado"})
		return
	}

	dbConn := db.GetDB()

	_, err := dbConn.Exec("DELETE FROM videojuegos WHERE genero_id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al eliminar videojuegos del género: " + err.Error()})
		return
	}

	_, err = dbConn.Exec("DELETE FROM generos WHERE id = $1", id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al eliminar género: " + err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Género eliminado exitosamente"})
}

func GeneroListaHTML(w http.ResponseWriter, r *http.Request) {
	// Reemplaza "genero_lista.html" por el nombre exacto del archivo HTML
	err := tmpl.ExecuteTemplate(w, "genero_lista.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func VideojuegosPorGenero(w http.ResponseWriter, r *http.Request) {
	idGenero := r.URL.Query().Get("id")
	if idGenero == "" {
		http.Error(w, "ID de género no proporcionado", http.StatusBadRequest)
		return
	}

	query := `
		SELECT v.id, v.titulo, g.nombre AS genero_nombre, d.nombre AS desarrollador_nombre,
			   e.nombre AS editor_nombre, v.fecha_lanzamiento, v.precio
		FROM videojuegos v
		JOIN generos g ON v.genero_id = g.id
		JOIN desarrolladores d ON v.desarrollador_id = d.id
		JOIN editores e ON v.editor_id = e.id
		WHERE g.id = $1
	`
	rows, err := db.GetDB().Query(query, idGenero)
	if err != nil {
		http.Error(w, "Error consultando videojuegos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var juegos []map[string]interface{}
	for rows.Next() {
		var v struct {
			ID               int
			Titulo           string
			GeneroNombre     string
			Desarrollador    string
			Editor           string
			FechaLanzamiento time.Time
			Precio           float64
		}
		if err := rows.Scan(&v.ID, &v.Titulo, &v.GeneroNombre, &v.Desarrollador, &v.Editor, &v.FechaLanzamiento, &v.Precio); err != nil {
			http.Error(w, "Error leyendo videojuegos: "+err.Error(), http.StatusInternalServerError)
			return
		}

		juegos = append(juegos, map[string]interface{}{
			"id":                   v.ID,
			"titulo":               v.Titulo,
			"genero_nombre":        v.GeneroNombre,
			"desarrollador_nombre": v.Desarrollador,
			"editor_nombre":        v.Editor,
			"fecha_lanzamiento":    v.FechaLanzamiento,
			"precio":               v.Precio,
		})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(juegos)
}

// /////////////////////
// EDITORES
// /////////////////////
func CrearEditor(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nombre := r.FormValue("nombre")
		if nombre == "" {
			http.Error(w, "El nombre es obligatorio", http.StatusBadRequest)
			return
		}

		// Verificar si ya existe un editor con ese nombre
		var existe bool
		err := db.GetDB().QueryRow("SELECT EXISTS(SELECT 1 FROM editores WHERE nombre = $1)", nombre).Scan(&existe)
		if err != nil {
			http.Error(w, "Error al verificar editor: "+err.Error(), http.StatusInternalServerError)
			return
		}
		if existe {
			// Devolver JSON si es AJAX
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "El nombre del editor ya existe"})
			return
		}

		// Insertar si no existe
		_, err = db.GetDB().Exec("INSERT INTO editores (nombre) VALUES ($1)", nombre)
		if err != nil {
			http.Error(w, "Error al insertar editor: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"mensaje": "Editor creado exitosamente"})
		return
	}

	tmpl.ExecuteTemplate(w, "editor_crear.html", nil)
}

func ListarEditor(w http.ResponseWriter, r *http.Request) {
	editores, err := db.ObtenerEditores()
	if err != nil {
		http.Error(w, "Error al obtener editores: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(editores)
}

func EliminarEditor(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/editores", http.StatusSeeOther)
		return
	}

	id := r.FormValue("id")
	if id == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "ID no proporcionado"})
		return
	}

	dbConn := db.GetDB()

	// 1. Eliminar videojuegos relacionados
	_, err := dbConn.Exec("DELETE FROM videojuegos WHERE editor_id = $1", id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al eliminar videojuegos del editor: " + err.Error()})
		return
	}

	// 2. Eliminar el editor
	_, err = dbConn.Exec("DELETE FROM editores WHERE id = $1", id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al eliminar editor: " + err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Editor eliminado exitosamente"})
}

func EditorListaHTML(w http.ResponseWriter, r *http.Request) {
	// Reemplaza "genero_lista.html" por el nombre exacto del archivo HTML
	err := tmpl.ExecuteTemplate(w, "editor_lista.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

///////////////////////
// DESARROLLADORES
///////////////////////

func CrearDesarrollador(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nombre := r.FormValue("nombre")
		if nombre == "" {
			http.Error(w, "El nombre es obligatorio", http.StatusBadRequest)
			return
		}
		_, err := db.GetDB().Exec("INSERT INTO desarrolladores (nombre) VALUES ($1)", nombre)
		if err != nil {
			http.Error(w, "Error al insertar desarrollador: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/desarrollador/html", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "desarrollador_crear.html", nil)
}

func ListarDesarrollador(w http.ResponseWriter, r *http.Request) {
	desarrolladores, err := db.ObtenerDesarrollador()
	if err != nil {
		http.Error(w, "Error al obtener desarrollador: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(desarrolladores)
}

func EliminarDesarrollador(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/desarrolladores", http.StatusSeeOther)
		return
	}
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}
	dbConn := db.GetDB()
	// Eliminar los videojuegos que usan este desarrollador
	_, err := dbConn.Exec("DELETE FROM videojuegos WHERE desarrollador_id = $1", id)
	if err != nil {
		http.Error(w, "Error al eliminar videojuegos del desarrollador: "+err.Error(), http.StatusInternalServerError)
		return
	}
	// Eliminar el desarrollador
	_, err = dbConn.Exec("DELETE FROM desarrolladores WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error al eliminar desarrollador: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/desarrolladores", http.StatusSeeOther)
}
func DesarrolladorListaHTML(w http.ResponseWriter, r *http.Request) {
	// Reemplaza "genero_lista.html" por el nombre exacto del archivo HTML
	err := tmpl.ExecuteTemplate(w, "desarrollador_lista.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

///////////////////////
// PLATAFORMAS
///////////////////////

func CrearPlataforma(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nombre := r.FormValue("nombre")
		if nombre == "" {
			http.Error(w, "El nombre es obligatorio", http.StatusBadRequest)
			return
		}
		_, err := db.GetDB().Exec("INSERT INTO plataformas (nombre) VALUES ($1)", nombre)
		if err != nil {
			http.Error(w, "Error al insertar plataforma: "+err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/plataforma/html", http.StatusSeeOther)
		return
	}
	tmpl.ExecuteTemplate(w, "plataforma_crear.html", nil)
}
func ListarPlataforma(w http.ResponseWriter, r *http.Request) {
	plataformas, err := db.ObtenerPlataforma()
	if err != nil {
		http.Error(w, "Error al obtener plataformas: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(plataformas)
}
func EliminarPlataforma(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/plataformas", http.StatusSeeOther)
		return
	}
	id := r.FormValue("id")
	if id == "" {
		http.Error(w, "ID no proporcionado", http.StatusBadRequest)
		return
	}
	dbConn := db.GetDB()
	//Eliminar los videojuegos que usan esta plataforma
	_, err := dbConn.Exec("DELETE FROM videojuegos WHERE plataforma_id = $1", id)
	if err != nil {
		http.Error(w, "Error al eliminar videojuegos de la plataforma: "+err.Error(), http.StatusInternalServerError)
		return
	}
	//Eliminar la plataforma
	_, err = dbConn.Exec("DELETE FROM plataformas WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Error al eliminar plataforma: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/plataformas", http.StatusSeeOther)
}
func PlataformaListaHTML(w http.ResponseWriter, r *http.Request) {
	// Reemplaza "genero_lista.html" por el nombre exacto del archivo HTML
	err := tmpl.ExecuteTemplate(w, "plataforma_lista.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// /////////////////////
// VIDEOJUEGOS
// /////////////////////
func ListarVideojuegos(w http.ResponseWriter, r *http.Request) {
	generoIDStr := r.URL.Query().Get("genero_id")
	var videojuegos []models.Videojuego
	var err error
	if generoIDStr != "" {
		generoID, convErr := strconv.Atoi(generoIDStr)
		if convErr != nil {
			http.Error(w, "ID de género inválido", http.StatusBadRequest)
			return
		}
		videojuegos, err = db.ObtenerVideojuegosPorGenero(generoID)
	} else {
		videojuegos, err = db.ObtenerVideojuegos()
	}
	if err != nil {
		http.Error(w, "Error al obtener videojuegos: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(videojuegos)
}
func CrearVideojuego(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "videojuego_crear.html", nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	titulo := r.FormValue("titulo")
	descripcion := r.FormValue("descripcion")
	generoID := r.FormValue("genero_id")
	desarrolladorID := r.FormValue("desarrollador_id")
	editorID := r.FormValue("editor_id")
	fechaLanzamiento := r.FormValue("fecha_lanzamiento")
	precio := r.FormValue("precio")
	imagenURL := r.FormValue("imagen_url")
	if titulo == "" || descripcion == "" || generoID == "" || desarrolladorID == "" || editorID == "" || fechaLanzamiento == "" || precio == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Debe llenar todos los campos"})
		return
	}
	existe, err := db.ExisteVideojuegoPorTitulo(titulo)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error validando título"})
		return
	}
	if existe {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"error": "El título del videojuego ya existe"})
		return
	}
	_, err = db.GetDB().Exec(`
		INSERT INTO videojuegos (titulo, descripcion, genero_id, desarrollador_id, editor_id, fecha_lanzamiento, precio, imagen_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
	`, titulo, descripcion, generoID, desarrolladorID, editorID, fechaLanzamiento, precio, imagenURL)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Error al crear videojuego"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Videojuego creado exitosamente"})
}
func EliminarVideojuego(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	idStr := r.FormValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	err = db.EliminarVideojuego(id)
	if err != nil {
		http.Error(w, "Error al eliminar videojuego: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Videojuego eliminado correctamente"})
}
func VideojuegosListaHTML(w http.ResponseWriter, r *http.Request) {
	// Reemplaza "genero_lista.html" por el nombre exacto del archivo HTML
	err := tmpl.ExecuteTemplate(w, "videojuego_lista.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func EditarVideojuego(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}
	id, _ := strconv.Atoi(r.FormValue("id"))
	titulo := r.FormValue("titulo")
	descripcion := r.FormValue("descripcion")
	generoID, _ := strconv.Atoi(r.FormValue("genero_id"))
	desarrolladorID, _ := strconv.Atoi(r.FormValue("desarrollador_id"))
	editorID, _ := strconv.Atoi(r.FormValue("editor_id"))
	fechaLanzamiento := r.FormValue("fecha_lanzamiento")
	precio, _ := strconv.ParseFloat(r.FormValue("precio"), 64)
	imagenURL := r.FormValue("imagen_url")

	_, err := db.GetDB().Exec(`
		UPDATE videojuegos 
		SET titulo = $1, descripcion = $2, genero_id = $3, desarrollador_id = $4, 
			editor_id = $5, fecha_lanzamiento = $6, precio = $7, imagen_url = $8, updated_at = NOW()
		WHERE id = $9`,
		titulo, descripcion, generoID, desarrolladorID, editorID, fechaLanzamiento, precio, imagenURL, id)

	if err != nil {
		http.Error(w, "Error al actualizar videojuego: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"mensaje": "Videojuego actualizado correctamente"})
}
func ObtenerVideojuego(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}
	vj, err := db.ObtenerVideojuegoPorID(id)
	if err != nil {
		http.Error(w, "Error al obtener videojuego: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(vj)
}

//menu

func MenuHandler(w http.ResponseWriter, r *http.Request) {
	dbConn := db.GetDB()

	var data ConteoDashboard

	dbConn.QueryRow("SELECT COUNT(*) FROM videojuegos").Scan(&data.Videojuegos)
	dbConn.QueryRow("SELECT COUNT(*) FROM desarrolladores").Scan(&data.Desarrolladores)
	dbConn.QueryRow("SELECT COUNT(*) FROM editores").Scan(&data.Editores)
	dbConn.QueryRow("SELECT COUNT(*) FROM generos").Scan(&data.Generos)
	dbConn.QueryRow("SELECT COUNT(*) FROM plataformas").Scan(&data.Plataformas)

	err := tmpl.ExecuteTemplate(w, "menu.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
