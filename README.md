# Streamify

A modern, production-ready Go web application with PostgreSQL, JWT authentication, RESTful APIs, and a beautiful, developer-friendly UI/UX.

---

## 🚀 Features

- Clean Go backend (modular, testable, idiomatic)
- PostgreSQL with schema migrations (Goose)
- Secure JWT authentication & role-based access
- Modern REST API with Swagger (OpenAPI) docs
- Vendor support for reproducible builds
- Docker & Docker Compose for local and production
- Linting, formatting, and test coverage targets
- Ready for frontend integration (React, Next.js, etc.)

---

## 🏁 Quick Start (Docker)

```sh
git clone https://github.com/alireza-akbarzadeh/streamify.git
cd streamify
cp .env.example .env
# Edit .env for DB/JWT secrets

docker compose up -d  # Starts Postgres and pgAdmin
make migrate         # Run DB migrations
make build           # Build Go backend
make run             # Start API server
```

- API: http://localhost:8080
- Swagger UI: http://localhost:8080/swagger/index.html
- pgAdmin: http://localhost:5050

---

## 🛠️ Local Development

- Install Go 1.21+, Docker, and Goose CLI
- Set environment variables (see `.env.example`)
- Use `make dev` for hot reload (if using [air](https://github.com/cosmtrek/air) or similar)
- Use `make tidy-vendor` to sync dependencies

---

## 🗄️ Database Migrations

```sh
make migrate         # Apply all migrations
make migrate-down    # Roll back last migration
make generate        # Regenerate SQL code (sqlc)
```

---

## 🔒 Environment Variables

- `DB_URL` - PostgreSQL connection string
- `JWT_SECRET` - Required! Used for signing tokens
- `FRONTEND_URL` - CORS and redirect support
- `PORT` - API server port (default: 8080)

---

## 📖 API Documentation

- Run `make swagger` to generate docs
- Access Swagger UI at `/swagger/index.html`
- Example endpoints:
  - `POST /api/v1/auth/register` (User registration)
  - `POST /api/v1/auth/login` (Login, returns JWT)
  - `GET  /api/v1/users` (List users, admin only)

---

## 🧑‍💻 Developer Experience

- Lint: `golangci-lint run`
- Test: `make test` or `make test-vendor`
- Coverage: `make test-cover-html`
- Format: `gofmt -w .`
- VS Code: Recommended extensions in `.vscode/extensions.json`

---

## 🖥️ Modern UI/UX

- Designed for easy frontend integration (React, Next.js, Vue, etc.)
- CORS enabled for local and production
- API returns clear, consistent JSON errors and responses
- Ready for OAuth, SSO, and social login extensions

---

## 🚀 Production & Deployment

- Use `GOFLAGS=-mod=vendor` for reproducible builds
- Set all secrets via environment variables (never commit secrets)
- Use Docker Compose or Kubernetes for deployment
- Monitor with Prometheus, Grafana, and structured logs
- Regularly update dependencies and run security scans

---

## 🤝 Contributing

- Fork, branch, and submit PRs
- Write tests for new features
- Follow code style and commit guidelines

---

## 📄 License

MIT License. See [LICENSE](LICENSE) for details.

---

## 💡 Need help?

Open an issue or contact the maintainer at [devtools95@gmail.com](mailto:devtools95@gmail.com).
