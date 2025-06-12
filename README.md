### **Запуск**:
Перед запуском, создайте файл .env в директории с docker-compose.yml
```
SERVER_NETWORK=tcp
SERVER_PORT=50123

POSTGRES_HOST=
POSTGRES_PORT=
DB_MAX_ATMPS=
DB_DELAY_ATMPS_S=

POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=

SERVICE_USERS_HOST=localhost
SERVICE_USERS_PORT=
```

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

POSTGRES_HOST=
POSTGRES_PORT=
DB_MAX_ATMPS=
DB_DELAY_ATMPS_S=

POSTGRES_USER=
POSTGRES_PASSWORD=
POSTGRES_DB=

SERVICE_USERS_HOST=localhost
SERVICE_USERS_PORT=
```
```
cd service/tests
go test -v
```
