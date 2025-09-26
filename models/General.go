package models

import "time"

type Genero struct {
	ID     int
	Nombre string
}

type Editor struct {
	ID     int
	Nombre string
}
type Desarrollador struct {
	ID     int
	Nombre string
}
type Plataforma struct {
	ID     int
	Nombre string
}
type Videojuego struct {
	ID                  int       `json:"id"`
	Titulo              string    `json:"titulo"`
	Descripcion         string    `json:"descripcion"`
	GeneroID            int       `json:"genero_id"`
	GeneroNombre        string    `json:"genero_nombre"`
	FechaLanzamiento    time.Time `json:"fecha_lanzamiento"`
	DesarrolladorID     int       `json:"desarrollador_id"`
	DesarrolladorNombre string    `json:"desarrollador_nombre"`
	EditorID            int       `json:"editor_id"`
	EditorNombre        string    `json:"editor_nombre"`
	Precio              float64   `json:"precio"`
	ImagenURL           string    `json:"imagen_url"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
