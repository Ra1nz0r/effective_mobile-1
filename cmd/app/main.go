package main

import (
	_ "github.com/Ra1nz0r/effective_mobile-1/docs"
	"github.com/Ra1nz0r/effective_mobile-1/internal/server"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// @title Music Library API
// @version 1.0
// @description REST API для управления онлайн-библиотекой песен. Включает функции добавления, обновления, удаления и поиска песен, а также взаимодействие с внешними сервисами для получения дополнительной информации о композициях.
// @termsOfService http://swagger.io/terms/

// @contact.name Artem Rylskii
// @contact.url https://t.me/Rainz0r
// @contact.email n52rus@gmail.com

// @host localhost:7654
// @BasePath /
func main() {
	server.Run()
}
