# URL Shortener Service (L3.2wb)

Высокопроизводительный сервис сокращения ссылок с аналитикой переходов, написанный на Go. Проект демонстрирует применение принципов **Clean Architecture**, **Dependency Injection** и **Graceful Shutdown**.

---

## Особенности

- **Сокращение ссылок**: Генерация уникальных коротких кодов (base64 URL encoding).
- **Мгновенный редирект**: Обработка перенаправлений с минимальной задержкой.
- **Аналитика**:
    - Подсчет общего количества кликов.
    - Учет уникальных посетителей (по IP).
    - Сбор User-Agent и времени перехода.
- **Архитектура**: Четкое разделение на слои (Transport, Service, Repository).
- **Инфраструктура**:
    - Полная контейнеризация (Docker & Docker Compose).
    - Собственная CLI утилита для миграций БД.
    - Структурированное логирование (Zerolog).
    - Graceful Shutdown для безопасного завершения работы.

---

## Технический стек

| Категория | Технологии |
|-----------|------------|
| **Language** | Go (Golang) |
| **Framework** | Gin Web Framework |
| **Database** | PostgreSQL 15 |
| **Migration** | Golang-Migrate (Custom CLI wrapper) |
| **Logging** | Zerolog |
| **Config** | Viper (YAML + Env) |
| **Deploy** | Docker, Docker Compose |
| **Testing** | Testify, Mockery |

---

## Архитектура проекта

Проект следует структуре **Standard Go Project Layout**:

```
├── cmd/
│ ├── service/ # Точка входа основного приложения (DI container)
│ └── migrator/ # CLI утилита для управления миграциями БД
├── internal/
│ ├── app/ # Логика жизненного цикла приложения (Run, Shutdown)
│ ├── config/ # Загрузка и валидация конфигурации
│ ├── domain/ # Доменные модели и кастомные ошибки
│ ├── repository/ # Слой работы с БД (PostgreSQL)
│ ├── service/ # Бизнес-логика
│ └── transport/ # HTTP хендлеры, роутинг, middleware
├── migrations/ # SQL файлы миграций
└── docker-compose.yaml
```


---

## Middlewares
- **Request Logger**: Логирование каждого запроса (метод, путь, статус, время выполнения, IP).
- **Recovery**: Перехват паники (panic) для предотвращения падения сервера.
- **Error Handler**: Централизованная обработка ошибок и приведение их к единому JSON-формату ответа.

