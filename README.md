# good_api_test

## Описание

Этот проект представляет собой Go-приложение, реализующее RESTful API для управления товарами (goods). Он использует PostgreSQL для хранения данных, Redis для кэширования и NATS для обмена сообщениями.

## Используемые технологии

*   **Go**
*   **Gin Web Framework**
*   **PostgreSQL**
*   **Redis**
*   **NATS**
*   **Docker & Docker Compose**

## Запуск проекта

Проект можно запустить с использованием Docker Compose, который настроит все необходимые сервисы (PostgreSQL, Redis, NATS) и Go-приложение.

### Требования

*   [Docker](https://docs.docker.com/get-docker/)
*   [Docker Compose](https://docs.docker.com/compose/install/)

### Шаги по запуску

1.  **Клонируйте репозиторий:**
    ```bash
    git clone https://github.com/Alexandrjob/good_api_test.git
    ```

2.  **Запустите сервисы с помощью Docker Compose:**
    ```bash
    docker-compose up --build
    ```
    
3.  **Проверка работы:**
    После запуска всех контейнеров, API будет доступен по адресу `http://localhost:8080/api/v1`.
    
### Тестирование
Вы можете импортировать следующую коллекцию Postman для удобного тестирования API:

1. **Postman Collection:**
```http request
https://web.postman.co/workspace/My-Workspace~a8f3e728-46ba-4698-8343-9d6e0f18fa6c/collection/37898445-49c254e0-f520-4a59-b2cf-fdc602e03097?action=share&source=copy-link&creator=37898445
```
Или
```
https://.postman.co/workspace/My-Workspace~a8f3e728-46ba-4698-8343-9d6e0f18fa6c/collection/37898445-49c254e0-f520-4a59-b2cf-fdc602e03097?action=share&creator=37898445
```

## Переменные окружения

Проект использует следующие переменные окружения. Вы можете задать их в файле `.env` в корне проекта (рядом с `go.mod`) или через системные переменные окружения.
*   `POSTGRES_USER`: Пользователь PostgreSQL.
*   `POSTGRES_PASSWORD`: Пароль PostgreSQL.
*   `POSTGRES_DB`: Имя базы данных PostgreSQL.
*   `POSTGRES_HOST`: Хост PostgreSQL.
*   `REDIS_ADDR`: Адрес Redis.
*   `NATS_URL`: URL NATS сервера.

## API Эндпоинты

Все запросы предполагают базовый URL: `http://localhost:8080/api/v1`

### 1. Получить информацию о товаре (GET /good)

*   **Описание:** Получает информацию о товаре по его ID и ProjectID.
*   **Метод:** `GET`
*   **URL:** `/good?id={id}&projectId={projectId}`
*   **Пример URL:** `http://localhost:8080/api/v1/good?id=1&projectId=123`
*   **Тело запроса (JSON):** Отсутствует

### 2. Создать новый товар (POST /good/create)

*   **Описание:** Создает новый товар.
*   **Метод:** `POST`
*   **URL:** `/good/create?projectId={projectId}`
*   **Пример URL:** `http://localhost:8080/api/v1/good/create?projectId=123`
*   **Тело запроса (JSON):**
    ```json
    {
        "name": "Название товара",
        "description": "Описание товара"
    }
    ```

### 3. Обновить существующий товар (PATCH /good/update)

*   **Описание:** Обновляет информацию о существующем товаре по его ID и ProjectID.
*   **Метод:** `PATCH`
*   **URL:** `/good/update?id={id}&projectId={projectId}`
*   **Пример URL:** `http://localhost:8080/api/v1/good/update?id=1&projectId=123`
*   **Тело запроса (JSON):**
    ```json
    {
        "name": "Обновленное название",
        "description": "Обновленное описание"
    }
    ```

### 4. Удалить товар (DELETE /good/remove)

*   **Описание:** Удаляет товар по его ID и ProjectID.
*   **Метод:** `DELETE`
*   **URL:** `/good/remove?id={id}&projectId={projectId}`
*   **Пример URL:** `http://localhost:8080/api/v1/good/remove?id=1&projectId=123`
*   **Тело запроса (JSON):** Отсутствует

### 5. Изменить приоритет товара (PATCH /good/reprioritize)

*   **Описание:** Изменяет приоритет товара по его ID и ProjectID.
*   **Метод:** `PATCH`
*   **URL:** `/good/reprioritize?id={id}&projectId={projectId}`
*   **Пример URL:** `http://localhost:8080/api/v1/good/reprioritize?id=1&projectId=123`
*   **Тело запроса (JSON):**
    ```json
    {
        "newPriority": 5
    }
    ```

### 6. Получить список всех товаров (GET /goods/list)

*   **Описание:** Получает список всех товаров с возможностью пагинации.
*   **Метод:** `GET`
*   **URL:** `/goods/list?limit={limit}&offset={offset}`
*   **Пример URL:** `http://localhost:8080/api/v1/goods/list?limit=10&offset=0`
*   **Тело запроса (JSON):** Отсутствует
