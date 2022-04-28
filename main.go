package main

import (
	// "database/sql"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	// "time"
)

// views

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func Meetings(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/meetings.html"))
	tmpl.Execute(w, nil)
}

func MeetingsDetails(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/meeting-details.html"))
	tmpl.Execute(w, nil)
}

func main() {

	// routes & static files

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("static"))

	r.HandleFunc("/", Home)
	r.HandleFunc("/meetings/", Meetings)
	r.HandleFunc("/meetings/details/", MeetingsDetails)

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// server config
	http.ListenAndServe(":80", r)
	fmt.Println("Server Started")
}
