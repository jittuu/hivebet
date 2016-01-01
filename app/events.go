package app

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jittuu/hivebet"
	"github.com/jittuu/hivebet/oddsUtil"

	"golang.org/x/net/context"
)

func getEventsIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	league := vars["league"]
	season := vars["season"]

	db := &hivebet.Db{ctx}
	events, err := db.GetAllEvents(league, season)
	if err != nil {
		return err
	}

	eventViews := make([]*EventView, len(events))
	for i, e := range events {
		eventViews[i] = &EventView{e}
	}

	return hivebet.RenderTemplate(w, eventViews, "templates/events.html")
}

func getEventsUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	return hivebet.RenderTemplate(w, nil, "templates/update.html")
}

func postEventsUpdate(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	season := r.FormValue("season")
	league := r.FormValue("league")

	imp := &hivebet.ImportEvents{
		Context: ctx,
		League:  hivebet.League(league),
		Season:  season,
	}

	err := imp.Do()
	if err != nil {
		return err
	}

	http.Redirect(w, r, "/events/update", http.StatusFound)
	return nil
}

// EventView represent event view model
type EventView struct {
	*hivebet.Event
}

// FormatMyanmarOdds return formatted Myanmar odds
func (e *EventView) FormatMyanmarOdds() string {
	goals, odds := e.MyanmarOdds()

	return fmt.Sprintf("%d%+.0f", goals, odds*100)
}

// MyanmarOdds returns Myanmar odds
func (e *EventView) MyanmarOdds() (goals int, odds float64) {
	var handicap, euro float64
	if e.AvgAHOdds.IsHomeFavorite() {
		handicap, euro = e.AvgAHOdds.Handicap, e.AvgAHOdds.Home
	} else {
		handicap, euro = e.AvgAHOdds.Handicap*-1, e.AvgAHOdds.Away
	}

	return oddsUtil.ConvertEuroToMyanmar(handicap, euro)
}

// Handicap return AH handicap
func (e *EventView) Handicap() float64 {
	return e.AvgAHOdds.Handicap
}

// HomeOdds return home odds in malay odds
func (e *EventView) HomeOdds() float64 {
	hk := oddsUtil.ConvertToHK(e.AvgAHOdds.Home)
	return oddsUtil.ConvertToMalay(hk)
}

// AwayOdds return away odds in malay odds
func (e *EventView) AwayOdds() float64 {
	hk := oddsUtil.ConvertToHK(e.AvgAHOdds.Away)
	return oddsUtil.ConvertToMalay(hk)
}

func (e *EventView) FavoriteWinLoss() float64 {
	goals, odds := e.MyanmarOdds()
	var diff int
	if e.AvgAHOdds.IsHomeFavorite() {
		diff = (e.HGoal - goals - e.AGoal)
	} else {
		diff = (e.AGoal - goals - e.HGoal)
	}

	switch {
	case diff > 0:
		return 1
	case diff < 0:
		return -1
	default:
		return odds
	}
}
