package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("views"))
	http.Handle("/", fs)
	log.Println("listening...")
	http.ListenAndServe(":8080", nil)
}
