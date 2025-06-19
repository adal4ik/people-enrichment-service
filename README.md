# people-enrichment-service

Сервис для обогащения информации о людях по ФИО через открытые API (возраст, пол, национальность) с сохранением в PostgreSQL.

## Возможности

- Добавление нового человека (обогащение через публичные API)
- Получение списка людей с фильтрами и пагинацией
- Получение информации о человеке по id
- Обновление информации о человеке
- Удаление человека по id
- Логирование (zap)
- Swagger-документация (docs/swagger.yml)
- Конфигурация через .env

## Используемые публичные API

- [Agify.io (возраст)](https://api.agify.io)
- [Genderize.io (пол)](https://api.genderize.io)
- [Nationalize.io (национальность)](https://api.nationalize.io)

## Быстрый старт

### 1. Клонируйте репозиторий

```sh
git clone https://github.com/your-username/people-enrichment-service.git
cd people-enrichment-service
```

### 2. Заполните `.env`

Пример содержимого:
```
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=people
APP_PORT=8080
```

### 3. Запустите сервис и БД через Docker

```sh
make up
```

Сервис будет доступен на [http://localhost:8080](http://localhost:8080)

### 4. Остановить сервис

```sh
make down
```

### 5. Локальный запуск (без Docker)

- Убедитесь, что PostgreSQL запущен и параметры совпадают с .env
- Выполните миграции (например, через psql или мигратор)
- Запустите приложение:

```sh
make run
```

### 6. Справка по make-командам

Для просмотра всех доступных команд Makefile выполните:

```sh
make help
```

Вы увидите список команд, которые можно использовать для управления сервисом, контейнерами и сборкой проекта.

## Документация API

Swagger-описание доступно в файле [`docs/swagger.yml`](docs/swagger.yml).  
Можно визуализировать через [Swagger Editor](https://editor.swagger.io/).

## Как посмотреть Swagger

1. Запустите сервис командой:
   ```sh
   make run
   ```
2. Перейдите в браузере на [http://localhost:8081](http://localhost:8081)
3. Введите в поле поиска `/swagger.yml` и нажмите Enter — откроется документация вашего API.


## Примеры запросов

**Добавить человека:**
```http
POST /person
Content-Type: application/json

{
  "name": "Dmitriy",
  "surname": "Ushakov",
  "patronymic": "Vasilevich"
}
```

**Получить список:**
```http
GET /person?limit=10&offset=0&name=Dmitriy
```

## Тесты

Для запуска unit-тестов:
```sh
go test ./...
```

## Структура проекта

```
.
├── cmd/app/                # Точка входа приложения
├── internal/
│   ├── handler/            # HTTP-ручки (контроллеры)
│   ├── service/            # Бизнес-логика
│   ├── repository/         # Работа с БД
│   ├── models/             # Описания структур
│   └── config/             # Загрузка конфигурации
├── migrations/             # SQL-миграции для БД
├── docs/                   # Swagger и другая документация
├── utils/                  # Вспомогательные функции
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

---

**Автор:**  
[Ваше имя или ник]
