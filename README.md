# Effective Mobile Test API 

## Требования

Проект разрабоатывался на **Golang** версии 1.21.5.

## Установка

1. Клонируйте репозиторий с проектом:
   ```bash
   git clone https://github.com/Dermofet/medods
   ```
   
2. Перейдите в директорию проекта:
    ```bash
    cd medods
    ```
    
3. Установите зависимости:
    ```bash
    go mod download
    ```
    
## Запуск

Для запуска Docker контейнеров выполните следующие команды:

```bash
docker compose -f ./dev/docker-compose.yml up -d --build
```

## Конфигурация

Для конфигурации проекта используется файл **.env**.

#### Переменные окружения

Для настройки приложения:
- APP_NAME
- APP_VERSION
- LOG_LEVEL
- DEBUG
- PATH_LOG

Для настройки аутентификации:
- API_KEY

Для настройки http сервера:
- HTTP_HOST
- HTTP_PORT

Для настройки базы данных:
- DB_HOST
- DB_PORT
- DB_NAME
- DB_USER
- DB_PASS

## API

### Эндпоинт 1: Получение новых токенов

**Путь**: /auth/new-tokens

**Метод**: POST

**Описание**: Этот эндпоинт предназначен для создания новых токенов по guid пользователя.

**Параметры запроса:**
- guid: guid пользователя

**Пример запроса:**
```text
POST /auth/new-tokens?guid=<guid>
```

**Примеры ответов**
- Статус 200 OK
    ```json
    {
      "access_token": "<access_token>",
      "refresh_token": "<refresh_token>"
    }
    ```
- Статус 409 Conflict
- Статус 422 UnprocessableEntity
- Статус 500 InternalServerError

### Эндпоинт 2: Обновление токенов

**Путь**: /auth/refresh-tokens

**Метод**: POST

**Описание**: Этот эндпоинт предназначен для обновления токенов по refresh токену.

**Тело запроса:**
```json
{
  "refresh_token": "<refresh_token>"
}
```

**Пример запроса:**
```text
POST /auth/refresh-tokens
Content-Type: application/json

{
  "refresh_token": "<refresh_token>"
}
```

**Примеры ответов**
- Статус 200 OK
    ```json
    {
      "access_token": "<access_token>",
      "refresh_token": "<refresh_token>"
    }
    ```
- Статус 409 Conflict
- Статус 422 UnprocessableEntity
- Статус 500 InternalServerError
        
