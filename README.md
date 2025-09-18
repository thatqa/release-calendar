# 📅 Release Calendar

![Release Calendar UI](https://thatqa.com/release-calendar-v2.png)

Release Calendar is a simple yet powerful web application for managing software releases.  
It provides a **calendar view**, detailed **release pages**, **comments**, **links**, and an **AI-powered summary** of discussions and notes.

---

## ✨ Features

- 📌 Calendar view with daily status markers:
    - 🟥 **Failed**
    - 🟩 **Success**
    - 🟦 **Planned**
- 🔍 Filter releases by status or duty user
- 📝 Notes, duty users, and editable external links per release
- 💬 Comment system (create, edit, delete)
- 🤖 AI summarizer (optional, requires `OPENAI_API_KEY`)
- ⚡ REST API built with **Go (Gin + GORM)**
- 🎨 Frontend built with **Next.js 14 (App Router)**

---

## 🏗️ Architecture
- **Frontend**: Next.js 14 (App Router), shadcn/ui
- **Backend**: Go (Gin + GORM)
- **Database**: MariaDB/MySQL
- **Migrations**: Goose + SQL files (embedded in backend)
- **API Proxy**: frontend → `/api/*` → backend

---

## 🚀 Installation

You can run Release Calendar either via **Helm (Kubernetes)** or **Docker Compose**.

### Option A — Helm (Kubernetes)

> ⚠️ The Helm chart **does not include a database**.  
> You must use your own MariaDB/MySQL instance and provide its connection settings.

Add the chart repository:

```bash
helm repo add thatqa https://thatqa.github.io/helm
helm repo update
helm search repo thatqa/release-calendar
```

Install:

```bash
helm install rc thatqa/release-calendar \
  --version 0.1.4 \
  --set backend.env.DB_HOST=mariadb.svc.cluster.local \
  --set backend.env.DB_PORT=3306 \
  --set backend.env.DB_NAME=release_calendar \
  --set backend.env.DB_USER=release_user \
  --set backend.env.DB_PASSWORD=secret
```

Enable AI summary (optional):
```bash
--set backend.env.AI_API_KEY=sk-your-key
--set backend.env.AI_TEMPERATURE=temp
--set backend.env.AI_MODEL=model
--set backend.env.AI_MAX_TOKENS=max_tokens
--set backend.env.AI_URL=url
```

Check resources:
```bash
kubectl get pods
kubectl get svc
kubectl get ingress
```
- Configure your Ingress host and TLS in chart values if needed.
- Ensure your DB is reachable from the cluster and the schema is migrated (use your preferred migration flow or the image’s migrate command if you run a Job).

### Option A — Helm (Kubernetes)

> 🐳 The Docker Compose setup includes a MariaDB database out of the box.    
> ⚠️ There are no health checks configured. To avoid race conditions, start services sequentially in the following order:
> 1. Database
> 2. Migrations
> 3. Backend
> 4. Frontend

Example docker-compose.yml:
```bash
version: "3.8"

services:
  db:
    image: mariadb:10.11
    restart: always
    environment:
      MARIADB_ROOT_PASSWORD: root
      MARIADB_DATABASE: release_calendar
      MARIADB_USER: release_user
      MARIADB_PASSWORD: secret
    ports:
      - "3306:3306"

  migrate:
    build: ./backend
    command: ["migrate"]
    environment:
      DB_HOST: db
      DB_PORT: 3306
      DB_NAME: release_calendar
      DB_USER: release_user
      DB_PASSWORD: secret
      DB_PARAMS: "charset=utf8mb4&parseTime=true&loc=UTC"
    depends_on:
      - db

  backend:
    build: ./backend
    command: ["serve"]
    environment:
      DB_HOST: db
      DB_PORT: 3306
      DB_NAME: release_calendar
      DB_USER: release_user
      DB_PASSWORD: secret
      DB_PARAMS: "charset=utf8mb4&parseTime=true&loc=UTC"
      # AI_API_KEY: "sk-your-key"   # optional
      # AI_TEMPERATURE: "temp"   # optional
      # AI_MODEL: "model"   # optional
      # AI_MAX_TOKENS: "max_tokens"   # optional
      # AI_URL: "url"   # optional
    ports:
      - "8080:8080"
    depends_on:
      - db

  frontend:
    build: ./frontend
    environment:
      NEXT_PUBLIC_API_BASE: /api
    ports:
      - "3000:3000"
    depends_on:
      - backend
```

Start sequentially:
```bash
docker compose up -d db
# wait until DB is accepting connections (e.g., 10–20s)

docker compose run --rm migrate
docker compose up -d backend frontend
```

Open → http://localhost:3000

## ⚙️ Configuration

### Backend Environment
| Variable         | Default                                            | Description                  |
|------------------|----------------------------------------------------|------------------------------|
| `DB_HOST`        | —                                                  | Database host                |
| `DB_PORT`        | 3306                                               | Database port                |
| `DB_NAME`        | —                                                  | Database name                |
| `DB_USER`        | —                                                  | Database user                |
| `DB_PASSWORD`    | —                                                  | Database password            |
| `DB_PARAMS`      | charset=utf8mb4&parseTime=true&loc=UTC             | DSN extra params             |
| `AI_TEMPERATURE` | empty                                              | Optional: enable AI summaries |
| `AI_API_KEY `    | empty                                              | Optional: enable AI summaries |
| `AI_MODEL `      | empty                                              | Optional: enable AI summaries |
| `AI_MAX_TOKENS ` | empty                                              | Optional: enable AI summaries |
| `AI_URL `        | empty                                              | Optional: enable AI summaries |

### Frontend Environment
| Variable               | Default | Description                          |
|------------------------|---------|--------------------------------------|
| `NEXT_PUBLIC_API_BASE` | /api    | API base path proxied by the frontend |

---

## 🔌 API Quick Reference

- `GET /api/releases` — list releases (supports filters: `date=YYYY-MM-DD`, `status`, `duty`)
- `POST /api/releases` — create a release
- `GET /api/releases/:id` — get release by id
- `PUT /api/releases/:id` — update a release
- `DELETE /api/releases/:id` — delete a release
- `GET /api/releases/:id/comments` — list comments
- `POST /api/releases/:id/comments` — add comment (`author`, `message`)
- `PUT /api/releases/:id/comments/:commentId` — update comment
- `DELETE /api/releases/:id/comments/:commentId` — delete comment
- `GET /api/releases/:id/summary` — AI summary (if `OPENAI_API_KEY` and other env present)
- `GET /api/release-days?from=YYYY-MM-DD&to=YYYY-MM-DD` — calendar status markers (per-day statuses in range)

### Examples:

```bash
# List releases for a day
curl -s "http://localhost:8080/api/releases?date=2025-09-18" | jq

# Create a release
curl -s -X POST http://localhost:8080/api/releases \
  -H "Content-Type: application/json" \
  -d '{
    "title":"Payments v3 rollout",
    "date":"2025-09-18T10:00:00",
    "status":"planned",
    "notes":"Blue/green deploy",
    "dutyUsers":["alice","bob"],
    "links":[{"name":"pipeline","url":"https://ci/p/123"}]
  }' | jq

# AI summary
curl -s "http://localhost:8080/api/releases/1/summary" | jq