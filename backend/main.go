package main

import (
	"backend/routes"
	"backend/config"
	"log"
	"net/http"
)

func main(){
	config.LoadEnv()
	r := routes.StockRoutes()
	port := config.LoadPort()

	log.Println("Server running in: http://localhost:"+port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start the server: %v", err)
	}
}