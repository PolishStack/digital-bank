# Project Structure Example

```plaintext
bitka/
├── go.work                     # Workspace linking all modules
├── README.md
├── Makefile                    # build, test, lint, etc.
├── deploy/                     # DevOps & infra
│   ├── k8s/                    # All Kubernetes manifests (by service)
│   │   ├── base/
│   │   └── overlays/
│   ├── terraform/              # IaC for AWS + on-prem
│   └── github-actions/         # CI/CD pipelines
├── services/                   # All microservices
│   ├── auth/
│   │   ├── go.mod              # module bitka/auth-service
│   │   ├── cmd/server/         # main.go here
│   │   ├── internal/
│   │   │   ├── config/         # service config
│   │   │   ├── domain/         # domain models + business rules
│   │   │   ├── usecase/        # application logic
│   │   │   ├── repo/           # interfaces for data layer
│   │   │   └── transport/      # HTTP/GRPC handlers
│   │   ├── pkg/                # service-only helpers (not shared)
│   │   └── migrations/
│   │
│   ├── wallet/
│   │   ├── go.mod              # module bitka/wallet
│   │   └── ... same layout ...
│   │
│   └── ledger/
│       ├── go.mod              # module bitka/ledger
│       └── ... same layout ...
├── libs/                       # Shared libraries (reusable across services)
│   ├── logger/
│   │   ├── go.mod              # module bitka/logger
│   │   └── logger.go           # zerolog-based
│   │
│   ├── jwt/
│   │   ├── go.mod              # module bitka/jwt
│   │   └── jwt.go              # common JWT utilities
│   │
│   ├── db/
│   │   ├── go.mod              # module bitka/db
│   │   └── postgres.go         # pgx connection + helpers
│   │
│   ├── config/
│   │   ├── go.mod              # module bitka/config
│   │   └── config.go           # env loader + struct mapping
│   │
│   └── common/
│       ├── go.mod              # module bitka/common
│       └── errors.go           # shared custom errors (optional)
└── tools/                      # Dev tools (lint configs, codegen, scripts)
    ├── scripts/
    └── lint/
```