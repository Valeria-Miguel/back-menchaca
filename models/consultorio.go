package models

type Consultorio struct {
	ID         int    `json:"id_consultorio"`
	Tipo       string `json:"tipo"`
	Nombre     string `json:"nombre"`
}
