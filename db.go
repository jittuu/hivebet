package hivebet

import (
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type Db struct {
	Context context.Context
}

func (db *Db) GetAllEvents(league, season string) ([]*Event, error) {
	g := goon.FromContext(db.Context)
	query := datastore.NewQuery("Event").
		Filter("League =", league).
		Filter("Season =", season).
		Order("-StartTime")

	var events []*Event
	_, err := g.GetAll(query, &events)
	if err != nil {
		return nil, err
	}

	return events, nil
}
