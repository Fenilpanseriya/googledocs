package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/fenilpanseriya/docs2.0/db"
	"github.com/fenilpanseriya/docs2.0/models"
	"go.mongodb.org/mongo-driver/bson"
)

func UserDetails(w http.ResponseWriter, r *http.Request) {
	var user models.User
	email := r.Context().Value("email")
	if email == nil {
		http.Error(w, "something went wrong not found try again", http.StatusBadRequest)
		return
	}
	userCollection := db.Client.Database("docs").Collection("users")
	filter := bson.M{"email": email}
	userCollection.FindOne(r.Context(), filter).Decode(&user)
	if user.Email == "" {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	publicUser := map[string]interface{}{
		"email": user.Email,
		"id":    user.Id,
	}
	response := map[string]interface{}{
		"message": "User details fetched successfully",
		"user":    publicUser,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
