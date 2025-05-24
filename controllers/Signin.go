package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/fenilpanseriya/docs2.0/db"
	"github.com/fenilpanseriya/docs2.0/helpers"
	"github.com/fenilpanseriya/docs2.0/models"
	"go.mongodb.org/mongo-driver/bson"
)

func Signin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	jwtKey := os.Getenv("JWT_KEY")
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	filter := bson.M{"email": user.Email, "password": user.Password}
	var userCollection = db.Client.Database("docs").Collection("users")
	count, err := userCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	tokenString, err := helpers.GenerateToken(&user, jwtKey, time.Hour*24*7)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Access_token = tokenString
	user.Updated_at = time.Now()
	user.Refresh_token = tokenString
	//update the user
	_, err = userCollection.UpdateOne(context.Background(), filter, bson.M{
		"$set": bson.M{
			"access_token":  user.Access_token,
			"updated_at":    user.Updated_at,
			"refresh_token": user.Refresh_token,
		},
	},
	)
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}
	//fmt.Println("users", users)
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * 7 * time.Hour),
		HttpOnly: true,
	})
	response := map[string]string{
		"message": "Signin successful",
		"email":   user.Email,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
