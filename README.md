Здравствуйте. Меня зовут Харсиев Мухаммед и это мое тестовое задание для компании Avito на позицию golang backend разработчик. Проект был упакован в docker контейнер для удобства работы. Также были подключены миграции для удобства тестирования. Для запуска и работы с проектом вам понадобятся:
-установленный docker на вашем компьютере.
-postman для работы с запросами. 

1. Сборка и запуск в Docker
Требования:
-Docker
-Docker Compose

Шаги:
Клонируйте репозиторий:
   bash
   git clone https://github.com/chikalet/pvz-service.git
   cd pvz-service
   
Соберите образы:
   bash
   docker compose build

Запустите сервисы:
   bash
   docker compose up -d

   
Файлы:
docker-compose.yml - описывает сервисы (приложение, БД, Prometheus, Grafana)
Dockerfile - двухэтапная сборка Go-приложения

2. Используемые технологии
Технология	Назначение
-Go 1.23	Бэкенд на языке Go
-Fiber v3	Высокопроизводительный веб-фреймворк
-PostgreSQL 15	Основное хранилище данных
-Prometheus	Сбор метрик (порт 9000)
-Grafana	Визуализация метрик (порт 3000)
-JWT	Аутентификация пользователей
-pgx	Драйвер PostgreSQL для Go


4. Пользователи и роли
Роли:
Модератор:
-Создание/удаление ПВЗ
-Просмотр всей статистики

Сотрудник ПВЗ:
-Управление приёмками товаров
-Добавление/удаление товаров

Клиент (реализация через /register):
-Просмотр статуса заказов

Создание пользователей:
   bash
   curl "http://localhost:8080/api/v1/dummyLogin?role=moderator"

Регистрация нового пользователя:
   curl -X POST "http://localhost:8080/api/v1/register" \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "password": "secret", "role": "employee"}'

     
4. Работа с данными
Пункты выдачи заказов (ПВЗ):

Создание ПВЗ (только модератор):
   curl -X POST "http://localhost:8080/api/v1/pvz" \
     -H "Authorization: Bearer <token>" \
     -d '{"city": "Москва"}'

Получение списка ПВЗ:

   curl "http://localhost:8080/api/v1/pvz?start_date=2025-01-01&end_date=2025-12-31" \
     -H "Authorization: Bearer <token>"

     
Приёмки товаров:


   curl -X POST "http://localhost:8080/api/v1/intake" \
     -H "Authorization: Bearer <token>" \
     -d '{"pvz_id": 1}'

Добавить товар:

   curl -X POST "http://localhost:8080/api/v1/pvz/1/items" \
     -H "Authorization: Bearer <token>" \
     -d '{"type": "electronics", "quantity": 5}'

Закрыть приемку:

   curl -X POST "http://localhost:8080/api/v1/pvz/1/close" \
     -H "Authorization: Bearer <token>"
     
Ограничения:

ПВЗ можно создавать только в Москве, СПб и Казани

Товары можно добавлять только в открытую приёмку

Удаление товаров только по принципу LIFO

5. Мониторинг
Prometheus (http://localhost:9000):

Метрики приложения:

text
http_requests_total{method="POST", path="/pvz", status="201"}
pvz_created_total
intakes_opened_total
Grafana (http://localhost:3000):


6. Производительность
RPS: ≥ 1000 запросов/сек
Латентность: < 100 мс (p99)
Доступность: 99.99% SLI

7. Безопасность
Все запросы (кроме /dummyLogin) требуют JWT-токена
Пароли хранятся в виде bcrypt-хешей
Автоматическая очистка старых сессий
