package models

import "time"

type Consulta struct {
	ID             int        `json:"id_consulta"`
	IDPaciente     int        `json:"id_paciente"`
	Tipo           string     `json:"tipo"`
	IDReceta       *int       `json:"id_receta"`       // puede ser null
	IDHorario      int        `json:"id_horario"`
	IDConsultorio  int        `json:"id_consultorio"`
	Diagnostico    string     `json:"diagnostico"`

	
	Costo          float64    `json:"costo"`
	FechaHora      time.Time  `json:"fecha_hora"`
}
