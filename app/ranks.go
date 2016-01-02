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
	teamLast5Ranks := make([]*TeamRank, len(teams))

	var i int
	for t, tEvents := range teams {
		teamRanks[i] = getTeamRanks(t, tEvents)
		teamLast5Ranks[i] = getTeamRanks(t, tEvents[:5])
		i++
	}

	sort.Sort(ByProfit(teamRanks))
	sort.Sort(ByProfit(teamLast5Ranks))

	return hivebet.RenderTemplate(w,
		&TeamRanksView{TeamAllRank: teamRanks, TeamLast5Rank: teamLast5Ranks},
		"templates/ranks.html")
}

func getTeamRanks(name string, events []*hivebet.Event) *TeamRank {
	rnk := &TeamRank{Name: name}

	for _, e := range events {
		ew := &EventView{e}
		favWL := ew.FavoriteWinLoss()
		if name == ew.Home {
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

	return rnk
}

// TeamRanksView is view model for ranks
type TeamRanksView struct {
	TeamAllRank   []*TeamRank
	TeamLast5Rank []*TeamRank
}

// TeamRank is view model for a team rank
type TeamRank struct {
	Name   string
	Profit float64
}

// ByProfit is sorter for []*TeamRank
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
