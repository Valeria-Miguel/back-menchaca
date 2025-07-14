# CHANGELOG

Todas las modificaciones importantes de este proyecto se documentan aqui

---
## [1.1] - 01-07-28



## [1.0] - 2025-06-28

###  Versión inicial

- Se publica la primera versión completa del sistema hospitalario.
- Se integran los módulos:
  - Pacientes
  - Empleados
  - Expedientes clínicos
  - Antecedentes médicos
  - Consultas
  - Consultorios
  - Horario
  - Recetas
  - Historial clínico
  - Aviso de privacidad y consentimiento
- Autenticación JWT para rutas protegidas
- Hash de contraseñas con `bcrypt`
- Conexión segura con Supabase (PostgreSQL)
- Rutas protegidas según roles (`paciente`, `empleado`)

### Nuevas funcionalidades
- CRUD de consultas
- CRUD de expediente, antecedentes y historial clinico
- Validaciones 

### Cambios
- Se integraron mas validaciones al utils/validate.go
---

## [0.2] - 2025-06-27

### Nuevas funcionalidades

- CRUD de empleado y pasiente  y login con validación de roles 
- Implementación de middleware JWT
- Se integran validaciones de campos obligatorios
- CRUD de consultorio y horario

### Cambios

- Se configuró conexión a Supabase usando `.env`

---

## [0.1] - 2025-06-26

### Estructura base

- Se crea la estructura base del proyecto con:
  - `config/`, `models/`, `routes/`, `handlers/`, `utils/`
- Se realiza la conexión inicial a Supabase
- Se crea repositorio en GitHub: `back-menchaca`
- Se añade archivo `.env` y se configura en `.gitignore`

---

## Branches principales

- `main` → Rama principal de despliegue 

---

