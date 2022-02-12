package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type httpError struct {
	Code string `json:"code"`;
	Msg string `json:"msg"`
}
func Error(rw http.ResponseWriter, status int, code string, msg string) error {
	rw.WriteHeader(status);
	errorRes := httpError{Code: code, Msg: msg}
	res, err := json.Marshal(errorRes)
	if err != nil {
		return err
	}
	rw.Write(res)
	return nil
}
func connectDB()  {
	_, err := db.Exec(`
		DROP TABLE IF EXISTS users;
		CREATE TABLE users (
			login varchar PRIMARY KEY,
			pass varchar,
			role int
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
