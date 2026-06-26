# Sistema de Gestión de Eventos y Entradas (Ticketera)

Práctico Integrador 2026 — Desarrollo de Software, Facultad de Ingeniería UCC.

Sistema tipo Ticketek para la gestión de eventos y venta de entradas. Permite a los usuarios explorar un catálogo de eventos, comprar entradas (con cupones de descuento opcionales), y gestionar su historial de compras (cancelación y traspaso a terceros). Los administradores cuentan con un panel completo para gestionar eventos, tipos de entrada, cupones promocionales y reportes de ocupación. El backend está desarrollado en Go (Gin + GORM) y el frontend en React (Vite + TypeScript).

## Tabla de Contenidos

- [Tecnologías Utilizadas](#tecnologías-utilizadas)
- [Requisitos Previos](#requisitos-previos)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Ejecución con Docker (recomendado)](#ejecución-con-docker-recomendado)
- [Ejecución Local](#ejecución-local)
- [Roles y Usuarios de Prueba](#roles-y-usuarios-de-prueba)
- [Funcionalidades](#funcionalidades)
- [Diagrama de Base de Datos](#diagrama-de-base-de-datos)
- [Decisiones de Diseño](#decisiones-de-diseño)
- [Testing](#testing)
- [Capturas de Pantalla](#capturas-de-pantalla)

## Tecnologías Utilizadas

- **Backend:** Go, Gin-Gonic, GORM, MySQL, JWT (golang-jwt), bcrypt
- **Frontend:** React, Vite, TypeScript, React Router DOM
- **Testing:** Go testing + httptest + SQLite en memoria
- **DevOps:** Docker, Docker Compose, Nginx
- **Control de versiones:** Git, GitHub, GitFlow

## Requisitos Previos

Para ejecución con Docker (recomendado):
- [Docker](https://www.docker.com/) y Docker Compose

Para ejecución local:
- [Go](https://go.dev/) 1.25 o superior
- [Node.js](https://nodejs.org/) 22 o superior
- [MySQL](https://www.mysql.com/) 8.0 corriendo localmente

## Estructura del Proyecto
.

├── backend/              # API REST en Go (Gin + GORM)

│   ├── controllers/      # Capa HTTP (manejo de requests/responses)

│   ├── services/         # Lógica de negocio

│   ├── dao/              # Acceso a datos (consultas GORM)

│   ├── domain/           # Modelos de dominio (entidades)

│   ├── dto/              # Objetos de transferencia de datos

│   ├── middlewares/      # Auth, autorización por roles y CORS

│   ├── utils/            # JWT y hasheo de contraseñas

│   ├── Dockerfile

│   └── main.go

├── frontend/             # Aplicación web en React + TypeScript

│   ├── src/

│   │   ├── pages/        # Vistas (incluye carpeta admin/)

│   │   ├── components/   # Navbar, ProtectedRoute

│   │   ├── context/      # AuthContext (manejo de sesión)

│   │   └── api.ts        # Cliente HTTP centralizado

│   ├── Dockerfile

│   └── nginx.conf

├── docs/                 # Diagrama de base de datos y capturas

├── docker-compose.yml

└── README.md

## Ejecución con Docker (recomendado)

Con Docker y Docker Compose instalados, desde la raíz del proyecto:

```bash
docker compose up --build
```

Esto levanta los tres servicios automáticamente:
- **Base de datos** (MySQL) en el puerto `3307` (externo)
- **Backend** (Go) en el puerto `8080`
- **Frontend** (Nginx) en el puerto `80`

Una vez levantado:
- Frontend: `http://localhost`
- Backend: `http://localhost:8080`

Para cargar datos de prueba, conectarse a la base con un cliente MySQL (host `localhost`, puerto `3307`, usuario `root`, contraseña `password`) y ejecutar el script de seed (ver sección de usuarios de prueba).

Para detener todo:

```bash
docker compose down
```

## Ejecución Local

### Base de datos

Crear la base de datos en MySQL:

```sql
CREATE DATABASE IF NOT EXISTS ticketera_db;
```

### Backend

1. Entrar a la carpeta del backend:
```bash
   cd backend
```

2. Crear un archivo `.env` en `backend/` (ajustar usuario y contraseña de MySQL):
DB_USER=root

DB_PASSWORD=password

DB_HOST=localhost

DB_PORT=3306

DB_NAME=ticketera_db

JWT_SECRET=desarrolo_de_software

3. Instalar dependencias y levantar:
```bash
   go mod tidy
   go run main.go
```

El backend queda en `http://localhost:8080`. Las tablas se crean automáticamente vía AutoMigrate.

### Frontend

1. En otra terminal:
```bash
   cd frontend
   npm install
   npm run dev
```

El frontend queda en `http://localhost:5173`.

## Roles y Usuarios de Prueba

El sistema maneja dos roles: **CLIENTE** (compra y gestiona sus entradas) y **ADMIN** (gestiona eventos, cupones y reportes).

Tras ejecutar el script de seed, los usuarios disponibles son (contraseña de todos: `123456`):

| Email | Rol |
|-------|-----|
| admin@test.com | ADMIN |
| gaston@test.com | CLIENTE |
| lucia@test.com | CLIENTE |
| martin@test.com | CLIENTE |
| sofia@test.com | CLIENTE |

## Funcionalidades

### Cliente
- Catálogo de eventos con búsqueda por nombre y filtro por estado.
- Detalle del evento con tipos de entrada, precios y stock.
- Compra de entradas con descuento por cupón opcional.
- Historial de entradas ("Mis Entradas").
- Cancelación de entradas (restaura el stock).
- Traspaso de entradas a otro usuario.

### Administrador
- Panel de gestión de eventos (tabla con acciones).
- Creación y edición de eventos con sus tipos de entrada.
- Cancelación de eventos.
- Reporte de ocupación/ventas por evento (métricas y compradores).
- Gestión de cupones promocionales (crear, editar, desactivar).

### Autenticación y Autorización
- Registro y login con hasheo bcrypt.
- JWT para autenticación.
- Autorización por roles: los endpoints de administrador validan que el token corresponda a un usuario ADMIN (responden 403 en caso contrario).

## Diagrama de Base de Datos

El modelo contempla cinco entidades: `usuarios`, `eventos`, `tipo_entradas`, `tickets` y `cupons`. Cada evento puede tener múltiples tipos de entrada con precio y stock propios. Cada ticket pertenece a un usuario y a un tipo de entrada, y opcionalmente a un cupón de descuento. Cada cupón está asociado a un evento. El archivo fuente del diagrama se encuentra en la carpeta `docs/`.

## Decisiones de Diseño

**1. Arquitectura en capas (MVC).** El backend separa responsabilidades en capas: controladores (HTTP), servicios (lógica de negocio), DAO (acceso a datos) y dominio (entidades). Esto facilita el testing aislado de cada capa y permite cambiar la implementación de una sin afectar las demás.

**2. Transacciones atómicas para operaciones críticas.** La compra y la cancelación de entradas se ejecutan dentro de transacciones de base de datos (`DB.Transaction`). Esto garantiza que la actualización del stock y la creación/modificación del ticket se realicen de forma atómica: si una operación falla, la otra se revierte automáticamente, evitando estados inconsistentes (por ejemplo, stock descontado sin ticket generado).

**3. Cancelación mediante cambio de estado (soft state).** Al cancelar una entrada o un evento, el registro no se elimina físicamente: se actualiza su campo `estado` a `CANCELADO`. Esto preserva el historial completo para reportes administrativos y evita romper relaciones de claves foráneas con registros referenciados.

**4. Autorización por roles con middleware encadenado.** La validación de permisos se implementa con dos middlewares: `AuthMiddleware` valida el token JWT y extrae el rol, y `AdminMiddleware` verifica que ese rol sea ADMIN. Los endpoints de administrador encadenan ambos, devolviendo 401 si falta el token y 403 si el usuario no tiene permisos suficientes.

**5. Precio congelado en el ticket.** El ticket guarda el `precio_pagado` al momento de la compra (con descuento de cupón ya aplicado). Si luego cambia el precio del tipo de entrada, el ticket conserva lo que el usuario efectivamente pagó.

## Testing

El backend cuenta con pruebas unitarias (servicios) y de integración HTTP (controladores, usando `httptest`), cubriendo casos de éxito y de error. Se utiliza SQLite en memoria para aislar los tests de la base de datos real.

Para correr los tests con reporte de cobertura:

```bash
cd backend
go test ./services/... ./controllers/... -cover
```

Cobertura alcanzada: superior al 80% tanto en la capa de servicios como en la de controladores.