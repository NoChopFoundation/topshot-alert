/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import "github.com/onflow/flow-go-sdk/client"

type IntervalStats struct {
	EventsProcessed int
	BlockHeight     uint64
}

func (interval *IntervalStats) ResetInterval() {
	interval.EventsProcessed = 0
	interval.BlockHeight = 0
}

func (interval *IntervalStats) CountBlockEvents(blockEvents []client.BlockEvents) {
	for _, block := range blockEvents {
		interval.EventsProcessed += len(block.Events)
	}
}

type Process struct {
	Config             *Configuration
	QueryApi           *TopshotQueryApi
	ErrorControl       *BackoffControl
	LastBlockProcessed uint64
	Emitter            *EventEmitter
	CurrentInterval    IntervalStats
}
