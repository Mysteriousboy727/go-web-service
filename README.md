# ⚡ GoForge — Production API Blueprint

> **One project. Nine stages. From `main.go` to something engineers at Google, Uber, and Docker would actually respect.**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-Primary_DB-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://postgresql.org)
[![Redis](https://img.shields.io/badge/Redis-Cache_Layer-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io)
[![Prometheus](https://img.shields.io/badge/Prometheus-Observability-E6522C?style=for-the-badge&logo=prometheus&logoColor=white)](https://prometheus.io)
[![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://docker.com)
[![JWT](https://img.shields.io/badge/JWT-Auth-black?style=for-the-badge&logo=jsonwebtokens&logoColor=white)](https://jwt.io)

---

## 🗺️ System Architecture — Full Request Flow

```mermaid
flowchart TD
    CLIENT(["🌐 HTTP Client / curl / Frontend"]):::client

    subgraph MW ["⚡ Middleware Stack"]
        direction LR
        RL["🛡️ Rate Limiter\nToken Bucket"]:::mw_rl
        AUTH["🔐 JWT Auth\nValidation"]:::mw_auth
        LOG["📋 Logger\nStructured slog"]:::mw_log
    end

    ROUTER["🔀 HTTP Router\nnet/http + ServeMux"]:::router

    subgraph APP ["📦 Application Layers"]
        direction TB
        HANDLER["🎯 Handler Layer\nParse Request → Call Service → Respond"]:::handler
        SERVICE["⚙️ Service Layer\nBusiness Logic & Validation"]:::service
        REPO["🗄️ Repository Layer\nDatabase Queries Only"]:::repo
    end

    subgraph STORE ["💾 Storage"]
        direction LR
        PG[("🐘 PostgreSQL\nWrite + Read")]:::postgres
        REDIS[("⚡ Redis\nRead Cache")]:::redis
    end

    PROM["📊 Prometheus\n/metrics"]:::prometheus

    CLIENT -->|"HTTPS Request"| MW
    RL --> AUTH --> LOG
    MW -->|"Validated Request"| ROUTER
    ROUTER --> HANDLER
    HANDLER --> SERVICE
    SERVICE --> REPO
    REPO -->|"Cache Miss"| PG
    REPO -->|"Cache Hit ⚡"| REDIS
    PG -.->|"Populate Cache"| REDIS
    PG -.->|"metrics"| PROM
    REDIS -.->|"metrics"| PROM

    classDef client fill:#1a1a2e,stroke:#e94560,stroke-width:3px,color:#e94560,font-weight:bold
    classDef mw_rl fill:#16213e,stroke:#f5a623,stroke-width:2px,color:#f5a623
    classDef mw_auth fill:#16213e,stroke:#7bed9f,stroke-width:2px,color:#7bed9f
    classDef mw_log fill:#16213e,stroke:#70a1ff,stroke-width:2px,color:#70a1ff
    classDef router fill:#0f3460,stroke:#e94560,stroke-width:2px,color:#ffffff,font-weight:bold
    classDef handler fill:#533483,stroke:#c77dff,stroke-width:2px,color:#c77dff
    classDef service fill:#1b4332,stroke:#52b788,stroke-width:2px,color:#52b788
    classDef repo fill:#1c3144,stroke:#48cae4,stroke-width:2px,color:#48cae4
    classDef postgres fill:#003049,stroke:#4ecdc4,stroke-width:3px,color:#4ecdc4,font-weight:bold
    classDef redis fill:#3d0000,stroke:#ff6b6b,stroke-width:3px,color:#ff6b6b,font-weight:bold
    classDef prometheus fill:#2d1b00,stroke:#ffd166,stroke-width:2px,color:#ffd166
```

---

## 🏗️ Nine Stages — Build Progression

```mermaid
flowchart LR
    S1["1\n🌱 Basic HTTP\nServer\n---\nnet/http\nslog\nMiddleware"]:::s1
    S2["2\n🐘 CRUD\nPostgreSQL\n---\npgx/v5\nMigrations\nConn Pool"]:::s2
    S3["3\n🔐 JWT\nAuth\n---\nbcrypt\nToken Pairs\nMiddleware"]:::s3
    S4["4\n⚡ Redis\nCaching\n---\nCache-Aside\nTTL\nInvalidation"]:::s4
    S5["5\n🛡️ Rate\nLimiting\n---\nToken Bucket\nPer-IP\nRedis Sync"]:::s5
    S6["6\n🐳 Docker\nCompose\n---\nMulti-Stage\nOrchestration\nDev Env"]:::s6
    S7["7\n🪂 Graceful\nShutdown\n---\nos.Signal\nContext\nDrain"]:::s7
    S8["8\n🧪 Testing\n---\nUnit Tests\nIntegration\nMock Repos"]:::s8
    S9["9\n📊 Prometheus\nMetrics\n---\n/metrics\nHistograms\nGauges"]:::s9

    S1 --> S2 --> S3 --> S4 --> S5 --> S6 --> S7 --> S8 --> S9

    classDef s1 fill:#2d6a4f,stroke:#52b788,stroke-width:2px,color:#b7e4c7
    classDef s2 fill:#1d3557,stroke:#457b9d,stroke-width:2px,color:#a8dadc
    classDef s3 fill:#4a1942,stroke:#c77dff,stroke-width:2px,color:#e0aaff
    classDef s4 fill:#641220,stroke:#e63946,stroke-width:2px,color:#ffb3c1
    classDef s5 fill:#7b2d00,stroke:#f4a261,stroke-width:2px,color:#ffd8b1
    classDef s6 fill:#023e8a,stroke:#0096c7,stroke-width:2px,color:#90e0ef
    classDef s7 fill:#283618,stroke:#606c38,stroke-width:2px,color:#dda15e
    classDef s8 fill:#370617,stroke:#f48c06,stroke-width:2px,color:#faa307
    classDef s9 fill:#3a0ca3,stroke:#7209b7,stroke-width:2px,color:#f72585
```

---

## 📐 Dependency Injection and Layer Flow

```mermaid
flowchart TD
    subgraph ENTRY ["cmd/server/main.go"]
        MAIN["🚀 main()"]:::main_node
    end

    subgraph CFG ["configs/"]
        CONFIG["⚙️ Config\nEnv Variables\n12-Factor App"]:::config_node
    end

    subgraph INT ["internal/  —  compile-time boundary"]
        subgraph HND ["handler/"]
            UH["🎯 UserHandler"]:::h_node
            AH["🔑 AuthHandler"]:::h_node
        end
        subgraph SVC ["service/"]
            US["⚙️ UserService interface"]:::svc_node
        end
        subgraph REP ["repository/"]
            UR["📋 UserRepository\ninterface"]:::repo_iface
            IMR["🧠 InMemory\nStage 1"]:::repo_impl
            PGR["🐘 Postgres\nStage 2+"]:::repo_impl
        end
        subgraph MWS ["middleware/"]
            M1["🛡️ RequestID"]:::mw
            M2["📋 Logger"]:::mw
            M3["🔐 Auth JWT"]:::mw
            M4["⏱️ RateLimit"]:::mw
            M5["💥 Recovery"]:::mw
        end
        subgraph AUT ["auth/"]
            JWT["🔏 JWTService\nHS256"]:::jwt_node
        end
        subgraph CAC ["cache/"]
            RC["⚡ RedisCache\nCache-Aside"]:::cache_node
        end
    end

    subgraph PKG ["pkg/"]
        RESP["📦 APIResponse\nEnvelope"]:::pkg_node
    end

    subgraph INFRA ["Infrastructure"]
        PGDB[("🐘 PostgreSQL")]:::db_node
        RDB[("⚡ Redis")]:::redis_node
        PROM["📊 Prometheus"]:::prom_node
    end

    MAIN --> CONFIG
    MAIN --> UH & AH
    MAIN --> M1 & M2 & M3 & M4 & M5
    UH --> US
    AH --> JWT
    US --> UR & RC
    UR -.->|"Stage 1"| IMR
    UR -.->|"Stage 2+"| PGR
    PGR --> PGDB
    RC --> RDB
    UH & AH --> RESP
    PGDB & RDB -.-> PROM

    classDef main_node fill:#e63946,stroke:#ff6b6b,color:#fff,font-weight:bold
    classDef config_node fill:#457b9d,stroke:#a8dadc,color:#fff
    classDef h_node fill:#7b2d8b,stroke:#c77dff,color:#e0aaff
    classDef svc_node fill:#2d6a4f,stroke:#52b788,color:#b7e4c7
    classDef repo_iface fill:#1d3557,stroke:#457b9d,color:#a8dadc,stroke-dasharray:5 5
    classDef repo_impl fill:#023e8a,stroke:#0096c7,color:#90e0ef
    classDef mw fill:#7b2d00,stroke:#f4a261,color:#ffd8b1
    classDef jwt_node fill:#3a0ca3,stroke:#7209b7,color:#f72585
    classDef cache_node fill:#641220,stroke:#e63946,color:#ffb3c1
    classDef pkg_node fill:#370617,stroke:#f48c06,color:#faa307
    classDef db_node fill:#003049,stroke:#4ecdc4,color:#4ecdc4,font-weight:bold
    classDef redis_node fill:#3d0000,stroke:#ff6b6b,color:#ff6b6b,font-weight:bold
    classDef prom_node fill:#2d1b00,stroke:#ffd166,color:#ffd166
```

---

## 🗄️ Read vs Write Data Flow

```mermaid
sequenceDiagram
    participant C as 🌐 Client
    participant M as ⚡ Middleware
    participant H as 🎯 Handler
    participant S as ⚙️ Service
    participant RC as ⚡ Redis
    participant PG as 🐘 PostgreSQL

    rect rgb(20, 40, 20)
        Note over C,PG: GET /api/v1/users/:id — Cache-Aside Read
        C->>+M: GET /users/abc Bearer token
        M->>M: Validate JWT + Rate Limit
        M->>+H: Authorized Request
        H->>+S: GetUser("abc")
        S->>+RC: GET user:abc
        alt Cache Hit (~0.1ms)
            RC-->>S: User JSON
            S-->>H: model.User
        else Cache Miss (~10ms)
            RC-->>S: redis.Nil
            S->>+PG: SELECT FROM users WHERE id=$1
            PG-->>-S: Row Data
            S->>RC: SET user:abc TTL=5min
            S-->>H: model.User
        end
        H-->>-C: 200 success true data
    end

    rect rgb(40, 20, 20)
        Note over C,PG: PUT /api/v1/users/:id — Write + Cache Invalidate
        C->>+M: PUT /users/abc Bearer token
        M->>+H: Authorized Request
        H->>+S: UpdateUser("abc", req)
        S->>+PG: UPDATE users SET name=$1 WHERE id=$2
        PG-->>-S: RowsAffected 1
        S->>RC: DEL user:abc
        Note over RC: Cache Invalidated
        S-->>H: Updated model.User
        H-->>-C: 200 success true data
    end
```

---

## 🛡️ Token Bucket Rate Limiter

```mermaid
flowchart LR
    REQ(["📨 Incoming Request\nfrom IP x.x.x.x"]):::req

    subgraph BUCKET ["Token Bucket — per IP address"]
        direction TB
        T1["🪙 token"]:::token
        T2["🪙 token"]:::token
        T3["🪙 token"]:::token
        T4["🪙 token"]:::token
        T5["🪙 token"]:::token
        FILL(["♻️ +r tokens/sec\nup to burst b"]):::fill
    end

    ALLOW(["✅ 200 OK\nConsume 1 token\nForward to Handler"]):::allow
    DENY(["❌ 429 Too Many Requests\nRetry-After: 1\nDrop Request"]):::deny

    REQ --> BUCKET
    BUCKET -->|"Token available"| ALLOW
    BUCKET -->|"Bucket empty"| DENY
    FILL --> BUCKET

    classDef req fill:#1d3557,stroke:#457b9d,color:#a8dadc,font-weight:bold
    classDef token fill:#f4a261,stroke:#e76f51,color:#1a1a1a,font-weight:bold
    classDef fill fill:#2d6a4f,stroke:#52b788,color:#b7e4c7
    classDef allow fill:#1b4332,stroke:#52b788,color:#b7e4c7,font-weight:bold
    classDef deny fill:#641220,stroke:#e63946,color:#ffb3c1,font-weight:bold
```

---

## 🔐 JWT Authentication Flow

```mermaid
sequenceDiagram
    participant C as 🌐 Client
    participant A as 🔑 Auth Handler
    participant S as ⚙️ Service
    participant DB as 🐘 PostgreSQL
    participant J as 🔏 JWT Service

    rect rgb(20, 20, 45)
        Note over C,J: POST /api/v1/auth/login
        C->>+A: email + password
        A->>+S: Login(email, password)
        S->>+DB: SELECT WHERE email=$1
        DB-->>-S: User row with hashed password
        S->>S: bcrypt.Compare cost=12 ~250ms
        alt Valid credentials
            S->>+J: GenerateTokenPair(userID, email)
            J-->>-S: accessToken 15min + refreshToken 7d
            A-->>C: 200 access_token refresh_token
        else Invalid credentials
            A-->>-C: 401 invalid credentials
            Note over A,C: Same error for wrong email OR wrong password
        end
    end

    rect rgb(20, 35, 20)
        Note over C,J: Subsequent Protected Request
        C->>+A: GET /users/abc Authorization Bearer token
        A->>+J: ValidateToken(tokenStr)
        J->>J: Verify HS256 signature
        J->>J: Check expiry + signing method
        J-->>-A: Claims UserID Email
        A->>A: Inject UserID into context
        A-->>-C: Forward to Handler
    end
```

---

## 📊 Technology Stack

| Component | Technology | Why |
|-----------|-----------|-----|
| 🐹 **Language** | Go 1.21+ | Concurrency primitives, fast compilation, Docker/K8s native |
| 🐘 **Database** | PostgreSQL + pgx/v5 | MVCC, 2–3x faster than lib/pq, used at Uber and Cloudflare |
| ⚡ **Cache** | Redis + go-redis/v9 | TTL, pub/sub, rich types — used at Twitter and GitHub |
| 🔀 **HTTP Router** | net/http stdlib | Learn primitives first; production-ready without a framework |
| 📋 **Logging** | slog (Go 1.21 stdlib) | Structured JSON → Fluentd / Loki / Grafana pipeline |
| 📊 **Metrics** | Prometheus | Standard for Kubernetes; pull-based; open-source |
| 🔐 **Auth** | JWT (golang-jwt/v5) | Stateless, no shared session store — microservices-friendly |
| 🔒 **Passwords** | bcrypt cost=12 | ~250ms per check — brute-force infeasible |
| 🐳 **Containers** | Docker multi-stage | CGO_ENABLED=0 static binary; tiny runtime image |

---

## 📁 Project Structure

```
go-industry-server/
├── cmd/
│   └── server/
│       └── main.go               ← binary entry point only
├── internal/                     ← Go enforced: no external imports allowed
│   ├── handler/
│   │   ├── user.go
│   │   └── auth.go
│   ├── service/
│   │   └── user.go
│   ├── repository/
│   │   ├── user.go               ← InMemory (Stage 1)
│   │   └── postgres_user.go      ← PostgreSQL (Stage 2+)
│   ├── middleware/
│   │   ├── middleware.go         ← RequestID, Logger, Recovery, Chain
│   │   ├── auth.go               ← JWT middleware
│   │   └── ratelimit.go          ← token bucket per IP
│   ├── model/
│   │   └── user.go
│   ├── auth/
│   │   └── jwt.go
│   └── cache/
│       └── redis.go
├── configs/
│   └── config.go                 ← 12-Factor env config
├── pkg/
│   └── response/
│       └── response.go           ← standard APIResponse envelope
├── migrations/
│   ├── 001_create_users.up.sql
│   └── 001_create_users.down.sql
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── go.sum
```

> **Why `internal/` matters** — Go enforces that packages inside `internal/` cannot be imported by code outside the module. This is a compile-time boundary that prevents accidental coupling. Google and Uber enforce this in their Go monorepos.

---

## ⚡ Quick Start

```bash
git clone https://github.com/yourname/go-industry-server
cd go-industry-server
go mod download

# Spin up PostgreSQL + Redis
docker-compose up -d

# Run the server
go run ./cmd/server/main.go

# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com","password":"secret123"}'

# Health check
curl http://localhost:8080/health
```

---

## 🌱 Stage Summary

| Stage | What You Build | Key Concepts |
|-------|---------------|-------------|
| 1 🌱 | Basic HTTP Server | Clean arch, DI, middleware chain, slog, `internal/` |
| 2 🐘 | CRUD + PostgreSQL | pgx/v5, migrations, connection pooling |
| 3 🔐 | JWT Auth | bcrypt cost=12, token pairs, algorithm confusion prevention |
| 4 ⚡ | Redis Caching | Cache-aside, TTL strategy, invalidation on write |
| 5 🛡️ | Rate Limiting | Token bucket, per-IP, Redis-backed for multi-pod |
| 6 🐳 | Docker + Compose | Multi-stage builds, static binaries, dev orchestration |
| 7 🪂 | Graceful Shutdown | `os.Signal`, context cancel, drain in-flight requests |
| 8 🧪 | Testing | Unit + integration, mock repos, testcontainers |
| 9 📊 | Prometheus Metrics | `/metrics`, histograms for latency, gauges for pool size |

---

**Built with ❤️ using Go — the language of Docker, Kubernetes, and Uber.**

*Every architectural decision mirrors what's running in production at scale.*
