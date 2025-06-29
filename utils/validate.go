package utils

import (
	"html"
	"errors"
	"regexp"
	"strings"
	"unicode"
	"time"
	"back-menchaca/models"
	"back-menchaca/config"
)
func SanitizarInput(s string) string {
	// Quita espacios y caracteres HTML peligrosos
	s = strings.TrimSpace(s)
	s = html.EscapeString(s) // <- evita que el string tenga etiquetas <script> etc.
	return s
} 
func ValidarCorreo(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func ValidarTextoLetras(input string) bool {
	re := regexp.MustCompile(`^[A-Za-zÁÉÍÓÚáéíóúÑñ\s]+$`)
	return re.MatchString(input)
}

// ✅ Validar contraseña manualmente
func ValidarContrasena(pw string) error {
	if len(pw) < 8 {
		return errors.New("La contraseña debe tener al menos 8 caracteres")
	}
	var hasUpper, hasLower, hasNumber, hasSpecial bool
	for _, c := range pw {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsDigit(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errors.New("La contraseña debe incluir mayúsculas, minúsculas, número y símbolo")
	}
	return nil
}

func ValidarPaciente(nombre, appaterno, correo, contrasena string) error {
	if strings.TrimSpace(nombre) == "" || !ValidarTextoLetras(nombre) {
		return errors.New("Nombre inválido")
	}
	if strings.TrimSpace(appaterno) == "" || !ValidarTextoLetras(appaterno) {
		return errors.New("Apellido paterno inválido")
	}
	if correo == "" || !ValidarCorreo(correo) {
		return errors.New("Correo inválido")
	}
	if err := ValidarContrasena(contrasena); err != nil {
		return err
	}
	return nil
}

func ValidarEmpleado(nombre, appaterno, tipo, area, correo, contrasena string) error {
	if strings.TrimSpace(nombre) == "" || !ValidarTextoLetras(nombre) {
		return errors.New("nombre inválido")
	}
	if strings.TrimSpace(appaterno) == "" || !ValidarTextoLetras(appaterno) {
		return errors.New("apellido paterno inválido")
	}
	if strings.TrimSpace(tipo) == "" || !ValidarTextoLetras(tipo) {
		return errors.New("tipo de empleado inválido")
	}
	if tipo != "doctor" && tipo != "enfermera" {
		return errors.New("tipo de empleado inválido. Solo se permite 'doctor' o 'enfermera'")
	}
	if strings.TrimSpace(area) == "" || !ValidarTextoLetras(area) {
		return errors.New("area inválida")
	}
	if correo == "" || !ValidarCorreo(correo) {
		return errors.New("correo inválido")
	}
	if err := ValidarContrasena(contrasena); err != nil {
		return err
	}
	return nil
}

func ValidarConsultorio(nombre, tipo string) error {
	if strings.TrimSpace(nombre) == "" {
		return errors.New("Nombre del consultorio requerido")
	}
	if strings.TrimSpace(tipo) == "" {
		return errors.New("Tipo de consultorio requerido")
	}
	
	return nil
}


func ValidarHorario(turno string, idEmpleado, idConsultorio int) error {
	turno = strings.ToLower(turno)
	if turno != "matutino" && turno != "vespertino" {
		return errors.New("Turno inválido. Debe ser 'matutino' o 'vespertino'")
	}
	if idEmpleado <= 0 {
		return errors.New("ID de empleado inválido")
	}
	if idConsultorio <= 0 {
		return errors.New("ID de consultorio inválido")
	}
	return nil
}


func ValidarConsulta(c models.Consulta) error {
	if c.IDPaciente == 0 || c.IDHorario == 0 || c.IDConsultorio == 0 {
		return errors.New("Faltan campos obligatorios")
	}
	if !ValidarTextoLetras(c.Tipo) {
		return errors.New("Tipo inválido")
	}
	if c.Costo < 0 {
		return errors.New("El costo no puede ser negativo")
	}
	return nil
}


func ExisteID(tabla string, columna string, id int) bool {
	var existe bool
	query := "SELECT EXISTS(SELECT 1 FROM " + tabla + " WHERE " + columna + " = $1)"
	err := config.DB.QueryRow(query, id).Scan(&existe)
	return err == nil && existe
}

func ExisteIDHisotial(tabla string, campo string, id int) bool {
	var count int
	err := config.DB.QueryRow("SELECT COUNT(*) FROM "+tabla+" WHERE "+campo+" = $1", id).Scan(&count)
	return err == nil && count > 0
}


func ValidarReceta(medicamento, dosis string, idConsultorio int) error {
	if strings.TrimSpace(medicamento) == "" {
		return errors.New("El medicamento es obligatorio")
	}
	if strings.TrimSpace(dosis) == "" {
		return errors.New("La dosis es obligatoria")
	}
	if idConsultorio <= 0 || !ExisteID("Consultorios", "id_consultorio", idConsultorio) {
		return errors.New("El ID de consultorio no es válido")
	}
	return nil
}


// ValidarSeguro verifica que el campo seguro no esté vacío y no contenga caracteres raros
func ValidarSeguro(seguro string) error {
    seguro = strings.TrimSpace(seguro)
    if seguro == "" {
        return errors.New("El campo 'seguro' no puede estar vacío")
    }
    // Puedes agregar más validaciones específicas si quieres (longitud, formato, etc.)
    return nil
}


func ExisteIDExped(id int) bool {
	var existe bool
	query := `SELECT EXISTS(SELECT 1 FROM Expediente WHERE id_expediente = $1)`
	err := config.DB.QueryRow(query, id).Scan(&existe)
	return err == nil && existe
}

// ValidarAntecedente valida los campos de Antecedentes
func ValidarAntecedente(diagnostico, descripcion string, fecha time.Time, idExpediente int) error {
    if idExpediente <= 0 {
        return errors.New("ID de expediente inválido")
    }
    diagnostico = strings.TrimSpace(diagnostico)
    if diagnostico == "" {
        return errors.New("Diagnóstico no puede estar vacío")
    }
    descripcion = strings.TrimSpace(descripcion)
    if descripcion == "" {
        return errors.New("Descripción no puede estar vacía")
    }
    if fecha.IsZero() {
        return errors.New("Fecha inválida o no proporcionada")
    }
    if fecha.After(time.Now()) {
        return errors.New("La fecha no puede ser futura")
    }
    return nil
}

