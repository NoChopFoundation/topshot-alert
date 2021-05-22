/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

const MOMENT_EVENTS_TYPE_FIELD_Purchased = "P"
const MOMENT_EVENTS_TYPE_FIELD_Listed = "L"

type MYSQL_IntervalStats struct {
	lock           *sync.Mutex
	successInserts int
	failedInserts  int
}

func (interval *MYSQL_IntervalStats) init() {
	interval.lock = &sync.Mutex{}
	interval.successInserts = 0
	interval.failedInserts = 0
}

func (interval *MYSQL_IntervalStats) add(add MYSQL_IntervalStats) {
	interval.lock.Lock()
	{
		interval.successInserts += add.successInserts
		interval.failedInserts += add.failedInserts
	}
	interval.lock.Unlock()
}

func (interval *MYSQL_IntervalStats) copyThenReset() MYSQL_IntervalStats {
	result := MYSQL_IntervalStats{}

	interval.lock.Lock()
	{
		result.successInserts = interval.successInserts
		result.failedInserts = interval.failedInserts

		interval.successInserts = 0
		interval.failedInserts = 0
	}
	interval.lock.Unlock()

	return result
}

func MySQL_EventEmitter_Updater_Listener(db *sql.DB, verbose bool) EventEmitter {
	var stats = &MYSQL_IntervalStats{}
	stats.init()
	return EventEmitter{
		EmitPriceChanged: func(e *QueryTopShot_MomentPurchasedListedChangedEvent) {
			log.WithFields(log.Fields{
				"MomentID": e.Details.MomentID(),
				"Price":    e.Event.Price(),
			}).Warn("DETECTED PRICE CHANGE")
			//TODO, see if we ever detect these
		},
		EmitPurchased: func(e *QueryTopShot_MomentPurchasedListedChangedEvent) {
			s := e.Details.MomentDetails()
			doMomentInsert(db, stats, verbose, "purchase "+strconv.Itoa(int(e.Details.MomentID())),
				getMomentInsert(), MOMENT_EVENTS_TYPE_FIELD_Purchased, int(e.Details.MomentID()),
				e.BlockInfo.Height, s.PlayID(), s.SerialNumber(), s.SetID(), e.Event.Price(), e.BlockInfo.BlockID.Hex(),
				e.Event.Seller().Hex())
		},
		EmitListed: func(e *QueryTopShot_MomentPurchasedListedChangedEvent) {
			s := e.Details.MomentDetails()
			doMomentInsert(db, stats, verbose, "listed "+strconv.Itoa(int(e.Details.MomentID())),
				getMomentInsert(), MOMENT_EVENTS_TYPE_FIELD_Listed, int(e.Details.MomentID()),
				e.BlockInfo.Height, s.PlayID(), s.SerialNumber(), s.SetID(), e.Event.Price(), e.BlockInfo.BlockID.Hex(),
				e.Event.Seller().Hex())
		},
		EmitWithdrawn: func(event *QueryTopShot_MomentWithdrawnEvent) {
			//TODO
		},
		EmitCutPercentageChanged: func(event *QueryTopShot_MomentCutPercentageChangedEvent) {
			//TODO
		},
		EmitIntervalStats: func(collectorId uint16, current IntervalStats, err error) {
			var sqlStats = stats.copyThenReset()

			var state string
			if err != nil { //
				state = "E" // Error
			} else {
				if sqlStats.failedInserts > 0 {
					// Must be issue with the INSERTs?
					state = "E"
				} else {
					state = "U" // Up
				}
			}
			insertStmt, err := db.Query("INSERT into `moment_events_collectors` ( `CollectorId`, `State`, `UpdatesInInterval`, `BlockHeight`) values (?, ?, ?, ?)",
				collectorId, state, sqlStats.successInserts, current.BlockHeight)
			if err != nil {
				fmt.Printf("SQL INSERT %v \n", err)
			} else {
				insertStmt.Close()
			}
		},
	}
}

func doMomentInsert(db *sql.DB, stats *MYSQL_IntervalStats, verbose bool, typeDebug string, insertSql string, insertArgs ...interface{}) {
	insertStmt, err := db.Query(insertSql, insertArgs...)
	if err != nil {
		log.WithFields(log.Fields{
			"typeDebug": typeDebug,
		}).Errorf("SQL INSERT ERROR %v", err)
		stats.add(MYSQL_IntervalStats{failedInserts: 1})
	} else {
		stats.add(MYSQL_IntervalStats{successInserts: 1})
		insertStmt.Close()
		if verbose {
			log.WithFields(log.Fields{
				"typeDebug": typeDebug,
			}).Debug("SQL INSERT OK")
		}
	}
}

func getMomentInsert() string {
	return "INSERT into `moment_events` (`type`, `MomentId`, `BlockHeight`, `PlayId`, `SerialNumber`, `SetId`, `Price`, `BlockId`, `SellerAddr`) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE created_at=NOW()"
}

// TODO replace with code we used previously in functions
func MySQL_Database_Create(ctx context.Context, config *Configuration) (*sql.DB, error) {
	db, err := sql.Open("mysql", config.MySqlConnection)
	if err != nil {
		return nil, err
	}
	//defer db.Close() TODO defer in main
	return db, nil
}
