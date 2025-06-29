package handlers

import (
	"back-menchaca/config"
	"github.com/gofiber/fiber/v2"
)

func ReporteDetalleConsultasPorPaciente(c *fiber.Ctx) error {
	var body struct {
		IDPaciente int `json:"id_paciente"`
	}
	if err := c.BodyParser(&body); err != nil || body.IDPaciente == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID de paciente inválido"})
	}

	query := `
		SELECT 
			c.id_consulta,
			p.nombre AS paciente,
			em.nombre AS empleado,
			h.turno,
			co.nombre AS consultorio,
			c.tipo,
			c.diagnostico,
			c.costo,
			c.fecha_hora
		FROM Consultas c
		JOIN Paciente p ON c.id_paciente = p.id_paciente
		JOIN Horarios h ON c.id_horario = h.id_horario
		JOIN Empleado em ON h.id_empleado = em.id_empleado
		JOIN Consultorios co ON c.id_consultorio = co.id_consultorio
		WHERE p.id_paciente = $1
		ORDER BY c.fecha_hora DESC
	`

	rows, err := config.DB.Query(query, body.IDPaciente)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener consultas del paciente"})
	}
	defer rows.Close()

	type DetallePacienteConsulta struct {
		ID          int     `json:"id_consulta"`
		Paciente    string  `json:"paciente"`
		Empleado    string  `json:"empleado"`
		Turno       string  `json:"turno"`
		Consultorio string  `json:"consultorio"`
		Tipo        string  `json:"tipo"`
		Diagnostico string  `json:"diagnostico"`
		Costo       float64 `json:"costo"`
		FechaHora   string  `json:"fecha_hora"`
	}

	var resultados []DetallePacienteConsulta
	for rows.Next() {
		var d DetallePacienteConsulta
		if err := rows.Scan(&d.ID, &d.Paciente, &d.Empleado, &d.Turno, &d.Consultorio, &d.Tipo, &d.Diagnostico, &d.Costo, &d.FechaHora); err == nil {
			resultados = append(resultados, d)
		}
	}
	return c.JSON(resultados)
}



func ReporteConsultasPorArea(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT e.area, COUNT(*) FROM Consultas c JOIN Horarios h ON c.id_horario = h.id_horario JOIN Empleado e ON h.id_empleado = e.id_empleado GROUP BY e.area`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener reporte"})
	}
	defer rows.Close()

	var data []fiber.Map
	for rows.Next() {
		var area string
		var total int
		if err := rows.Scan(&area, &total); err == nil {
			data = append(data, fiber.Map{"area": area, "total": total})
		}
	}
	return c.JSON(data)
}

func ReporteConsultasPorTurno(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT h.turno, COUNT(*) FROM Consultas c JOIN Horarios h ON c.id_horario = h.id_horario GROUP BY h.turno`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener reporte"})
	}
	defer rows.Close()

	var data []fiber.Map
	for rows.Next() {
		var turno string
		var total int
		if err := rows.Scan(&turno, &total); err == nil {
			data = append(data, fiber.Map{"turno": turno, "total": total})
		}
	}
	return c.JSON(data)
}

func ReporteIngresosPorConsultorio(c *fiber.Ctx) error {
	rows, err := config.DB.Query(`SELECT cons.nombre, SUM(c.costo) FROM Consultas c JOIN Consultorios cons ON c.id_consultorio = cons.id_consultorio GROUP BY cons.nombre`)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener reporte"})
	}
	defer rows.Close()

	var data []fiber.Map
	for rows.Next() {
		var consultorio string
		var total float64
		if err := rows.Scan(&consultorio, &total); err == nil {
			data = append(data, fiber.Map{"consultorio": consultorio, "total": total})
		}
	}
	return c.JSON(data)
}

func ReporteDetallesConsultaExpediente(c *fiber.Ctx) error {
	var body struct {
		IDExpediente int `json:"id_expediente"`
	}
	if err := c.BodyParser(&body); err != nil || body.IDExpediente == 0 {
		return c.Status(400).JSON(fiber.Map{"error": "ID de expediente inválido"})
	}

	rows, err := config.DB.Query(`
		SELECT a.diagnostico, a.descripcion, a.fecha, c.tipo, c.fecha_hora, c.diagnostico
		FROM Historial_Clinico h
		JOIN Consultas c ON h.id_consultas = c.id_consulta
		JOIN Antecedentes a ON a.id_expediente = h.id_expediente
		WHERE h.id_expediente = $1
	`, body.IDExpediente)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener detalles del historial clínico"})
	}
	defer rows.Close()

	var resultados []fiber.Map
	for rows.Next() {
		var diagAnt, descAnt, tipo, diagCons string
		var fechaAnt, fechaCons string
		err := rows.Scan(&diagAnt, &descAnt, &fechaAnt, &tipo, &fechaCons, &diagCons)
		if err == nil {
			resultados = append(resultados, fiber.Map{
				"diagnostico_antecedente": diagAnt,
				"descripcion": descAnt,
				"fecha_antecedente": fechaAnt,
				"tipo_consulta": tipo,
				"fecha_consulta": fechaCons,
				"diagnostico_consulta": diagCons,
			})
		}
	}
	return c.JSON(resultados)
}

func ObtenerDetalleSimpleConsultas(c *fiber.Ctx) error {
	query := `
		SELECT
			c.id_consulta,
			p.nombre AS paciente,
			em.nombre AS empleado,
			c.tipo,
			c.diagnostico,
			c.costo,
			c.fecha_hora
		FROM Consultas c
		LEFT JOIN Paciente p ON c.id_paciente = p.id_paciente
		LEFT JOIN Horarios h ON c.id_horario = h.id_horario
		LEFT JOIN Empleado em ON h.id_empleado = em.id_empleado
	`

	rows, err := config.DB.Query(query)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Error al obtener detalle de consultas"})
	}
	defer rows.Close()

	type DetalleSimpleConsulta struct {
		ID          int     `json:"id_consulta"`
		Paciente    string  `json:"paciente"`
		Empleado    string  `json:"empleado"`
		Tipo        string  `json:"tipo"`
		Diagnostico string  `json:"diagnostico"`
		Costo       float64 `json:"costo"`
		FechaHora   string  `json:"fecha_hora"`
	}

	var detalles []DetalleSimpleConsulta
	for rows.Next() {
		var d DetalleSimpleConsulta
		if err := rows.Scan(&d.ID, &d.Paciente, &d.Empleado, &d.Tipo, &d.Diagnostico, &d.Costo, &d.FechaHora); err == nil {
			detalles = append(detalles, d)
		}
	}
	return c.JSON(detalles)
}
