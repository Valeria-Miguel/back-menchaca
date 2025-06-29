package models

type Empleado struct {
	ID         int    `json:"id_empleado"`
	Nombre     string `json:"nombre"`
	Appaterno  string `json:"appaterno"`
	Apmaterno  string `json:"apmaterno"`
	Tipo       string `json:"tipo_empleado"`
	Area       string `json:"area"`
	Correo     string `json:"correo"`
	Contrasena string `json:"contrasena,omitempty"`
}
