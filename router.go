package main

import (
	"authApi/types"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/jwtauth"
	jwtCreator "github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)


func register( rw http.ResponseWriter, r *http.Request){
	var body types.RegisterDTO;
	err := json.NewDecoder(r.Body).Decode(&body);
	if err != nil {
		Error(rw, 400, "MALFORMED_BODY", "Body couldn't be parsed")
		return
	}
	if(len(body.Login)<3 || len(body.Password)<6 || body.Role == 0){
		Error(rw, 400, "BAD_REQUEST", "Login should be at least 3 chars. Password should be at least 6. Roles start from 1")
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.MinCost)
	if err != nil {
		Error(rw, 500, "INTERNAL_SERVER_ERROR", "")
		return
	}
	_, err = db.Exec("INSERT INTO users VALUES	(?,?,?)", body.Login, hash, body.Role);
	if err != nil {
		Error(rw, 409, "EXISTING_USER","This login is used")
		return
	}
}

func login(rw http.ResponseWriter, r *http.Request){
	var body types.LoginDTO;
	json.NewDecoder(r.Body).Decode(&body);
	userRow := db.QueryRow("SELECT login, pass, role FROM users WHERE login=?", body.Login);
	if userRow.Err() != nil {
		Error(rw, 500, "INTERNAL_SERVER_ERROR", "")
		return
	}

	var dbUser types.LoginDB;
	errRow := userRow.Scan(&dbUser.Login, &dbUser.Hash, &dbUser.Role)
	errHash := bcrypt.CompareHashAndPassword([]byte(dbUser.Hash), []byte(body.Password))
	if errRow != nil || errHash != nil {
		Error(rw, 401, "WRONG_LOGIN", "Wrong login or password")
		return
}
	token := jwtCreator.New(jwtCreator.GetSigningMethod("HS256"))
	token.Claims = 
		jwtCreator.MapClaims{
			"login": body.Login,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(time.Minute).Unix(),
			"role": dbUser.Role,
		}
	tokenString, err := token.SignedString(JWT_SECRET)
	if err != nil {
		Error(rw, 500, "INTERNAL_SERVER_ERROR", "")
	}
	rw.Write([]byte(tokenString))
}

func check(rw http.ResponseWriter, r *http.Request){
	 _, token, _ := jwtauth.FromContext(r.Context())
	// if !success {
	// 	Error(rw, 404, "NO_ROLE", "JWT doesn't contain a role")
	// 	return
	// }
	tokenstr, err := json.Marshal(token)
	if err != nil {
		Error(rw, 500, "INTERNAL_SERVER_ERROR", "")
		return
	}
	rw.Write([]byte(tokenstr))
}

