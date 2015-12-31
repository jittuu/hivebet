package app

import (
	"net/http"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/jittuu/hivebet"
)

func init() {
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.Handle("/", hivebet.AppHandler(home))

	events := r.PathPrefix("/events").Subrouter()
	events.Handle("/{league}/{season}", hivebet.AppHandler(getEventsIndex)).Methods("GET")
	events.Handle("/update", hivebet.AppHandler(getEventsUpdate)).Methods("GET")
	events.Handle("/update", hivebet.AppHandler(postEventsUpdate)).Methods("POST")

	http.Handle("/", r)
}

func home(c context.Context, w http.ResponseWriter, r *http.Request) error {
	return hivebet.RenderTemplate(w, nil, "templates/home.html")
}
