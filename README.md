# Лабораторная работа №11: Контейнеризация мультиязычных приложений

**Студент:** Ражина Маргарита Александровна
**Группа:** 220032-11
**Вариант:** 8

## Выполненные задания

### Средняя сложность:
1. **Задание 2:** Написать Dockerfile для Go-приложения с многоэтапной сборкой.

**Эндпоинты:**

- `GET /health` — проверка статуса
- `POST /echo` — возврат принятого JSON
- `POST /user` — создание пользователя (name required, age 18-100)
- `GET /swagger/index.html` — [Swagger UI](http://localhost:8080/swagger/index.html)

**Запуск без Docker:**
```bash
cd task1_2
go mod tidy
go run main.go
```

**Swagger UI:** http://localhost:8080/swagger/index.html

**Запуск тестов:**
```bash
cd task1_2
go test 
```

**Сборка Docker-образа:**
```bash
cd task1_2
docker build -t lab11-go .
```

**Запуск контейнера:**
```bash
docker run -d --name go-app -p 8080:8080 lab11-go
```

**Проверка статуса:**
```bash
docker ps
```

**Просмотр логов:**
```bash
docker logs go-app
```

**Остановка контейнера:**
```bash
docker stop go-app
```

**Удаление контейнера:**
```bash
docker rm go-app
```

**Удаление образа:**
```bash
docker rmi lab11-go
```

**Просмотр размера образа:**
```bash
docker images lab11-go
```


2. **Задание 8:** Добавить healthcheck для каждого сервиса.


**Эндпоинты:**

- `GET /health` — проверка статуса
- `GET /` — приветствие
- `POST /data` — эхо-ответ для JSON


**Запуск тестов:**

Go:
```bash
cd task2_8/go
go test -v ./...
```

Python:
```bash
cd task2_8/python
python -m pytest test_app.py -v
```

Rust:
```bash
cd task2_8/rust
cargo test
```

**Сборка образов:**

```bash
cd task2_8/go
docker build -t lab11-go-task2 .

cd task2_8/python
docker build -t lab11-python-task2 .

cd task2_8/rust
docker build -t lab11-rust-task2 .
```

**Запуск через Docker Compose:**

```bash
cd task2_8
docker compose up -d
```

**Проверка эндпоинтов:**

```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
curl http://localhost:8083/health
```

**Проверка статуса сервисов:**

```bash
docker compose -f task2_8/docker-compose.yml ps
```

**Просмотр логов:**

```bash
docker compose -f task2_8/docker-compose.yml logs go-service
docker compose -f task2_8/docker-compose.yml logs python-service
docker compose -f task2_8/docker-compose.yml logs rust-service
```

**Остановка:**

```bash
cd task2_8
docker compose down
```

**Остановка с удалением образов:**

```bash
cd task2_8
docker compose down --rmi all
```

**Просмотр размеров образов:**

```bash
docker image ls --filter "reference=lab11-*-task2"
```

3. **Задание 10:** Использовать переменные окружения для конфигурации.


**Эндпоинты:**

- `GET /health` — проверка статуса
- `GET /config` — текущая конфигурация + переменные окружения (секреты маскируются `****`)
- `POST /echo` — эхо-ответ с метаданными сервера

**Переменные окружения:**

| Переменная | По умолчанию | Описание |
|-----------|-------------|----------|
| `PORT` | `8080` | Порт сервера |
| `APP_ENV` | `development` | Окружение |
| `APP_NAME` | `go-config-service` | Имя приложения |
| `MAX_BODY_SIZE` | `10` | Максимальный размер тела (МБ) |
| `READ_TIMEOUT` | `5s` | Таймаут чтения |
| `WRITE_TIMEOUT` | `10s` | Таймаут записи |
| `DB_PASSWORD` | — | Пароль БД (маскируется) |
| `API_KEY` | — | API-ключ (маскируется) |
| `SECRET_KEY` | — | Секретный ключ (маскируется) |

**Запуск без Docker:**

```bash
cd task3_10
go mod tidy
go run main.go
```

**Запуск с переменными:**

```bash
# Windows PowerShell
$env:APP_ENV="production"; $env:DB_PASSWORD="secret"; $env:API_KEY="key-123"; go run main.go

# Windows CMD
set APP_ENV=production && set DB_PASSWORD=secret && set API_KEY=key-123 && go run main.go
```

**Запуск тестов:**

```bash
cd task3_10
go test -v ./...
```

**Сборка Docker-образа:**

```bash
cd task3_10
docker build -t lab11-go-config .
```

**Запуск контейнера:**

```bash
docker run -d --name go-config -p 8080:8080 \
  -e APP_ENV=production \
  -e DB_PASSWORD=secret123 \
  -e API_KEY=key-abc \
  lab11-go-config
```

**Запуск через Docker Compose:**

```bash
cd task3_10
cp .env.example .env
# отредактируйте .env
docker compose up -d
```

**Проверка эндпоинтов:**

```bash
curl http://localhost:8080/health
curl http://localhost:8080/config
curl -X POST http://localhost:8080/echo -H "Content-Type: application/json" -d "{\"message\":\"hello\"}"
```

**Проверка статуса:**

```bash
docker ps
```

**Просмотр логов:**

```bash
docker logs go-config
```

**Остановка:**

```bash
docker compose down
```

**Остановка с удалением образов:**

```bash
docker compose down --rmi all
```

**Просмотр размера образа:**

```bash
docker images lab11-go-config
```

### Повышенная сложность:

1. **Задание 2:** Собрать Rust-приложение с поддержкой musl для полностью статической сборки.

**Эндпоинты:**

- `GET /health` — проверка статуса
- `GET /info` — информация о сервисе
- `GET /hello` — приветственное сообщение

**Запуск без Docker:**

```bash
cd task4_2
cargo run
```

**Запуск тестов:**

```bash
cd task4_2
cargo test
```

**Сборка Docker-образа (static musl):**

```bash
cd task4_2
docker build -t rust-web-server .
```

**Запуск контейнера:**

```bash
docker run -d --name rust-server -p 8080:8080 rust-web-server
```

**Проверка эндпоинтов:**

```bash
curl http://localhost:8080/health
curl http://localhost:8080/info
curl http://localhost:8080/hello
```

**Проверка статуса:**

```bash
docker ps
```

**Просмотр логов:**

```bash
docker logs rust-server
```

**Остановка контейнера:**

```bash
docker stop rust-server
docker rm rust-server
```

**Удаление образа:**

```bash
docker rmi rust-web-server
```

**Просмотр размера образа:**

```bash
docker images rust-web-server
```

**Проверка, что бинарник статический (без динамических зависимостей):**

1. Извлечь бинарник из образа:
```bash
docker create --name temp rust-web-server
docker cp temp:/server task4_2/server
docker rm temp
```

2. Запустить Alpine-контейнер для проверки:
```bash
docker run -d --name check-runner rust:alpine tail -f /dev/null
docker cp task4_2/server check-runner:/server
docker exec check-runner apk add --no-cache file musl-dev
docker exec check-runner file /server
docker exec check-runner sh -c "ldd /server 2>&1 || echo 'STATIC: no dynamic dependencies'"
```

3. Ожидаемый результат `file`:
```
/server: ELF 64-bit LSB pie executable, x86-64, version 1 (SYSV), static-pie linked, stripped
```

Ключевая фраза — **`static-pie linked`** — подтверждает, что бинарник полностью статический.

Команда `ldd` либо выдаст `not a dynamic executable`, либо покажет только `/lib/ld-musl-x86_64.so.1` (встроенный musl loader).

4. Очистка:
```bash
docker rm -f check-runner
rm task4_2/server
```

2. **Задание 8:** Настроить автоматическое обновление контейнеров (watchtower).

**Эндпоинты:**

- `GET /version` — возвращает текущую версию приложения

**Запуск без Docker:**

```bash
cd task5_8
go run main.go
```

**Запуск тестов:**

```bash
cd task5_8
go test -v ./...
```

**Сборка Docker-образа:**

```bash
cd task5_8
docker build -t lab11-go-watchtower .
```

**Запуск через Docker Compose (приложение + watchtower):**

```bash
cd task5_8
docker compose up -d
```

**Проверка эндпоинта:**

```bash
curl http://localhost:8080/version
```

**Проверка статуса контейнеров:**

```bash
docker compose -f task5_8/docker-compose.yml ps
```

**Просмотр логов watchtower:**

```bash
docker compose -f task5_8/docker-compose.yml logs watchtower
```

**Остановка:**

```bash
cd task5_8
docker compose down
```

**Остановка с удалением образов:**

```bash
cd task5_8
docker compose down --rmi all
```

**Как работает Watchtower:**

Watchtower проверяет наличие новых образов каждые **30 секунд**. Если в реестре появляется новая версия образа с меткой `com.centurylinklabs.watchtower.enable=true`, контейнер автоматически обновляется и перезапускается.

**Проверка обновления версии:**

1. Убедитесь, что compose запущен:
```bash
cd task5_8
docker compose up -d
```

2. Измените версию в `main.go`:
```go
c.JSON(200, gin.H{
    "version": "2.0.0",
})
```

3. Пересоберите и обновите образ:
```bash
cd task5_8
docker compose build app
docker compose up -d app
```

4. Проверьте, что версия обновилась:
```bash
curl http://localhost:8080/version
# Ожидаемый ответ: {"version":"2.0.0"}
```

5. Watchtower автоматически обнаружит новый образ и обновит контейнер при следующей проверке (каждые 30 секунд).