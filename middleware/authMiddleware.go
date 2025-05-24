package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fenilpanseriya/docs2.0/db"
	"github.com/fenilpanseriya/docs2.0/models"
	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomClaims struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

// AuthMiddleware validates JWT token and checks user existence
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtKey := os.Getenv("JWT_KEY")
		// Get the token from the cookie
		cookie, err := r.Cookie("auth_token")
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		tokenStr := cookie.Value

		// Parse the token
		claims := &CustomClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		})
		fmt.Println("token", token)
		fmt.Println("valid", token.Valid)
		fmt.Println("claims", claims)
		if err != nil {
			if strings.Contains(err.Error(), "expired") {
				responseMessage := map[string]string{
					"message": "Token expired. Please login again.",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(responseMessage)
				return
			}
			responseMessage := map[string]string{
				"message": "Invalid token",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responseMessage)
			return
		}

		if !token.Valid {
			responseMessage := map[string]string{
				"message": "Invalid token",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responseMessage)
			return
		}

		// Token is valid, extract the email from claims
		email := claims.Email
		if email == "" {
			responseMessage := map[string]string{
				"message": "Unauthorized: Email not found in token",
			}
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(responseMessage)
			return
		}

		// Find user in database
		var user models.User
		filter := bson.M{"email": email}
		var userCollection = db.Client.Database("docs").Collection("users")

		err = userCollection.FindOne(context.Background(), filter).Decode(&user)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				responseMessage := map[string]string{
					"message": "User not found",
				}
				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(responseMessage)
				return
			}
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Store the user ID and email in context for next handler
		ctx := context.WithValue(r.Context(), "userId", user.Id)
		ctx = context.WithValue(ctx, "email", user.Email)

		// Pass to the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
