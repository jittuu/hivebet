package hivebet

import (
	"encoding/csv"
	"io"
	"strconv"
	"time"
)

// CsvParser parse csv file to events
type CsvParser struct {
	File io.Reader
}

// Parse csv file to events
func (p *CsvParser) Parse() ([]*Event, error) {
	r := csv.NewReader(p.File)
	r.TrailingComma = true
	lines, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	var events []*Event
	h := csvHeader(lines[0]).HeaderIndex()

	for _, line := range lines[1:] {
		startTime, _ := time.Parse("02/01/06", line[h.StartTime])
		hGoal, err := strconv.ParseInt(line[h.HGoal], 10, 32)
		if err != nil {
			continue
		}
		aGoal, err := strconv.ParseInt(line[h.AGoal], 10, 32)
		if err != nil {
			continue
		}
		mxH, _ := strconv.ParseFloat(line[h.MxH], 64)
		avH, _ := strconv.ParseFloat(line[h.AvH], 64)
		mxD, _ := strconv.ParseFloat(line[h.MxD], 64)
		avD, _ := strconv.ParseFloat(line[h.AvD], 64)
		mxA, _ := strconv.ParseFloat(line[h.MxA], 64)
		avA, _ := strconv.ParseFloat(line[h.AvA], 64)
		hAH, _ := strconv.ParseFloat(line[h.AHh], 64)
		mxAHH, _ := strconv.ParseFloat(line[h.MxAHH], 64)
		mxAHA, _ := strconv.ParseFloat(line[h.MxAHA], 64)
		avAHH, _ := strconv.ParseFloat(line[h.AvAHH], 64)
		avAHA, _ := strconv.ParseFloat(line[h.AvAHA], 64)
		hOU := 2.5
		mxOUO, _ := strconv.ParseFloat(line[h.MxOUO], 64)
		mxOUU, _ := strconv.ParseFloat(line[h.MxOUU], 64)
		avOUO, _ := strconv.ParseFloat(line[h.AvOUO], 64)
		avOUU, _ := strconv.ParseFloat(line[h.AvOUU], 64)

		event := &Event{
			StartTime: startTime,
			Home:      line[h.Home],
			Away:      line[h.Away],
			HGoal:     int(hGoal),
			AGoal:     int(aGoal),
			MaxOdds: MatchOdds{
				Home: mxH,
				Draw: mxD,
				Away: mxA,
			},
			AvgOdds: MatchOdds{
				Home: avH,
				Draw: avD,
				Away: avA,
			},
			MaxAHOdds: AHOdds{
				Handicap: hAH,
				Home:     mxAHH,
				Away:     mxAHA,
			},
			AvgAHOdds: AHOdds{
				Handicap: hAH,
				Home:     avAHH,
				Away:     avAHA,
			},
			MaxOUOdds: OUOdds{
				Handicap: hOU,
				Over:     mxOUO,
				Under:    mxOUU,
			},
			AvgOUOdds: OUOdds{
				Handicap: hOU,
				Over:     avOUO,
				Under:    avOUU,
			},
		}

		events = append(events, event)
	}

	return events, nil
}

type csvHeader []string

func (header csvHeader) HeaderIndex() *headerIndex {
	h := &headerIndex{}
	for i, col := range header {
		switch col {
		case "Date":
			h.StartTime = i
		case "HomeTeam":
			h.Home = i
		case "AwayTeam":
			h.Away = i
		case "FTHG":
			h.HGoal = i
		case "FTAG":
			h.AGoal = i
		case "BbMxH":
			h.MxH = i
		case "BbMxD":
			h.MxD = i
		case "BbMxA":
			h.MxA = i
		case "BbAvH":
			h.AvH = i
		case "BbAvD":
			h.AvD = i
		case "BbAvA":
			h.AvA = i
		case "BbAHh":
			h.AHh = i
		case "BbMxAHH":
			h.MxAHH = i
		case "BbMxAHA":
			h.MxAHA = i
		case "BbAvAHH":
			h.AvAHH = i
		case "BbAvAHA":
			h.AvAHA = i
		case "BbMx>2.5":
			h.MxOUO = i
		case "BbMx<2.5":
			h.MxOUU = i
		case "BbAv>2.5":
			h.AvOUO = i
		case "BbAv<2.5":
			h.AvOUU = i
		}
	}

	return h
}

type headerIndex struct {
	StartTime     int
	Home, Away    int
	HGoal, AGoal  int
	MxH, MxD, MxA int
	AvH, AvD, AvA int
	AHh           int
	MxAHH, MxAHA  int
	AvAHH, AvAHA  int
	OUh           int
	MxOUO, MxOUU  int
	AvOUO, AvOUU  int
}
