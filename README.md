<div align="center"> <h1 align="center"> ТЗ: Реализация онлайн библиотеки песен 🎶 </h1> </div>

__Go веб-сервер и агент. Сервер собирает метрики и данные, которые отправляет агент, обновление и отправка происходит с указанным интервалом времени.__

[Инструкция по локальному запуску и информация по приложению.](#local)\
[Инструкция по созданию Docker образа и запуску контейнера.](#docker)\
[Инструкция по запуску PostgreSQL.](#postgresql)

***
#### Инструкция по локальному запуску и информация по приложению.

По-умолчанию приложение запускается на ```localhost:8080```

- Программу можно запускать двумя способами через терминал.
    - Обычные команды. 
    - Короткими командами из TaskFile.
<div>

- ___Для запуска сервера в терминале.___\
```go run ./cmd/server``` или ```task run_s```
- ___Для запуска сервера в терминале.___\
```go run ./cmd/agent``` или ```task run_a```
- ___Для запуска тестов в терминале.___\
```go test -v ./... -count=1``` или ```task test```

<a name="docker"></a>
***
#### Инструкция по созданию Docker образа и запуску контейнера.

1) TODO

<a name="postgresql"></a>
***
 #### Инструкция по запуску PostgreSQL.

1) TODO

***