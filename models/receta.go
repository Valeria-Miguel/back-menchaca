package models

import "time"

type Receta struct {
	ID           int       `json:"id_receta"`
	Fecha        time.Time `json:"fecha"`
	Medicamento  string    `json:"medicamento"`
	Dosis        string    `json:"dosis"`
	IDConsultorio int      `json:"id_consultorio"`
}
