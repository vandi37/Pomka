### **Запуск**:
Перед запуском, создайте файл .env в директории с docker-compose.yml
```
SERVER_NETWORK=tcp
SERVER_PORT=50123

DB_HOST=postgres
DB_PORT=5432
DB_MAX_ATMPS=5
DB_DELAY_ATMPS_S=5
DB_USER=<USER>
DB_PASSWORD=<PASSWORD>
DB_NAME=<NAME>

SERVICE_USERS_HOST=localhost
SERVICE_USERS_PORT=<PORT>
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
Перед тестами, создайте файл config.env в директории ./service/tests. Чтобы тесты смогли поднять бд, запустите докер. Тесты покрывают сервис Promos, описанный в proto/promos.
```
SERVER_NETWORK=tcp
SERVER_PORT=50123

DB_HOST=localhost
DB_PORT=5432
DB_MAX_ATMPS=5
DB_DELAY_ATMPS_S=5
DB_USER=<USER>
DB_PASSWORD=<PASSWORD>
DB_NAME=<NAME>

SERVICE_USERS_HOST=DONT_TOUCH
SERVICE_USERS_PORT=DONT_TOUCH
```
```
cd service/tests
go test -v
```
