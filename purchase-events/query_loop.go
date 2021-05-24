/**
 * author: NoChopFoundation@gmail.com
 * Inspired by https://medium.com/@eric.ren_51534/polling-nba-top-shot-p2p-market-purchase-events-from-flow-blockchain-using-flow-go-sdk-3ec80119e75f
 */
package main

import (
	"context"
	"database/sql"
	"flag"
	"time"

	"github.com/NoChopFoundation/topshot-alert/purchase-events/topshot"
	log "github.com/sirupsen/logrus"
)

// How far back do we ask the query service for details
// The query service will only have a certain amount of recent blocks
// available to run scripts against the state at that point in time
const MAX_PREV_BLOCKS_QUERIED = 100

// Main entry point for running the TopShot query and database populate.
// We persist our processing state in JobState file incase we
// shutdown and need to restart.
func main() {
	topshot.Initialize()

	isDatabaseModePtr := flag.Bool("useDb", false, "attempt to use a database, furthur env variables are required")
	dumpSetPlaysPtr := flag.Bool("dumpSchema", false, "console print database schema and latest moment data in TopShot contract")
	interationsPtr := flag.Int("i", 1, "the number of iterations to execute before exiting (0 forever)")
	quietPtr := flag.Bool("quiet", false, "do not display events to the console")
	flag.Parse()

	var eventEmitters []topshot.EventEmitter
	var endPoint *topshot.Configuration
	var initErr error
	if *isDatabaseModePtr {
		endPoint, initErr = topshot.Configuration_MainNet_withMySql()
		panicOnError(initErr)
	} else {
		endPoint, initErr = topshot.Configuration_MainNet()
		panicOnError(initErr)
	}

	if !*quietPtr {
		eventEmitters = append(eventEmitters, topshot.EventEmitter_DebugConsoleListener())
	}

	var db *sql.DB
	if len(endPoint.MySqlConnection) > 0 {
		var dbErr error
		db, dbErr = topshot.MySQL_Database_Create(context.Background(), endPoint)
		panicOnError(dbErr)
		eventEmitters = append(eventEmitters, topshot.MySQL_EventEmitter_Updater_Listener(db, true))
	} else {
		db = nil
	}

	if *dumpSetPlaysPtr {
		topshot.TopShotUtil_dumpSetPlaysData(endPoint)
	} else {
		prevJobstate := topshot.LoadJobState()
		interations := *interationsPtr
		process := topshot.Process{
			Config:             endPoint,
			ErrorControl:       topshot.Backoff_Create(topshot.Backoff_AlgorithmNoBlast(topshot.Backoff_ConsoleLogger())),
			QueryApi:           nil,
			LastBlockProcessed: prevJobstate.LastBlockProcessed,
			Emitter:            topshot.EventEmitter_Create(eventEmitters),
			CurrentInterval:    topshot.IntervalStats{},
		}
		if interations == 0 {
			for i := 0; true; i++ {
				err := loop(prevJobstate, &process)
				loopCompleted(&process, err)
			}
		} else {
			for i := 0; i < interations; i++ {
				err := loop(prevJobstate, &process)
				loopCompleted(&process, err)
			}
		}
	}
}

func loop(prevJobstate topshot.JobState, process *topshot.Process) error {
	process.CurrentInterval.ResetInterval()

	var err error
	if process.QueryApi == nil {
		process.QueryApi, err = topshot.Connection(context.Background(), process.Config)
		if err != nil {
			return err
		}
	}

	log.WithFields(log.Fields{
		"prevJobstate.LastBlockProcessed": prevJobstate.LastBlockProcessed,
	}).Info("loop main")

	latestBlock, err := process.QueryApi.GetLatestBlock()
	if err != nil {
		return err
	}

	process.CurrentInterval.BlockHeight = latestBlock.Height

	var startHeight uint64
	if latestBlock.Height < process.LastBlockProcessed {
		// We previously processed a block in the future?
		// Must be they are replaying/reloading the blockchain?
		startHeight = latestBlock.Height - MAX_PREV_BLOCKS_QUERIED
	} else {
		if (latestBlock.Height - process.LastBlockProcessed) < MAX_PREV_BLOCKS_QUERIED {
			startHeight = process.LastBlockProcessed
		} else {
			startHeight = latestBlock.Height - MAX_PREV_BLOCKS_QUERIED
		}
	}
	if startHeight == latestBlock.Height {
		// Skip, wait for more blocks
		return nil
	}

	log.WithFields(log.Fields{
		"latestBlock.Height": (latestBlock.Height - startHeight),
	}).Info("calling query")

	err = topshot.QueryTopShot_ProcessEvents(process, startHeight, latestBlock.Height)
	if err != nil {
		return err
	}

	process.LastBlockProcessed = latestBlock.Height
	prevJobstate.LastBlockProcessed = latestBlock.Height
	topshot.SaveJobState(&prevJobstate)

	return nil
}

func loopCompleted(process *topshot.Process, errArg error) {
	if errArg != nil {
		sleepMs, err := topshot.Backoff_HandleReturn(process.ErrorControl, errArg)
		if err != nil {
			// Error handler returned error, this one is fatal
			panic(err)
		} else {
			time.Sleep(time.Duration(sleepMs) * time.Microsecond)
		}
	} else {
		// No errors, maybe see if we are processing max query length
		// and don't delay, or delay some if we are retriving less than max
		time.Sleep(15000 * time.Millisecond)
	}

	// After the sleep internal, report all activity which occurred during the internal
	currentInterval := process.CurrentInterval
	process.Emitter.EmitIntervalStats(process.Config.CollectorId, currentInterval, errArg)
}

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}
