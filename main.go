package main

import (
	// "database/sql"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// types
type ContactDetails struct {
	Name    string
	Email   string
	Subject string
	Message string
}

func (c ContactDetails) validateForm() (string, bool) {
	name := c.Name
	email := c.Email
	subject := c.Subject
	message := c.Message
	if name == "" {
		return "name is empty", false
	}
	if email == "" {
		return "email is empty", false
	}
	if subject == "" {
		return "subject is empty", false
	}
	if message == "" {
		return "message is empty", false
	}
	return "valid", true
}

func (c *ContactDetails) validateEmail() (string, bool) {
	email := c.Email
	if emailRegex.MatchString(email) {
		return "valid", true
	} else {
		return "email is invalid", false
	}

}

// views

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("templates/index.html"))
	tmpl.ExecuteWriter(pongo2.Context{"title": "Education - Golang-based web app"}, w)
}

func Meetings(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("templates/meetings.html"))
	tmpl.ExecuteWriter(pongo2.Context{"title": "Education - Meetings list"}, w)
}

func MeetingsDetails(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("templates/meeting-details.html"))
	tmpl.ExecuteWriter(pongo2.Context{"title": "Education - Meeting details"}, w)
}
func ContactView(w http.ResponseWriter, r *http.Request) {
	tmpl := pongo2.Must(pongo2.FromFile("templates/contact-response.html"))
	formFields := ContactDetails{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Message: r.FormValue("message"),
	}
	value, state := formFields.validateForm()
	if !state {
		tmpl.ExecuteWriter(pongo2.Context{"success": false, "error": value}, w)
	}
	evalue, estate := formFields.validateEmail()
	if !estate {
		tmpl.ExecuteWriter(pongo2.Context{"success": false, "error": evalue}, w)
	}
	file, err := os.OpenFile("user_messages.csv", os.O_WRONLY|os.O_APPEND, 0644)

	if err != nil {
		fmt.Printf("failed opening file: %s", err)
	}

	defer file.Close()

	fmf := formFields
	_, err2 := file.WriteString(strings.Join([]string{fmf.Name, fmf.Email, fmf.Subject, fmf.Message}, ","))

	if err2 != nil {
		fmt.Println(err2)
	}
	tmpl.ExecuteWriter(pongo2.Context{"success": true, "error": "", "SuccessMessage": "message was sent successfully."}, w)

}

func main() {

	// routes & static files

	r := mux.NewRouter()
	fs := http.FileServer(http.Dir("static"))

	r.HandleFunc("/", Home).Name("Home")
	r.HandleFunc("/meetings/", Meetings).Name("meetings")
	r.HandleFunc("/meetings/details/", MeetingsDetails).Name("meetings-details") 
	r.HandleFunc("/contact/", ContactView).Methods("POST").Name("contact")

	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs)).Name("static")

	// server config
	fmt.Println("Server Started")
	http.ListenAndServe(":80", r)
}
