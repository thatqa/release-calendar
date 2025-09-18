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