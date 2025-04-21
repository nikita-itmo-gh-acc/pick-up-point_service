# pick-up-point_service
Avito tech internship backend assignment task

# Как запустить в докере?
## Требования:
- Установленный docker
- Docker Desktop (для Windows)
## 1 Собираем контейнеры: 
```shell
docker compose up -d --build
```
## 2 Чтобы посмотреть логи приложения:
```shell
docker compose logs app
```
## 3 Чтобы заполнить таблицы БД минимальными данными (опционально):
```shell
docker exec -it db_container pg_restore -U postgres -d songs_db /backups/songsData.sql
```
