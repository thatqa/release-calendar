# Release Calendar

Полностью готовое приложение по ТЗ:
- **Frontend**: Next.js 14 (App Router) + Tailwind + лёгкие компоненты в стиле shadcn/ui.
- **Backend**: Go (Gin) + GORM + goose миграции (встроены через `go:embed`).
- **DB**: MariaDB.
- **Infra**: Nginx как reverse-proxy, всё упаковано в Docker + docker-compose.

## Запуск
```bash
cp .env.example .env
docker-compose up --build
```
Открой: `http://localhost:8088`

## API
- `GET    /api/releases?date=YYYY-MM-DD&status=&duty=`
- `POST   /api/releases` (body: {title, date (ISO), status, notes, dutyUsers[], links[]})
- `GET    /api/releases/:id`
- `PUT    /api/releases/:id` (обновляет и синхронизирует links)
- `DELETE /api/releases/:id`

Комментарии:
- `GET    /api/releases/:id/comments`
- `POST   /api/releases/:id/comments` ({author, message})
- `PUT    /api/releases/:id/comments/:commentId`
- `DELETE /api/releases/:id/comments/:commentId`

## Заметки
- Поле dutyUsers хранится в JSON; фильтрация по duty происходит на приложении при листинге.
- Links управляются только через POST/PUT release, как ты просил.
- UI: календарь слева, справа карточка релиза + CRUD.
