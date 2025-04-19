package main

import (
	"log"
	"net/http"

	"protein-shop-api/api"
	"protein-shop-api/db"
	"protein-shop-api/middleware"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	db.InitDB()
	defer db.CloseDB()

	router := mux.NewRouter()

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	}).Methods("GET")

	router.HandleFunc("/auth/register", api.RegisterHandler).Methods("POST")
	router.HandleFunc("/auth/login", api.LoginHandler).Methods("POST")
	router.HandleFunc("/products", api.ProductsHandler).Methods("GET")

	router.HandleFunc("/orders", middleware.AuthMiddleware(api.CreateOrderHandler)).Methods("POST")
	router.HandleFunc("/orders", middleware.AuthMiddleware(api.GetOrdersHandler)).Methods("GET")
	router.HandleFunc("/loyalty/points", middleware.AuthMiddleware(api.LoyaltyPointsHandler)).Methods("GET")
	router.HandleFunc("/profile", middleware.AuthMiddleware(api.ProfileHandler)).Methods("GET")

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
