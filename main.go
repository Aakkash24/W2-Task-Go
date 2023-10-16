package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"payment/configs"
	"payment/routes"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func EnvMongoURI() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv("MONGOURL")
}

const port = 4000

func main() {
	router := mux.NewRouter()
	routes.UserRoute(router)

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	_ = configs.ConnectDB()
	fmt.Println("Server started at port:", port)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("Error", err)
	}
}
