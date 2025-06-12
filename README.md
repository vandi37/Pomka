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
Перед запуском, создайте файл config.env в директории ./service/tests
```
SERVER_NETWORK=tcp
SERVER_PORT=50123

DB_HOST=
DB_PORT=
DB_MAX_ATMPS=
DB_DELAY_ATMPS_S=
DB_USER=
DB_PASSWORD=
DB_NAME=

SERVICE_USERS_HOST=
SERVICE_USERS_PORT=
```
```
cd service/tests
go test -v
```
