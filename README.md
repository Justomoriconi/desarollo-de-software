# desarollo-de-software

# Sistema de Gestión de Eventos y Entradas (Ticketera)

Práctico Integrador 2026 — Desarrollo de Software, Facultad de Ingeniería UCC.

Sistema tipo Ticketek para la gestión de eventos y venta de entradas, con backend en Go (Gin + GORM) y frontend en React (Vite + TypeScript).

## Requisitos previos

- [Go](https://go.dev/) 1.25 o superior
- [Node.js](https://nodejs.org/) 22.14 o superior
- [MySQL](https://www.mysql.com/) corriendo localmente

## Estructura del proyecto
.

├── backend/    # API REST en Go (Gin + GORM)

└── frontend/   # Aplicación web en React + TypeScript

## Base de datos

Crear la base de datos en MySQL (las tablas se crean solas vía AutoMigrate):

```sql
CREATE DATABASE IF NOT EXISTS ticketera_db;
```

## Backend

1. Entrar a la carpeta del backend:
```bash
   cd backend
```

2. Crear un archivo `.env` en `backend/` con el siguiente contenido (ajustar usuario y contraseña de MySQL):
DB_USER=root

DB_PASSWORD=password

DB_HOST=localhost

DB_PORT=3306

DB_NAME=ticketera_db

JWT_SECRET=desarrolo_de_software

3. Instalar dependencias:
```bash
   go mod tidy
```

4. Levantar el servidor:
```bash
   go run main.go
```

El backend queda corriendo en `http://localhost:8080`. Al iniciar, se conecta a MySQL y crea automáticamente las tablas necesarias.

### Cargar datos de prueba (opcional)

```sql
USE ticketera_db;

INSERT INTO eventos (nombre, descripcion, fecha, lugar, estado)
VALUES ('Recital de prueba', 'Evento de test', '2026-08-15 21:00:00', 'Estadio Córdoba', 'ACTIVO');

INSERT INTO tipo_entradas (nombre, precio, stock_disponible, evento_id)
VALUES ('General', 5000.00, 100, 1);
```

## Frontend

1. En otra terminal, entrar a la carpeta del frontend:
```bash
   cd frontend
```

2. Instalar dependencias:
```bash
   npm install
```

3. Levantar el servidor de desarrollo:
```bash
   npm run dev
```

El frontend queda disponible en `http://localhost:5173` y se conecta automáticamente al backend en `http://localhost:8080`.

## Flujo de uso

1. Registrarse o iniciar sesión.
2. Explorar el catálogo de eventos (con búsqueda por nombre y filtro por estado).
3. Ver el detalle de un evento y comprar una entrada.
4. Ir a "Mis Entradas" para ver el historial, cancelar o transferir entradas a otro usuario.

## Testing

Para correr los tests del backend (servicios, con SQLite en memoria):

```bash
cd backend
go test ./services/... -cover
```

## Tecnologías

- **Backend:** Go, Gin-Gonic, GORM, MySQL, JWT, bcrypt
- **Frontend:** React, Vite, TypeScript, React Router DOM