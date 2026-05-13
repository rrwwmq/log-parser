# Log Parser

Микросервис на Go для парсинга лог-файлов InfiniBand сети. Принимает zip-архивы с логами, парсит их, сохраняет данные в PostgreSQL и предоставляет REST API для получения информации о топологии сети.

## Стек

- Go 1.22+
- PostgreSQL
- Docker / Docker Compose

## Структура проекта

```
log-parser/
├── cmd/logparser/       # точка входа
├── internal/
│   ├── core/            # общая инфраструктура
│   ├── parser/          # парсинг лог-файлов
│   └── features/        # бизнес-логика (logs, nodes, ports)
├── migrations/          # SQL миграции
└── data/                # папка для zip-архивов
```

## Как собрать и запустить

### 1. Клонировать репозиторий

```bash
git clone https://github.com/rrwwmq/log-parser
cd log-parser
```

### 2. Создать .env файл

```bash
cp .env.example .env
```

### 3. Запустить базу данных

```bash
make env-up
```

### 4. Применить миграции

```bash
make migrate-up
```

### 5. Запустить приложение

```bash
make logparser-run
```

Сервис запустится на `http://localhost:8080`.

### Запуск приложения через Docker

```bash
make app-up
```

Остановить:
```bash
make app-down
```

## Использование

Положите zip-архив с логами в папку `data/` и отправьте запрос на парсинг.

## API

### POST /api/v1/parse
Загрузить и распарсить лог-файл. Файл должен лежать в папке `data/`.

```bash
curl -X POST http://localhost:8080/api/v1/parse \
  -H "Content-Type: application/json" \
  -d '{"file_path": "data/log.zip"}'
```

Ответ:
```json
{
  "log_id": 1
}
```

### GET /api/v1/log/{log_id}
Мета-информация о логе.

```bash
curl http://localhost:8080/api/v1/log/1
```

Ответ:
```json
{
  "id": 1,
  "file_name": "data/log.zip",
  "status": "done",
  "uploaded_at": "2026-05-13T09:26:51Z",
  "node_count": 5,
  "port_count": 260
}
```

### GET /api/v1/topology/{log_id}
Список узлов топологии.

```bash
curl http://localhost:8080/api/v1/topology/1
```

Ответ:
```json
{
  "nodes": [
    {
      "id": 1,
      "node_guid": "0xswitch1",
      "node_desc": "SWITCH_1",
      "node_type": "switch",
      "num_ports": 65
    },
    {
      "id": 5,
      "node_guid": "0xhost1",
      "node_desc": "HOST_1",
      "node_type": "host",
      "num_ports": 1
    }
  ]
}
```

### GET /api/v1/node/{node_id}
Детали узла.

```bash
curl http://localhost:8080/api/v1/node/1
```

Ответ:
```json
{
  "id": 1,
  "node_guid": "0xswitch1",
  "node_desc": "SWITCH_1",
  "node_type": "switch",
  "num_ports": 65,
  "info": {
    "serial_number": "SOS123",
    "part_number": "MMM-MAV",
    "revision": "AA",
    "product_name": "Gorilla",
    "endianness": 10,
    "enable_endianness_per_job": 0,
    "reproducibility_disable": 0
  }
}
```

### GET /api/v1/port/{node_id}
Порты узла.

```bash
curl http://localhost:8080/api/v1/port/1
```

Ответ:
```json
{
  "ports": [
    {
      "id": 1,
      "port_guid": "0xswitch1",
      "port_num": 0,
      "port_state": 0,
      "lid": 22
    }
  ]
}
```

## Граф топологии

На основе nodes и ports можно построить граф сети:

- **Узлы** — hosts и switches
- **Связи** — switch соединяется с host или другим switch через порты. Связь определяется по совпадению `LID` у портов двух узлов — если порт switch и порт host имеют одинаковый ненулевой `LID`, они соединены.

## Переменные окружения

| Переменная | Описание | По умолчанию |
|---|---|---|
| `HTTP_ADDR` | Адрес сервера | `:8080` |
| `HTTP_SHUTDOWN_TIMEOUT` | Таймаут остановки | `30s` |
| `POSTGRES_HOST` | Хост БД | — |
| `POSTGRES_PORT` | Порт БД | `5432` |
| `POSTGRES_USER` | Пользователь БД | — |
| `POSTGRES_PASSWORD` | Пароль БД | — |
| `POSTGRES_DB` | Имя БД | — |
| `POSTGRES_TIMEOUT` | Таймаут запросов | — |
| `LOGGER_LEVEL` | Уровень логирования | `DEBUG` |