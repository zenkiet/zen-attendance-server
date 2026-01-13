golang-project/
├── cmd/                    # Entry points
│   └── api/
│       └── main.go
├── internal/               # Main Code (Private)
│   ├── handlers/           # HTTP Handlers (Controllers)
│   ├── models/             # Business Objects
│   ├── repository/         # DB Access Layer
│   └── service/            # Business Logic
├── pkg/                    # Library code (Public)
│   ├── utils/
│   ├── databases/
│   ├── migrations/
├── configs/                # (yaml, .env)
├── deployments/            # Dockerfile, docker-compose, k8s
├── scripts/                # support for building
├── go.mod
├── go.sum
├── Makefile                # Task runner
├── .golangci.yml           # Config Linter
└── lefthook.yml            # Config Git Hooks