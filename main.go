package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/joho/godotenv"
)

var (
	JWT_SECRET []byte
	jwt *jwtauth.JWTAuth
  db *sql.DB
)
func main() {
	var err error;

	godotenv.Load(".env");
	JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))
	jwt = jwtauth.New("HS256", JWT_SECRET, nil);
	
	db, err = sql.Open("sqlite3", "users.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close();
	connectDB()

	router := chi.NewRouter();
	router.Use(middleware.Logger)
	router.Post("/register", register)
	router.Post("/login", login)
	router.With(jwtauth.Verifier(jwt), jwtauth.Authenticator).Get("/check", check);
	http.ListenAndServe("127.0.0.1:3323", router)
}