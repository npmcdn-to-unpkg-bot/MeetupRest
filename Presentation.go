package MeetupRest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
)

const datastorePresentationskind = "Presentations"

type Presentation struct {
	Title       string
	Description string
	Speaker     string
	Voters      []string
}

// Get the handler which contains all the presentation handling routes and the corresponding handlers.
func RegisterPresentationRoutes(m *mux.Router) error {
	if m == nil {
		errors.New("m may not be nil when registering presentation routes")
	}
	m.HandleFunc("/", getPresentation).Methods("GET")
	m.HandleFunc("/", addPresentation).Methods("POST")
	m.HandleFunc("/form/add", addPresentationForm).Methods("GET")

	return nil
}

func getPresentation(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	params, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		log.Errorf(ctx, "Can't parse query: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	q := datastore.NewQuery(datastorePresentationskind).Limit(1)

	if title, ok := params["title"]; ok == true {
		q = q.Filter("Title=", title[0])
	}

	if speaker, ok := params["speaker"]; ok == true {
		q = q.Filter("Speaker=", speaker[0])
	}

	if description, ok := params["description"]; ok == true {
		q = q.Filter("Description=", description[0])
	}

	newCtx, _ := context.WithTimeout(ctx, time.Second*2)
	t := q.Run(newCtx)
	myPresentation := Presentation{}
	_, err = t.Next(&myPresentation)
	// No speaker retrieved
	if err == datastore.Done {
		fmt.Fprint(w, "No Presentation found.")
		return
	}
	// Some other error
	if err != nil {
		log.Errorf(ctx, "Can't get presentation: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(&myPresentation)
	if err != nil {
		log.Errorf(ctx, "Failed to serialize presentation: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	io.Copy(w, bytes.NewReader(data))
}

func addPresentation(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	err := r.ParseForm()
	if err != nil {
		log.Errorf(ctx, "Couldn't parse form: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p := &Presentation{}

	decoder := schema.NewDecoder()
	decoder.Decode(p, r.PostForm)

	if p.Title == "" || p.Speaker == "" || p.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Fields Title, Speaker and Description are mandatory!")
		return
	}

	key := datastore.NewKey(ctx, datastorePresentationskind, "", 0, nil)
	newCtx, _ := context.WithTimeout(ctx, time.Second*2)
	id, err := datastore.Put(newCtx, key, p)
	if err != nil {
		log.Errorf(ctx, "Can't create datastore object: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "%v", id.IntID())
}

func addPresentationForm(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>Adding Presentation Form</h1>"+
		"<form action=\"/presentation/\" method=\"POST\">"+
		"Title: <input type=\"text\" name=\"Title\"><br>"+
		"Description: <textarea name=\"Description\"></textarea><br>"+
		"Speaker: <input type=\"text\" name=\"Speaker\"><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>")
}
