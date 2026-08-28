# RAG Chat

[English](#english) | [Español](#español)

## English

### Overview

RAG Chat is a small Retrieval-Augmented Generation application written in Go. It combines a conversational LLM with documents stored in PostgreSQL through pgvector, allowing answers to use relevant excerpts from a local document collection.

The project provides both a terminal REPL and a browser chat. It supports streamed responses, automatic document ingestion, text and image uploads, conversation-aware query rewriting, and optional image captioning through a vision model.

This is intentionally a compact, local project. It focuses on showing the complete RAG flow without adding production infrastructure.

### Features

- Terminal and web chat interfaces.
- Streamed LLM responses.
- Retrieval from PostgreSQL with pgvector cosine similarity.
- Conversation-aware rewriting of follow-up questions.
- Automatic ingestion of `.txt`, `.md`, and `.markdown` files.
- Browser uploads for text documents and images.
- Manual image descriptions and optional automatic captions.
- Markdown rendering with HTML sanitization in the browser.
- Graceful fallback to regular chat when the vector store is unavailable.

### Technologies

| Area | Technology |
| --- | --- |
| Language | Go 1.26.3 |
| LLM and embeddings | OpenAI Go SDK v3 with OpenAI-compatible endpoints |
| HTTP routing | Chi v5 |
| Database | PostgreSQL 18 and pgvector |
| PostgreSQL client | pgx v5 |
| File watching | fsnotify |
| Configuration | godotenv and environment variables |
| Server-side UI | Go `html/template` and embedded templates |
| Browser UI | Tailwind CSS, Marked, DOMPurify, and vanilla JavaScript |
| Local infrastructure | Docker Compose |

### How It Works

The ingestion flow is:

```text
Text or image description
        -> chunking
        -> embedding generation
        -> PostgreSQL / pgvector
```

The question-answering flow is:

```text
Conversation
        -> standalone query rewrite
        -> query embedding
        -> top matching document chunks
        -> context added to the latest question
        -> streamed LLM response
```

The retriever currently requests the five closest chunks. Retrieved context includes source metadata so the model can cite filenames in its answer.

### Project Structure

```text
.
├── app/                    Application wiring, lifecycle, and shutdown
├── chat/                   Terminal REPL
├── cmd/rag/                Executable entry point
├── config/                 Environment configuration and defaults
├── documents/              Source, processed, and uploaded documents
├── ingest/                 Chunking, file watching, and document ingestion
├── llm/                    Chat, embedding, and vision model clients
├── prompts/                Optional custom system prompts
├── rag/                    Query rewriting, retrieval, and context formatting
├── vector/                 Vector store interface
│   └── pgvector/           PostgreSQL/pgvector implementation and migration
├── web/                    HTTP server and embedded browser UI
│   └── templates/          Go HTML templates
├── create-table.sql        Optional manual schema reference
├── docker-compose.yml      Local PostgreSQL/pgvector service
├── go.mod                  Go module and dependencies
└── README.md               Project documentation
```

### Requirements

- Go 1.26.3 or a compatible version.
- A chat completion endpoint compatible with the OpenAI API.
- An embedding endpoint compatible with the OpenAI API.
- Docker Compose, when using the included PostgreSQL/pgvector service.

The embedding model output size must match `EMBEDDING_DIM`. The included defaults use 768 dimensions.

### Configuration

Configuration is loaded from environment variables and, when present, a local `.env` file. The `.env` file is ignored by Git.

A typical local configuration looks like this:

```dotenv
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_API_KEY=replace-me
OPENAI_MODEL=gpt-4o-mini

# Example of a separate OpenAI-compatible embedding provider.
EMBEDDING_BASE_URL=http://localhost:11434/v1
EMBEDDING_API_KEY=
EMBEDDING_MODEL=nomic-embed-text
EMBEDDING_DIM=768

DATABASE_URL=postgres://rag:rag@localhost:5432/rag?sslmode=disable
HTTP_ADDR=:8080

INGEST_DIR=./documents
PROCESS_DIR=./documents/processed
IMAGE_DIR=./documents/images
SYSTEM_PROMPT_FILE=./prompts/system-custom.md

# Optional. Enables automatic image captions when supported by the chat provider.
VISION_MODEL=
```

Important behavior:

- `HTTP_ADDR` enables the web interface. If it is empty, only the terminal REPL runs.
- `DATABASE_URL` enables ingestion and retrieval. Without it, the application still works as a regular chat.
- If `EMBEDDING_BASE_URL` is empty, the application reuses `OPENAI_BASE_URL`.
- `VISION_MODEL` is optional. Images can still be indexed with a manual description when it is empty.
- The application creates the pgvector extension, table, and index automatically when it connects.

### Running the Project

Start PostgreSQL and pgvector:

```bash
docker compose up -d
```

Run the application:

```bash
go run ./cmd/rag
```

The terminal starts an interactive chat session. When `HTTP_ADDR=:8080`, the browser chat is available at:

```text
http://localhost:8080/
```

The root route redirects to `/chat`. The web server runs alongside the terminal REPL and stops when the main process exits.

### Adding Documents

There are two ways to add text documents:

1. Place a `.txt`, `.md`, or `.markdown` file in `INGEST_DIR`. After successful ingestion, the watcher moves it to `PROCESS_DIR`.
2. Use the attachment button in the browser chat.

Images are added through the browser. Each image needs a searchable description. When `VISION_MODEL` is configured, the interface can generate a draft caption for review before indexing it.

### Verification

```bash
go test ./...
go vet ./...
```

> This project is designed for local use and does not include authentication. Do not expose it publicly without adding access control and appropriate operational safeguards.

---

## Español

### Descripción

RAG Chat es una aplicación pequeña de Generación Aumentada por Recuperación escrita en Go. Combina un LLM conversacional con documentos almacenados en PostgreSQL mediante pgvector, permitiendo que las respuestas utilicen fragmentos relevantes de una colección local.

El proyecto ofrece un REPL por terminal y un chat en el navegador. Incluye respuestas por streaming, ingestión automática de documentos, carga de textos e imágenes, reescritura de consultas según la conversación y descripción opcional de imágenes mediante un modelo de visión.

El proyecto es intencionalmente compacto y está pensado para uso local. Su objetivo es mostrar el flujo RAG completo sin agregar infraestructura de producción.

### Funcionalidades

- Chat mediante terminal y navegador.
- Respuestas del LLM por streaming.
- Recuperación desde PostgreSQL utilizando similitud coseno con pgvector.
- Reescritura de preguntas de seguimiento según el contexto de la conversación.
- Ingestión automática de archivos `.txt`, `.md` y `.markdown`.
- Carga de documentos de texto e imágenes desde el navegador.
- Descripción manual de imágenes y captions automáticos opcionales.
- Renderizado de Markdown con sanitización de HTML en el navegador.
- Degradación controlada a un chat normal cuando el vector store no está disponible.

### Tecnologías

| Área | Tecnología |
| --- | --- |
| Lenguaje | Go 1.26.3 |
| LLM y embeddings | OpenAI Go SDK v3 con endpoints compatibles con OpenAI |
| Enrutamiento HTTP | Chi v5 |
| Base de datos | PostgreSQL 18 y pgvector |
| Cliente PostgreSQL | pgx v5 |
| Monitoreo de archivos | fsnotify |
| Configuración | godotenv y variables de entorno |
| Interfaz del servidor | `html/template` de Go y templates embebidos |
| Interfaz del navegador | Tailwind CSS, Marked, DOMPurify y JavaScript puro |
| Infraestructura local | Docker Compose |

### Cómo Funciona

El flujo de ingestión es:

```text
Texto o descripción de imagen
        -> división en fragmentos
        -> generación de embeddings
        -> PostgreSQL / pgvector
```

El flujo de preguntas y respuestas es:

```text
Conversación
        -> reescritura como consulta independiente
        -> embedding de la consulta
        -> fragmentos de documentos más cercanos
        -> contexto agregado a la última pregunta
        -> respuesta del LLM por streaming
```

Actualmente, el retriever solicita los cinco fragmentos más cercanos. El contexto recuperado incluye metadatos de origen para que el modelo pueda citar los nombres de archivo en su respuesta.

### Estructura del Proyecto

```text
.
├── app/                    Ensamblado, ciclo de vida y cierre de la aplicación
├── chat/                   REPL de terminal
├── cmd/rag/                Punto de entrada del ejecutable
├── config/                 Configuración de entorno y valores predeterminados
├── documents/              Documentos fuente, procesados y cargados
├── ingest/                 Fragmentación, monitoreo e ingestión de documentos
├── llm/                    Clientes de chat, embeddings y modelo de visión
├── prompts/                Prompts de sistema personalizados opcionales
├── rag/                    Reescritura, recuperación y formato del contexto
├── vector/                 Interfaz del vector store
│   └── pgvector/           Implementación y migración de PostgreSQL/pgvector
├── web/                    Servidor HTTP e interfaz web embebida
│   └── templates/          Templates HTML de Go
├── create-table.sql        Referencia opcional del esquema manual
├── docker-compose.yml      Servicio local de PostgreSQL/pgvector
├── go.mod                  Módulo y dependencias de Go
└── README.md               Documentación del proyecto
```

### Requisitos

- Go 1.26.3 o una versión compatible.
- Un endpoint de chat compatible con la API de OpenAI.
- Un endpoint de embeddings compatible con la API de OpenAI.
- Docker Compose para utilizar el servicio PostgreSQL/pgvector incluido.

El tamaño de salida del modelo de embeddings debe coincidir con `EMBEDDING_DIM`. Los valores predeterminados incluidos utilizan 768 dimensiones.

### Configuración

La configuración se carga desde variables de entorno y, cuando existe, desde un archivo local `.env`. El archivo `.env` está ignorado por Git.

Una configuración local típica es:

```dotenv
OPENAI_BASE_URL=https://api.openai.com/v1
OPENAI_API_KEY=reemplazar
OPENAI_MODEL=gpt-4o-mini

# Ejemplo de un proveedor de embeddings separado y compatible con OpenAI.
EMBEDDING_BASE_URL=http://localhost:11434/v1
EMBEDDING_API_KEY=
EMBEDDING_MODEL=nomic-embed-text
EMBEDDING_DIM=768

DATABASE_URL=postgres://rag:rag@localhost:5432/rag?sslmode=disable
HTTP_ADDR=:8080

INGEST_DIR=./documents
PROCESS_DIR=./documents/processed
IMAGE_DIR=./documents/images
SYSTEM_PROMPT_FILE=./prompts/system-custom.md

# Opcional. Habilita captions automáticos cuando el proveedor de chat lo soporta.
VISION_MODEL=
```

Comportamientos importantes:

- `HTTP_ADDR` habilita la interfaz web. Si está vacío, solamente se ejecuta el REPL de terminal.
- `DATABASE_URL` habilita la ingestión y recuperación. Sin esta variable, la aplicación continúa funcionando como un chat normal.
- Si `EMBEDDING_BASE_URL` está vacío, la aplicación reutiliza `OPENAI_BASE_URL`.
- `VISION_MODEL` es opcional. Sin esta variable, las imágenes todavía pueden indexarse con una descripción manual.
- La aplicación crea automáticamente la extensión pgvector, la tabla y el índice al conectarse.

### Ejecución

Iniciar PostgreSQL y pgvector:

```bash
docker compose up -d
```

Ejecutar la aplicación:

```bash
go run ./cmd/rag
```

La terminal inicia una sesión de chat interactiva. Con `HTTP_ADDR=:8080`, el chat web queda disponible en:

```text
http://localhost:8080/
```

La ruta raíz redirige a `/chat`. El servidor web se ejecuta junto al REPL de terminal y se detiene cuando finaliza el proceso principal.

### Agregar Documentos

Existen dos maneras de agregar documentos de texto:

1. Colocar un archivo `.txt`, `.md` o `.markdown` en `INGEST_DIR`. Después de una ingestión correcta, el watcher lo mueve a `PROCESS_DIR`.
2. Utilizar el botón de adjuntar en el chat del navegador.

Las imágenes se agregan desde el navegador. Cada imagen necesita una descripción que pueda utilizarse en las búsquedas. Cuando `VISION_MODEL` está configurado, la interfaz puede generar un caption preliminar para revisarlo antes de indexarlo.

### Verificación

```bash
go test ./...
go vet ./...
```

> Este proyecto está diseñado para uso local y no incluye autenticación. No debe exponerse públicamente sin agregar control de acceso y las protecciones operativas correspondientes.
