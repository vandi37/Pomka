### **Запуск**:
Перед запуском, заполните файл .env в директории с docker-compose.yml
```
docker compose up
```

### **Миграции**:
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

### **Тесты**:
Перед тестами, заполните файл config.env в директории ./service/tests. Чтобы тесты смогли поднять бд, запустите докер.
```
cd service/tests
go test -v
```
