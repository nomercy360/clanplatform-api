# Go api server for clanplatform app

## How to run

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate -source file://${PWD}/migrations -database postgres://postgres:mysecretpassword@localhost:5432/clanplatform\?sslmode=disable up
```

```bash
cp configs/config.api.example.yaml config.yaml
go run main.go
```