# Wallet Service

Wallet Service — это простое REST API-приложение на Go для управления кошельками. В рамках проекта реализованы операции пополнения и снятия средств, а также получение текущего баланса. Данные хранятся в PostgreSQL, а всё приложение контейнеризовано с помощью Docker и docker-compose.

## Функциональность

- **POST `/api/v1/wallet`**  
Выполняет операцию пополнения или снятия средств, изменяя баланс кошелька в базе данных.
  Принимает JSON-запрос:
  ```json
  {
    "walletId": "UUID",
    "operationType": "DEPOSIT" | "WITHDRAW",
    "amount": 1000
  }


- **GET `/api/v1/wallets/{walletId}`**
    Возвращает текущий баланс кошелька в формате:
  ```json
    {
      "walletId": "UUID",
      "balance": 1000.0
    }

## Технологии

  Язык программирования: Go
  База данных: PostgreSQL
  Контейнеризация: Docker, docker-compose
  Маршрутизация: Gorilla Mux
  Конфигурация: Переменные окружения считываются из файла config.env (с использованием godotenv)

## Структура проекта

wallet-service/
├── config.env               # Файл конфигурации (переменные окружения)
├── Dockerfile               # Сборка контейнера приложения
├── docker-compose.yml       # Подъём контейнеров (app и postgres)
├── go.mod                   # Go-модуль
├── main.go                  # Точка входа в приложение
├── handlers/                # HTTP-обработчики (endpoints)
│   └── wallet.go
├── models/                  # Модели данных (опционально)
│   └── wallet.go
├── repository/              # Доступ к базе данных
│   ├── db.go
│   └── wallet.go
├── service/                 # Бизнес-логика
│   └── wallet.go
└── tests/                   # Тесты (юнит- и интеграционные)
    └── wallet_test.go

## Настройка окружения

1. Создайте файл config.env в корне проекта и заполните его следующими переменными:

SERVER_PORT=8080
DB_HOST=postgres         # Для запуска через docker-compose (имя сервиса)
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres         # Пароль должен соответствовать настройкам в docker-compose.yml
DB_NAME=walletdb

2. Запуск приложения
Через Docker Compose

    Убедитесь, что Docker и docker-compose установлены.
    В корне проекта выполните:

    docker-compose up --build

    Это поднимет два контейнера:
        postgres — база данных PostgreSQL.
        app — Go-приложение.

Приложение будет доступно по адресу http://localhost:8080.

3. Тестирование
Локальное тестирование

Для запуска тестов ) выполните:

go test -v ./...

## Использование API

Пример запроса для операции (DEPOSIT)

curl -X POST -H "Content-Type: application/json" \
  -d '{"walletId": "test-wallet", "operationType": "DEPOSIT", "amount": 100}' \
  http://localhost:8080/api/v1/wallet

Пример запроса для получения баланса

curl http://localhost:8080/api/v1/wallets/test-wallet
