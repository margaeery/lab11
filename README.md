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


3. **Задание 10:** Использовать переменные окружения для конфигурации.


### Повышенная сложность:

1. **Задание 2:** Собрать Rust-приложение с поддержкой musl для полностью статической сборки. 

2. **Задание 8:** Настроить автоматическое обновление контейнеров (watchtower). 