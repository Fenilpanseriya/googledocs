package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/fenilpanseriya/docs2.0/db"
	"github.com/fenilpanseriya/docs2.0/models"
	"go.mongodb.org/mongo-driver/bson"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	userCollection := db.Client.Database("docs").Collection("users")
	// userID := r.Context().Value("userID")
	email := r.Context().Value("email")
	// fmt.Fprintf(w, "Hello, %s! Your UserID is %v", email, userID)
	fmt.Println("email", email)

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	filter := bson.M{"email": email}
	count, err := userCollection.CountDocuments(r.Context(), filter)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	update := bson.M{"$set": bson.M{"password": user.Password}}
	_, err = userCollection.UpdateOne(r.Context(), filter, update)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := map[string]string{
		"message": "Password reset successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
