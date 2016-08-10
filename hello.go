package MeetupRest

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

type Human struct {
	Name string
	Age  int
}

func init() {

	m := mux.NewRouter()
	m.HandleFunc("/secure/test", func(w http.ResponseWriter, r *http.Request) {

	})
	m.Handle("/speaker", GetSpeakerHandler())
	m.Handle("/presentation", GetPresentationHandler())
	m.Handle("/addSpeaker", GetFormsHandler())
	s := m.PathPrefix("/meetup").Subrouter()
	err := RegisterMeetupRoutes(s)
	if err != nil {
		panic(err)
	}

	http.Handle("/", m)
}

type dsHandler struct {
}

func (h dsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	variables := mux.Vars(r)
	e := &Human{}

	e.Name = variables["name"]
	age, _ := strconv.Atoi(variables["age"])
	e.Age = age

	k := datastore.NewKey(ctx, "People", "", 0, nil)
	newCtx, _ := context.WithTimeout(ctx, time.Second*2)
	id, err := datastore.Put(newCtx, k, e)
	if err != nil {
		log.Errorf(ctx, "Can't create datastore object: %v", err)
		return
	}

	fmt.Fprintf(w, "Your key: %v", id.IntID())
}
