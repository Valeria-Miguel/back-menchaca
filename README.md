# VErsionando -Codigo
sistema API REST, que gestiona  citas y 


# Fiber Task API

API RESTful desarrollada con Go y Fiber para la gestión de de citas y reportes utilizando autenticación JWT y almacenamiento en supabase.

---

## Tecnologías usadas

- **Go**
- **Fiber** (Framework web)
- **JWT** (Autenticación segura)
- **bcrypt** (Hash de contraseñas)

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