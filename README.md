# Backend + Base de datos
sistema API REST, que gestiona  un Sistema de Citas y reportes de un Hospital 

---
# Fiber Task API

API RESTful desarrollada con Go y Fiber para la gestión de de citas y reportes utilizando autenticación JWT y almacenamiento en supabase.

---

## Tecnologías usadas

- **Go**
- **Fiber**  - Framework web
- **JWT** - Autenticación segura
- **bcrypt** - Hash de contraseñas
-**Supabase** – Base de datos PostgreSQL en la nube
- **database/sql**  – Conexión directa
- **Validaciones manuales** - sin dependencias externas


---
## Requisitos

- Go 1.20 o superior
- Cuenta en Supabase 
- Editor como Visual Studio Code 

---
### Clonar el repositorio

```bash
git clone https://github.com/Valeria-Miguel/back-menchaca.git
```
# ir al proyecto
cd back-menchaca

# instalar dependencias
go mod tidy

# ejecutar el servidor
go run main.go

##  Configuración `.env`

```bash
DATABASE_URL= tu conexion a sudabase

JWT_SECRET=tu clave de JWT
```
---

## Estructura del Proyecto

├── main.go → Punto de entrada
├── go.mod → Módulo de Go
├── config/ → Configuración de MongoDB
├── models/ → Modelos (User, Task)
├── handlers/ → Lógica de endpoints
├── routes/ → Definición de rutas
├── middleware/ → Middleware JWT
├── utils/ → Funciones auxiliares (JWT)
├── test/ → Pruebas automatizadas (Aun no implementada)
└── .env → Variables de entorno

--- 

## Características principales

- Gestión completa de pacientes y empleados (registro, consulta, actualización, eliminación)
- CRUD de expedientes clínicos y antecedentes médicos
- Módulo de recetas y consultas médicas
- Historial clínico vinculado a consultas previas
- Registro de consentimiento con aviso de privacidad
- Autenticación con JWT y protección por roles (paciente / empleado)
- Validación de datos con funciones personalizadas
- Conexión a Supabase/PostgreSQL

---
## Seguridad y privacidad

Este proyecto implementa medidas  de seguridad:

- Encriptación de contraseñas con `bcrypt`
- Generación y verificación de tokens JWT con expiración
- Protección de rutas según rol del usuario (empleado o paciente)
- Aviso de privacidad accesible por ruta `/api/consentimiento/aviso-privacidad`
- Registro del consentimiento informado del paciente
