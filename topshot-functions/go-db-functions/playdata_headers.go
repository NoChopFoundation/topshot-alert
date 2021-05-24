package playdata

import (
	"database/sql"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type NoOp struct {
}

type DbHandlerParams struct {
	Opts       *GetPlayDataHTTPInit
	Db         *sql.DB
	W          http.ResponseWriter
	Request    *http.Request
	RouteParms httprouter.Params
}
type DbHandlerFunc func(p DbHandlerParams)

type GetPlayDataHTTPInit struct {
	dbConnectionStr                string
	defaultBlockHeightStreamSize   int // How far back we allow /stream/momentEvents/from/HEAD
	defaultCollectorPastSeconds    int
	defaultBlockHeightPlaySetQuery int // How far back query allowed /momentEvents/play/:PlayId/set/:SetId/from/HEAD
}

type PlayScan struct {
	PlayId            int
	NbaSeason         sql.NullString
	TeamAtMomentNBAID sql.NullString
	PlayCategory      sql.NullString
	JerseyNumber      sql.NullString
	PlayerPosition    sql.NullString
	DateOfMoment      sql.NullString
	PlayType          sql.NullString
	FullName          sql.NullString
	PrimaryPosition   sql.NullString
	TeamAtMoment      sql.NullString
}

type Collector struct {
	CollectorId       string
	State             string
	UpdatesInInterval int
	BlockHeight       uint64
	CreatedAt         string
}

type Play struct {
	PlayId            int
	NbaSeason         string
	TeamAtMomentNBAID string
	PlayCategory      string
	JerseyNumber      string
	PlayerPosition    string
	DateOfMoment      string
	PlayType          string
	FullName          string
	PrimaryPosition   string
	TeamAtMoment      string
	Sets              []int
	EditionCounts     []int
}

type SetPlay struct {
	PlayId   int
	FullName string
}

type Set struct {
	SetId     int
	Name      string
	SetSeries int
}

type PlaysInSets struct {
	PlayId       int
	SetId        int
	EditionCount int
}

type PayloadCollector struct {
	Data []Collector
}

type PayLoadPlaysInSets struct {
	Data []PlaysInSets
}

type MomentEvents struct {
	Type         string
	MomentId     int
	BlockHeight  uint64
	PlayId       int
	SerialNumber int
	SetId        int
	SellerAddr   string
	Price        float32
	Created_At   string
}

type PayloadMomentEvents struct {
	Data []MomentEvents
}

type PayloadPlays struct {
	Data []Play
}

type PayloadSets struct {
	Data []Set
}

type PayloadSetWithPlay struct {
	Data []SetPlay
}

const DBTABLE_plays_in_sets = "plays_in_sets"
const DBTABLE_sets = "sets"
const DBTABLE_plays = "plays"
const DBTABLE_moment_events = "moment_events"
const DBTABLE_moment_events_collectors = "moment_events_collectors"
