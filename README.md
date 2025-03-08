# MerchShop_Service
Сервис для покупки мерча.
## Основные задействованные сторонние библиотеки:
1. gorilla/mux для роутинга
2. Masterminds/squirrel для построения SQL-запросов
3. golang-jwt/jwt/v5 для работы с JWT-токенами
4. golang.org/x/crypto/bcrypt для шифрования паролей в БД
5. vektra/mockery/v2 для генерации моков

## Быстрый старт
```bash
make
```
```bash
make migrate-up
```

## Остановить приложение:
```bash
make stop
```
или
```bash
make remove
```

## Описание Makefile
1. Сборка и запуск приложения
```bash
make all
```

2. Собирает Docker-образ приложения
```bash
make build
```

3. Запускает Docker-контейнер с приложением
```bash
make run
```

4. Останавливает Docker-контейнер с приложением (без удаления контейнеров)
```bash
make stop
```

5. Останавливает Docker-контейнер с приложением (удаляя контейнеры)
```bash
make remove
```

6. Накатывает миграции базы данных
```bash
make migrate-up
```

7. Откатывает миграции базы данных
```bash
make migrate-down
```

> [!NOTE]
> Для работы с миграциями на Win предусмотрен флаг os=win      
> Пример:
> ```bash
> make os=win migrate-up
> ```
