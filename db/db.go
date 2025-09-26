package db

import (
	models "github.com/gatozenil/go-crud/Models"

	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("Falta variable de entorno DATABASE_URL")
	}
	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Error conectando a BD:", err)
	}
	if err = DB.Ping(); err != nil {
		log.Fatal("No se pudo conectar a la BD:", err)
	}
}
func GetDB() *sql.DB {
	return DB
}

// Funciones para acceder a datos
// Funcion generos
func ObtenerGeneros() ([]models.Genero, error) {
	rows, err := GetDB().Query("SELECT id, nombre FROM generos")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var generos []models.Genero
	for rows.Next() {
		var g models.Genero
		if err := rows.Scan(&g.ID, &g.Nombre); err != nil {
			return nil, err
		}
		generos = append(generos, g)
	}
	return generos, nil
}
func ExisteGeneroPorNombre(nombre string) (bool, error) {
	var count int
	err := GetDB().QueryRow("SELECT COUNT(*) FROM generos WHERE nombre = $1", nombre).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Funcion editores
func ObtenerEditores() ([]models.Editor, error) {
	rows, err := GetDB().Query("SELECT id, nombre FROM editores")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Editores []models.Editor
	for rows.Next() {
		var g models.Editor
		if err := rows.Scan(&g.ID, &g.Nombre); err != nil {
			return nil, err
		}
		Editores = append(Editores, g)
	}
	return Editores, nil
}

// Funcion desarrollarores
func ObtenerDesarrollador() ([]models.Desarrollador, error) {
	rows, err := GetDB().Query("SELECT id, nombre FROM desarrolladores")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Desarrollador []models.Desarrollador
	for rows.Next() {
		var g models.Desarrollador
		if err := rows.Scan(&g.ID, &g.Nombre); err != nil {
			return nil, err
		}
		Desarrollador = append(Desarrollador, g)
	}
	return Desarrollador, nil
}

// Funcion plataforma
func ObtenerPlataforma() ([]models.Plataforma, error) {
	rows, err := GetDB().Query("SELECT id, nombre FROM plataformas")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var Plataforma []models.Plataforma
	for rows.Next() {
		var g models.Plataforma
		if err := rows.Scan(&g.ID, &g.Nombre); err != nil {
			return nil, err
		}
		Plataforma = append(Plataforma, g)
	}
	return Plataforma, nil
}

// Funciones videojuegos
func ObtenerVideojuegos() ([]models.Videojuego, error) {
	query := `
	SELECT 
		v.id, v.titulo, v.descripcion, v.genero_id, g.nombre AS genero_nombre,
		v.fecha_lanzamiento, v.desarrollador_id, d.nombre AS desarrollador_nombre,
		v.editor_id, e.nombre AS editor_nombre,
		v.precio, v.imagen_url, v.created_at, v.updated_at
	FROM 
		videojuegos v
	JOIN generos g ON v.genero_id = g.id
	JOIN desarrolladores d ON v.desarrollador_id = d.id
	JOIN editores e ON v.editor_id = e.id
	ORDER BY v.id;
	`

	rows, err := GetDB().Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videojuegos []models.Videojuego
	for rows.Next() {
		var v models.Videojuego
		err := rows.Scan(
			&v.ID, &v.Titulo, &v.Descripcion, &v.GeneroID, &v.GeneroNombre,
			&v.FechaLanzamiento, &v.DesarrolladorID, &v.DesarrolladorNombre,
			&v.EditorID, &v.EditorNombre,
			&v.Precio, &v.ImagenURL, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		videojuegos = append(videojuegos, v)
	}
	return videojuegos, nil
}

// Crear videojuego
func CrearVideojuego(v models.Videojuego) error {
	_, err := GetDB().Exec(`INSERT INTO videojuegos (titulo, descripcion, genero_id, fecha_lanzamiento, desarrollador_id, editor_id, precio, imagen_url) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		v.Titulo, v.Descripcion, v.GeneroID, v.FechaLanzamiento, v.DesarrolladorID, v.EditorID, v.Precio, v.ImagenURL)
	return err
}

// Eliminar videojuego
func EliminarVideojuego(id int) error {
	_, err := GetDB().Exec("DELETE FROM videojuegos WHERE id = $1", id)
	return err
}

// Existe titulo de videojuego?
func ExisteVideojuegoPorTitulo(titulo string) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM videojuegos WHERE titulo = $1", titulo).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
func ObtenerVideojuegoPorID(id int) (models.Videojuego, error) {
	query := `
	SELECT 
		v.id, v.titulo, v.descripcion, v.genero_id, g.nombre AS genero_nombre,
		v.fecha_lanzamiento, v.desarrollador_id, d.nombre AS desarrollador_nombre,
		v.editor_id, e.nombre AS editor_nombre,
		v.precio, v.imagen_url, v.created_at, v.updated_at
	FROM 
		videojuegos v
	JOIN generos g ON v.genero_id = g.id
	JOIN desarrolladores d ON v.desarrollador_id = d.id
	JOIN editores e ON v.editor_id = e.id
	WHERE v.id = $1
	`
	var v models.Videojuego
	err := GetDB().QueryRow(query, id).Scan(
		&v.ID, &v.Titulo, &v.Descripcion, &v.GeneroID, &v.GeneroNombre,
		&v.FechaLanzamiento, &v.DesarrolladorID, &v.DesarrolladorNombre,
		&v.EditorID, &v.EditorNombre,
		&v.Precio, &v.ImagenURL, &v.CreatedAt, &v.UpdatedAt,
	)
	return v, err
}

// Obtener videojuegos por género
func ObtenerVideojuegosPorGenero(generoID int) ([]models.Videojuego, error) {
	query := `
	SELECT 
		v.id, v.titulo, v.descripcion, v.genero_id, g.nombre AS genero_nombre,
		v.fecha_lanzamiento, v.desarrollador_id, d.nombre AS desarrollador_nombre,
		v.editor_id, e.nombre AS editor_nombre,
		v.precio, v.imagen_url, v.created_at, v.updated_at
	FROM 
		videojuegos v
	JOIN generos g ON v.genero_id = g.id
	JOIN desarrolladores d ON v.desarrollador_id = d.id
	JOIN editores e ON v.editor_id = e.id
	WHERE v.genero_id = $1
	ORDER BY v.id;
	`

	rows, err := GetDB().Query(query, generoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var videojuegos []models.Videojuego
	for rows.Next() {
		var v models.Videojuego
		err := rows.Scan(
			&v.ID, &v.Titulo, &v.Descripcion, &v.GeneroID, &v.GeneroNombre,
			&v.FechaLanzamiento, &v.DesarrolladorID, &v.DesarrolladorNombre,
			&v.EditorID, &v.EditorNombre,
			&v.Precio, &v.ImagenURL, &v.CreatedAt, &v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		videojuegos = append(videojuegos, v)
	}
	return videojuegos, nil
}
