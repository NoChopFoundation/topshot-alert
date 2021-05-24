package playdata

import (
	"database/sql"
	"sync"
	"time"
)

// GCP Functions may keep global variables across frequent access
type GlobalCache struct {
	mutex          *sync.Mutex
	playsCache     *PayloadPlays
	playsCacheTime time.Time

	// Attempt to reuse database connection
	dbPool *sql.DB
}

func withLock(g *GlobalCache, logic func()) {
	g.mutex.Lock()
	logic()
	g.mutex.Unlock()
}

func GlobalCache_GetDb(g *GlobalCache) *sql.DB {
	var dbCache *sql.DB
	withLock(g, func() {
		if g.dbPool != nil {
			dbCache = g.dbPool
		}
	})
	return dbCache
}

// Possible we should execute the DB open under lock to avoid leak
func GlobalCache_SetDb(g *GlobalCache, db *sql.DB) {
	withLock(g, func() {
		g.dbPool = db
	})
}

func GlobalCacheAvailable(g *GlobalCache) *PayloadPlays {
	var playsCache *PayloadPlays
	withLock(g, func() {
		if g.playsCache != nil && g.playsCacheTime.After(time.Now()) {
			playsCache = g.playsCache
		}
	})
	return playsCache
}

func GlobalCacheSet(g *GlobalCache, playsCache *PayloadPlays) {
	withLock(g, func() {
		g.playsCache = playsCache
		g.playsCacheTime = time.Now().Add(1 * time.Hour) // Could make this env
	})
}

var Cache *GlobalCache
