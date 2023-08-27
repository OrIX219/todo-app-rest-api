# TODO lists app

REST API application for keeping track of TODO tasks organized in lists.

Requires registration with login and password. Logging in returns JWT token
which is used for later authentification for API requests.

## Prerequisites
- __Docker__ with _compose_ plugin

## Run
1. `make env`
2. `docker compose up -d`

## Technology stack
- Go ([gin](https://github.com/gin-gonic/gin),
      [sqlx](https://github.com/jmoiron/sqlx),
      [goose](https://github.com/pressly/goose),
      [testify](https://github.com/stretchr/testify))
- PostgreSQL
- Docker
