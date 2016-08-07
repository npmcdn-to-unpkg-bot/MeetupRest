package MeetupRest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"io"
	"net/http"
	"net/url"
	"time"
)

const kindSpeakers = "Speakers"

type Speaker struct {
	Name    string
	Surname string
	About   string
	Email   string
	Company string
}

func GetSpeakerHandler() http.Handler {
	m := mux.NewRouter()
	m.HandleFunc("/speaker", getSpeaker).Methods("GET")
	m.HandleFunc("/speaker", addSpeaker).Methods("POST")

	return m
}

func getSpeaker(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Errorf(ctx, "Can't parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Can't parse query: %v", err)
		return
	}

	q := datastore.NewQuery(kindSpeakers).Limit(1)

	if name, ok := params["name"]; ok == true {
		q = q.Filter("Name=", name[0])
	}

	if surname, ok := params["surname"]; ok == true {
		q = q.Filter("Surname=", surname[0])
	}

	if email, ok := params["email"]; ok == true {
		q = q.Filter("Email=", email[0])
	}

	newCtx, _ := context.WithTimeout(ctx, time.Second*2)
	t := q.Run(newCtx)

	mySpeaker := Speaker{}
	_, err = t.Next(&mySpeaker)
	// No speaker retrieved
	if err == datastore.Done {
		fmt.Fprint(w, "No speaker found.")
		return
	}
	// Some other error.
	if err != nil {
		log.Errorf(ctx, "Can't get speaker: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Can't get speaker: %v", err)
		return
	}
	data, err := json.Marshal(&mySpeaker)
	if err != nil {
		log.Errorf(ctx, "Failed to serialize speaker: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Failed to serialize speaker: %v", err)
		return
	}
	io.Copy(w, bytes.NewBuffer(data))
}

func addSpeaker(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "Couldn't parse form: Name")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Couldn't parse form: Name")
		return
	}

	s := &Speaker{}

	decoder := schema.NewDecoder()
	decoder.Decode(s, r.PostForm)

	if s.Name == "" || s.Surname == "" || s.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Name, surname and email are mandatory.")
		return
	}

	key := datastore.NewKey(ctx, kindSpeakers, "", 0, nil)
	newCtx, _ := context.WithTimeout(ctx, time.Second*2)
	id, err := datastore.Put(newCtx, key, s)
	if err != nil {
		log.Errorf(ctx, "Can't create datastore object: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Can't create datastore object: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v", id.IntID())
}
