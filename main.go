package main

import (
	"log"
	"net/http"
)

func main() {
	srv := http.NewServeMux()
	db_dms := DB{
		Host:     "localhost",
		Password: "password",
		Port:     "5432",
		Name:     "dms",
	}
	repo := CreateRepo(db_dms)
	NewRepo(repo)
	Routes(srv)
	log.Fatalln(http.ListenAndServe(":8080", srv))
}
