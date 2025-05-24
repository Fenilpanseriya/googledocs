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

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	userCollection := db.Client.Database("docs").Collection("users")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	filter := bson.M{"email": user.Email}
	count, err := userCollection.CountDocuments(context.Background(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	// Generate a password reset token
	tokenString, err := helpers.GenerateToken(&user, os.Getenv("JWT_KEY"), time.Minute*10)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	helpers.SendEmail(user.Email, "Password Reset link", "http://localhost:8080/resetpassword?token="+tokenString)
	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Password reset link sent to your email",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
