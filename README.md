# ДЗ Level 0

## Задание

Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе.

Данное задание предполагает создание небольшого микросервиса на `Go` с использованием базы данных и очереди сообщений. Сервис будет получать данные заказов из очереди `Kafka`, сохранять их в базу данных `PostgreSQL` и кэшировать в памяти для быстрого доступа.

## Запуск БД и Kafak

```shell
docker compose up
```

## Запуск миграций

```shell
liquibase update --changelog-file=./migration.sql --username=user --password=password --url=jdbc:postgresql://localhost:5432/db
```

## Запуск приложения
```shell
go run cmd/app/main.go 
```

## Приложение

По [http://localhost:8000/index.html](http://localhost:8000/index.html) можно открыть простой UI для просмотра заказов, сделанный на шаблоне

По [http://localhost:8000/static/indexjs.html](http://localhost:8000/static/indexjs.html) можно открыть простой UI для просмотра заказов, сделанные на JS скрипте

По [http://localhost:8000/api/v1/orders?uid=1](http://localhost:8000/api/v1/orders?uid=1) можно получить json с информацией о заказе по UID

## Kafka UI

По [http://localhost:8082/](http://localhost:8082/) доступен KafkaUi, в качестве сервера надо выбрать 
`http://kafka` с портом `9092`

Через этот сервис можно отправлять сообщения в брокер

## Демонстрация

[Видео](https://drive.google.com/file/d/1lPzT9bvUOHim2tqhh-sSiZ4z_pfdRZus/view?usp=sharing)
