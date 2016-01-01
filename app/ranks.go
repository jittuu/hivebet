package app

import (
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/jittuu/hivebet"

	"golang.org/x/net/context"
)

func getRanksIndex(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	league := vars["league"]
	season := vars["season"]

	db := &hivebet.Db{ctx}
	events, err := db.GetAllEvents(league, season)
	if err != nil {
		return err
	}

	teams := make(map[string][]*hivebet.Event)
	for _, e := range events {
		teams[e.Home] = append(teams[e.Home], e)
		teams[e.Away] = append(teams[e.Away], e)
	}

	teamRanks := make([]*TeamRank, len(teams))

	var i int
	for t, tEvents := range teams {
		rnk := &TeamRank{Name: t}
		teamRanks[i] = rnk
		i++

		for _, e := range tEvents {
			ew := &EventView{e}
			favWL := ew.FavoriteWinLoss()
			if t == ew.Home {
				if ew.AvgAHOdds.IsHomeFavorite() {
					rnk.Profit += favWL
				} else {
					rnk.Profit -= favWL
				}
			} else {
				if ew.AvgAHOdds.IsHomeFavorite() {
					rnk.Profit -= favWL
				} else {
					rnk.Profit += favWL
				}
			}
		}
	}

	sort.Sort(ByProfit(teamRanks))

	return hivebet.RenderTemplate(w, teamRanks, "templates/ranks.html")
}

type TeamRank struct {
	Name   string
	Profit float64
}

type ByProfit []*TeamRank

func (s ByProfit) Len() int {
	return len(s)
}
func (s ByProfit) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByProfit) Less(i, j int) bool {
	return s[i].Profit > s[j].Profit
}
