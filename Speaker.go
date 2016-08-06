package MeetupRest

import (
	"net/http"
	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"net/url"
	"google.golang.org/appengine/log"
	"fmt"
)

type Speaker struct {
	Name          string
	Surname       string
	About         string
	Email         string
	Company       string
	Presentations []string
}

func GetSpeakerHandler() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc(/...)
	m.Methods("GET").HandleFunc("/author/")

	return m
}

func getSpeaker(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Errorf(ctx, "Can't parse query: %v", err)
		fmt.Fprintf(w, "Can't parse query: %v", err)
		return
	}


	name 	:= params["Name"]
	surname := params["Surname"]
	email 	:= params["Email"]

	// Fillter Query

	q := datastore.NewQuery("Speaker").Filter("Name =", name).Filter("Surname =", surname).Filter("Email =",)
}
