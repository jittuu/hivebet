package hivebet

import (
	"net/http"

	"golang.org/x/net/context"
	"google.golang.org/appengine"
)

// AppHandler provide wrapper for http.Handler interface and it allows to return error
type AppHandler func(context.Context, http.ResponseWriter, *http.Request) error

func (h AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if err := h(c, w, r); err != nil {
		panic(err)
	}
}
