package models

import "time"

type Consentimiento struct {
	ID         int       `json:"id" db:"id"`
	IDPaciente int       `json:"id_paciente" db:"id_paciente"`
	FechaHora  time.Time `json:"fecha_hora" db:"fecha_hora"`
}
