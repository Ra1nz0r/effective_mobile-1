package main

import (
	_ "github.com/Ra1nz0r/effective_mobile-1/docs"
	"github.com/Ra1nz0r/effective_mobile-1/internal/server"
)

// @title Music Library
// @version 1.0
// @description Implementation of an online song library.
// @termsOfService http://swagger.io/terms/

// @contact.name Artem Rylskii
// @contact.url https://t.me/Rainz0r
// @contact.email n52rus@gmail.com

// @host localhost:7654
// @BasePath /
func main() {
	server.Run()
}
