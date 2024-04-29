package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	supa "github.com/nedpals/supabase-go"
)

var client *supa.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseAnonKey := os.Getenv("SUPABASE_ANON_KEY")

	if supabaseURL == "" || supabaseAnonKey == "" {
		log.Fatal("Error: SUPABASE_URL and SUPABASE_ANON_KEY must be set")
	}

	client = supa.CreateClient(supabaseURL, supabaseAnonKey)
}

func registerUser(w http.ResponseWriter, r *http.Request) {
	var user supa.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	authResponse, err := client.Auth.SignUp(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(authResponse)
}

func loginUser(w http.ResponseWriter, r *http.Request) {
	var user supa.UserCredentials
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	authResponse, err := client.Auth.SignIn(ctx, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(authResponse)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/register", registerUser).Methods("POST")
	r.HandleFunc("/login", loginUser).Methods("POST")
	log.Fatal(http.ListenAndServe(":8000", r))
}
