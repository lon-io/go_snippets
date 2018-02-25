package main

import (
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseFiles("index.gohtml"))
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/dog", dog)
	http.HandleFunc("/me", you)

	http.ListenAndServe(":8080", nil)
}

func home(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, "HOME")
}

func dog(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, "DOG")
}

func you(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, "YOU")
}
