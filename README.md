# Industry-Level Go Web Server

> **One project. Nine stages. From main.go to something engineers at Google, Uber, and Docker would actually respect.**

##  Project Overview

This is a production-grade **Users API microservice** built with Go, demonstrating industry-standard patterns used at companies like Google, Uber, Twitter, and Docker. Every architectural decision is intentional and battle-tested at scale.

### Technology Stack

| Component | Technology | Why This Choice |
|-----------|-----------|-----------------|
| **Language** | Go 1.21+ | Concurrency primitives, fast compilation, used by Docker/Kubernetes |
| **Database** | PostgreSQL | Superior JSON support, MVCC concurrency, ACID transactions (used by Uber, Instagram) |
| **Cache** | Redis | Rich data structures, pub/sub, persistence (used by Twitter, GitHub) |
| **HTTP Router** | net/http stdlib | Understanding primitives before frameworks; production-ready |
| **Logging** | slog (stdlib) | Structured JSON logging, zero dependencies |
| **Metrics** | Prometheus | Industry standard for Kubernetes environments |
| **Auth** | JWT | Stateless, microservices-friendly (used by Netflix, Uber) |

---

##  System Architecture

### High-Level Request Flow

`mermaid
graph TD
    A[HTTP Client / curl / Frontend] -->|HTTPS Request| B[Middleware Stack]
    B --> C[Rate Limiter]
    C --> D[Auth Middleware JWT]
    D --> E[Logger Middleware]
    E --> F[HTTP Router net/http]
    F --> G[Handler Layer]
    G --> H[Service Layer Business Logic]
    H --> I{Cache Check}
    I -->|Cache Hit| J[Redis Read Cache]
    I -->|Cache Miss| K[Repository Layer]
    K --> L[PostgreSQL Write + Read]
    L --> M[Update Cache]
    M --> J
    J --> N[Response]
    L -.->|metrics exported| O[Prometheus /metrics]
    
    style A fill:#e1f5ff
    style B fill:#fff4e1
    style G fill:#ffe1f5
    style H fill:#e1ffe1
    style J fill:#ff6b6b
    style L fill:#4ecdc4
    style O fill:#ffe66d
`

### Why This Architecture?

**Separation of Concerns** means:
-  Swap PostgreSQL for MongoDB without touching handlers
-  Add Redis caching without touching business logic
-  Test each layer in isolation
-  Scale components independently

---

##  Project Structure

`
go-industry-server/
 cmd/
    server/
        main.go               Binary entry point only
 internal/                     Cannot be imported externally (Go enforced)
    handler/                  HTTP handlers (parse req, call service, write resp)
       user.go
       auth.go
    service/                  Business logic lives here
       user.go
    repository/               DB queries only
       user.go
       postgres_user.go
    middleware/               Logging, auth, rate limit, recovery
       middleware.go
       auth.go
       ratelimit.go
    model/                    Shared structs
       user.go
    auth/                     JWT service
       jwt.go
    cache/                    Redis abstraction
        redis.go
 configs/
    config.go                 Environment-driven configuration
 pkg/                          Reusable code safe to export
    response/
        response.go           Standard API response envelope
 migrations/                   Database migrations
    001_create_users.up.sql
    001_create_users.down.sql
 docker-compose.yml            Local development orchestration
 Dockerfile                    Multi-stage production build
 go.mod
 go.sum
`

### Why internal/ Matters

Go **enforces** that packages inside internal/ cannot be imported by code outside the module. This is a **compile-time boundary** that prevents accidental coupling between services. Google and Uber enforce this in their Go monorepos.

