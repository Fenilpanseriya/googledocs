package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fenilpanseriya/docs2.0/controllers"
	"github.com/fenilpanseriya/docs2.0/db"
	"github.com/fenilpanseriya/docs2.0/middleware"
	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // load .env file
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	mongoURI := os.Getenv("MONGO_URI")
	fmt.Println("Mongo URI: ", mongoURI)
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	db.ConnectMongoDB(mongoURI)

	// Create a new ServeMux for routing
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/signup", controllers.Signup)
	mux.HandleFunc("/api/v1/signin", controllers.Signin)
	mux.HandleFunc("/api/v1/forgotpassword", controllers.ForgotPassword)
	mux.HandleFunc("/api/v1/resetpassword", middleware.AuthMiddleware(controllers.ResetPassword))
	mux.HandleFunc("/api/v1/me", middleware.AuthMiddleware(controllers.UserDetails))
	mux.HandleFunc("/", controllers.Welcome)

	// Set CORS headers
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:5173"})
	credentials := handlers.AllowCredentials()
	// Start the server with CORS support
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins, credentials)(mux)))
}
