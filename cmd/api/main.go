package main

import (
	"carthero/internal/server"
	"fmt"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

func main() {

	server := server.New()

	server.RegisterFiberRoutes()
	server.RegisterRiderRoutes()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	err := server.Listen(fmt.Sprintf(":%d", port))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
