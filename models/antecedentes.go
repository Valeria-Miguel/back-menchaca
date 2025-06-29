package models

import "time"

type Antecedente struct {
    ID            int       `json:"id_antecedente"`
    IDExpediente  int       `json:"id_expediente"`
    Diagnostico   string    `json:"diagnostico"`
    Descripcion   string    `json:"descripcion"`
    Fecha         time.Time `json:"fecha"`
}
