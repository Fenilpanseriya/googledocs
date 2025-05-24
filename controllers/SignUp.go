package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/fenilpanseriya/docs2.0/db"
	"github.com/fenilpanseriya/docs2.0/helpers"
	"github.com/fenilpanseriya/docs2.0/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	jwtKey := os.Getenv("JWT_SECRET")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	filter := bson.M{"email": user.Email}
	var userCollection = db.Client.Database("docs").Collection("users")
	count, er := userCollection.CountDocuments(context.Background(), filter)
	if er != nil {
		http.Error(w, er.Error(), http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	tokenString, err := helpers.GenerateToken(&user, jwtKey, time.Hour*24*7)
	fmt.Println("jwtKey", jwtKey, user.Email)
	fmt.Println("tokenString", tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	user.Id = primitive.NewObjectID()
	user.Access_token = tokenString
	user.Created_at = time.Now()
	user.Updated_at = time.Now()
	user.Refresh_token = tokenString

	//create a new user
	_, err = userCollection.InsertOne(context.Background(), user)
	if err != nil {
		fmt.Println("Error inserting user:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
	})
	response := map[string]string{
		"message": "Signup successful",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
