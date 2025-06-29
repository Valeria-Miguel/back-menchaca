package models

type Paciente struct {
	ID         int    `json:"id_paciente"`
	Nombre     string `json:"nombre"`
	Appaterno  string `json:"appaterno"`
	Apmaterno  string `json:"apmaterno"`
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena,omitempty"` // omitida en respuestas
}
