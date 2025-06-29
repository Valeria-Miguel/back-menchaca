package models

type Horario struct {
	ID             int    `json:"id_horario"`
	IDConsultorio  int    `json:"id_consultorio"`
	Turno          string `json:"turno"`
	IDEmpleado     int    `json:"id_empleado"`
}
