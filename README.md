запуск:
```
docker compose up
```

миграции:
```
cd service/internal/migrations
```
Up
```
goose postgres postgresql://USER:PASSWORD@HOST:PORT/postgres up
```
Down
```
goose postgres postgresql://USER:PASSWORD@HOST:PORT/postgres down
```

тесты:
```
cd service/tests
go test -v
```
