package hivebet

import (
	"fmt"
	"io"

	"github.com/mjibson/goon"

	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/urlfetch"
)

// ImportEvents download csv file and import events into db
type ImportEvents struct {
	Context context.Context
	League  League
	Season  string
}

// Do download csv file and import events into db
func (imp *ImportEvents) Do() error {
	body, err := imp.download()
	if err != nil {
		return err
	}
	defer body.Close()

	events, err := imp.parse(body)
	if err != nil {
		return err
	}

	err = imp.save(events)
	if err != nil {
		return err
	}

	return nil
}

func (imp *ImportEvents) download() (io.ReadCloser, error) {
	url := imp.URL()
	client := urlfetch.Client(imp.Context)

	log.Infof(imp.Context, "[start] getting data from %s", url)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	log.Infof(imp.Context, "[done] getting data from %s", url)

	return resp.Body, nil
}

func (imp *ImportEvents) parse(body io.Reader) ([]*Event, error) {
	log.Infof(imp.Context, "[start] parsing data")
	parser := &CsvParser{File: body}
	events, err := parser.Parse()
	log.Infof(imp.Context, "[done] parsing data with %d events", len(events))

	return events, err
}

func (imp *ImportEvents) save(events []*Event) error {
	query := datastore.NewQuery("Event").
		Filter("League =", imp.League).
		Filter("Season =", imp.Season)

	g := goon.FromContext(imp.Context)
	var existingEvents []*Event
	_, err := g.GetAll(query, &existingEvents)
	if err != nil {
		return err
	}

	var newCount int
	for _, e := range events {
		e.League = imp.League
		e.Season = imp.Season

		found := false
		for _, xe := range existingEvents {
			if e.Home == xe.Home && e.Away == xe.Away {
				found = true
				break
			}
		}

		if !found {
			_, err := g.Put(e)
			if err != nil {
				return err
			}

			newCount++
		}
	}

	log.Infof(imp.Context, "[done] new (%d) events are added", newCount)

	return nil
}

// URL return url to download
func (imp *ImportEvents) URL() string {
	l, s := imp.League, imp.Season
	spath := s[2:4] + s[7:] // 2015-2016 => 1516

	var fname string
	switch l {
	case EPL:
		fname = "E0.csv"
	case SerieA:
		fname = "I1.csv"
	case Bundesliga:
		fname = "D1.csv"
	case LaLiga:
		fname = "SP1.csv"
	}
	return fmt.Sprintf("http://www.football-data.co.uk/mmz4281/%s/%s", spath, fname)
}
