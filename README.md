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
```mermaid
flowchart LR
  A[Frontend (Next.js)] -->|HTTP/Ingress| B[Backend (Go + Gin + GORM)]
  B <--> C[(MariaDB/MySQL)]
  A -.->|/api proxy| B
```

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

---

Option B — Docker Compose

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