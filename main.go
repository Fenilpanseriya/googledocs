package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fenilpanseriya/docs2.0/controllers"
	"github.com/fenilpanseriya/docs2.0/db"
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

	http.HandleFunc("/", controllers.Welcome)
	http.HandleFunc("/signup", controllers.Signup)
	http.HandleFunc("/signin", controllers.Signin)
	// http.HandleFunc("/signin", Signin)
	// http.HandleFunc("/welcome", Welcome)
	//http.HandleFunc("/refresh", Refresh)
	//http.HandleFunc("/logout", Logout)

	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
