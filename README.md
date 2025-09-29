# SalesTracker

SalesTracker is a backend + frontend application for managing items (e.g., financial transactions, orders, events) and performing analytics on them. The service supports full CRUD operations and provides aggregated analytics such as sum, average, count, median, and percentiles.

All data is stored in PostgreSQL. Analytics queries are performed using SQL, including window functions for median and percentile calculations. The frontend is a simple HTML/JS interface (React/TS/Tailwind optional) to manage and visualize items and categories.

---

## Features

* **CRUD operations for items and categories**
* **Analytics endpoints**:

    * Sum (`/analytics/sum`)
    * Average (`/analytics/avg`)
    * Count (`/analytics/count`)
    * Median (`/analytics/median`)
    * Percentile (`/analytics/percentile`)
* **Frontend** for managing items and categories
* **Filter items by date, category, and kind**
* **Validation** of amounts, dates, and JSON metadata
* **PostgreSQL database with proper indexing for analytics**

---

## API Endpoints

### Categories

| Method | Endpoint              | Description           |
| ------ | --------------------- | --------------------- |
| POST   | `/api/categories`     | Create a new category |
| GET    | `/api/categories`     | List all categories   |
| GET    | `/api/categories/:id` | Get category by ID    |
| PUT    | `/api/categories/:id` | Update category by ID |
| DELETE | `/api/categories/:id` | Delete category by ID |

### Items

| Method | Endpoint         | Description                                                         |
| ------ | ---------------- | ------------------------------------------------------------------- |
| POST   | `/api/items`     | Create a new item                                                   |
| GET    | `/api/items`     | List all items (with optional filters: from, to, category_id, kind) |
| GET    | `/api/items/:id` | Get item by ID                                                      |
| PUT    | `/api/items/:id` | Update item by ID                                                   |
| DELETE | `/api/items/:id` | Delete item by ID                                                   |

### Analytics

| Method | Endpoint                    | Description                                   |
| ------ | --------------------------- | --------------------------------------------- |
| GET    | `/api/analytics/sum`        | Get sum of items in a period                  |
| GET    | `/api/analytics/avg`        | Get average amount                            |
| GET    | `/api/analytics/count`      | Get count of items                            |
| GET    | `/api/analytics/median`     | Get median amount                             |
| GET    | `/api/analytics/percentile` | Get N-th percentile (query: `percentile=0.9`) |

**Query parameters for analytics endpoints:**

* `from` (optional): start date (ISO8601 / RFC3339)
* `to` (optional): end date (ISO8601 / RFC3339)
* `category_id` (optional): filter by category UUID
* `kind` (optional): filter by item kind (`income`, `expense`, `transfer`, `refund`)
* `percentile` (optional, default 0.9): for percentile endpoint

---

## Project Structure

```
.
├── cmd/                 # Application entry points
├── config/              # Configuration files
├── internal/            # Internal application packages
│   ├── api/             # HTTP handlers, router, server
│   │   ├── request      # Helpers (ParseUUIDParam, ParseUUIDQuery, ParseFloatQuery, ParseTimeQuery etc.)
│   │   ├── response     # Response helpers (JSON, OK, Created etc.)
│   │   ├── router
│   │   └── server
│   ├── config/          # Config parsing logic
│   ├── model/           # Data models
│   ├── repository/      # Database repositories
│   └── service/         # Business logic
├── migrations/          # Database migrations
├── web/                 # Frontend UI (React + TS + TailwindCSS)
├── Dockerfile           # Backend Dockerfile
├── go.mod
├── go.sum
├── .env.example         # Example environment variables
├── docker-compose.yml   # Multi-service Docker setup
├── Makefile             # Development commands
└── README.md
```

---

## Running the project

### Using Makefile

* Build and start all Docker services:

```bash
make docker-up
```

* Stop and remove all Docker services and volumes:

```bash
make docker-down
```

### Default Ports

* Frontend: [http://localhost:3000](http://localhost:3000)
* Backend API: [http://localhost:8080/api](http://localhost:8080/api)

---

## Usage

* Use the frontend to:

    * Add new items and categories
    * View and filter items
    * Fetch analytics (sum, average, count, median, percentile)
* Use the API directly for automation or integration:

    * POST, GET, PUT, DELETE items and categories
    * Query analytics with optional filters (`from`, `to`, `category_id`, `kind`, `percentile`)

---

## Tech Stack

* **Backend:** Go, Gin, PostgreSQL
* **Database:** PostgreSQL
* **Docker:** Multi-service setup (backend + frontend + database)

---

## Notes

* All amounts are stored as decimal (`NUMERIC`) and validated to be non-negative.
* Item metadata is stored as JSONB and can hold any arbitrary JSON object.
* Analytics queries are performed in SQL with proper indexing for efficiency.
* Date and time filters should use ISO8601/RFC3339 format.