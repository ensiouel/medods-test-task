# REST API Сервис аутентификации

---

## Стек технологий

#### API
- Fiber
- JWT
#### Хранилище:
- MongoDB

---

## Deployment

**Build** application

```shell
docker compose build
```

---

**Run** application

```shell
docker compose up -d
```

---

## Примеры использования

### Получить пару access и refresh токенов
#### Запрос
```http request
GET http://localhost:8080/auth?guid=7e484497-a798-4655-b5d9-56b16e990e12
```

#### Ответ
```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI4NzA3MTEsImd1aWQiOiI3ZTQ4NDQ5Ny1hNzk4LTQ2NTUtYjVkOS01NmIxNmU5OTBlMTIiLCJzZXNzaW9uX2lkIjoiZmNiYTIwY2ItMjJiZC00NTY5LWFlZjgtYjJlZmU3NzEzZDY0In0.kYMrY-CmwNZTI8K6BFxu1aw4i4CFyYAeSEiO4is7cuqvf_lzKPaSWvnXtWJ3X6Cc_alZbClmioWn3wosJs4L2Q",
  "refresh_token": "Juk4ZGvfQY6ebFBUPPO0jg=="
}
```

---

### Обновить пару токенов с помощью refresh токена и id сессии
#### Запрос
```http request
POST http://localhost:8080/fcba20cb-22bd-4569-aef8-b2efe7713d64/refresh
Content-Type: application/json

{
  "refresh_token": "UMxZHuVoTnaRr5dFzPm77A=="
}
```

#### Ответ
```json
{
  "access_token": "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTI4NzEwMzAsImd1aWQiOiI3ZTQ4NDQ5Ny1hNzk4LTQ2NTUtYjVkOS01NmIxNmU5OTBlMTIiLCJzZXNzaW9uX2lkIjoiZmNiYTIwY2ItMjJiZC00NTY5LWFlZjgtYjJlZmU3NzEzZDY0In0.gKRSeT7--TcDoPT04yjJl8MawPwjXZ_9w9ZqfuW9fmeswcEm24-Hg05yUn6fVOH57upPjOO9CGEjihprGhBUXw",
  "refresh_token": "UMxZHuVoTnaRr5dFzPm77A=="
}
```

---

## Конфигурации

### Все параметры загружаются из файта **[.env](.env)**

```dotenv
SERVER_ADDR=:8080
MONGO_HOST=mongo
MONGO_PORT=27017
SESSION_ACCESS_TOKEN_EXPIRATION=1h
SESSION_REFRESH_TOKEN_EXPIRATION=720h
SESSION_SECRET=secret
```