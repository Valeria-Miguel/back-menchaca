package models

import "time"

type Expediente struct {
    ID           int       `json:"id_expediente"`
    IDPaciente   int       `json:"id_paciente"`
    Seguro       string    `json:"seguro"`
    FechaCreacion time.Time `json:"fecha_creacion"`
}
