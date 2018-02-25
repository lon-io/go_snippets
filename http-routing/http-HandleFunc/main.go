package main

import (
	"io"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Home")
}

func dog(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Dog")
}

func you(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "You")
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/dog", dog)
	http.HandleFunc("/me", you)

	http.ListenAndServe(":8080", nil)
}
