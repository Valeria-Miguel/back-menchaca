package models

type HistorialClinico struct {
	ID          int `json:"id_historial"`
	IDExpediente int `json:"id_expediente"`
	IDConsulta   int `json:"id_consulta"`
}
