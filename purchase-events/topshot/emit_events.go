/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	log "github.com/sirupsen/logrus"
)

type EventEmitter struct {
	EmitPriceChanged         func(*QueryTopShot_MomentPurchasedListedChangedEvent)
	EmitPurchased            func(*QueryTopShot_MomentPurchasedListedChangedEvent)
	EmitListed               func(*QueryTopShot_MomentPurchasedListedChangedEvent)
	EmitWithdrawn            func(*QueryTopShot_MomentWithdrawnEvent)
	EmitCutPercentageChanged func(*QueryTopShot_MomentCutPercentageChangedEvent)
	EmitIntervalStats        func(uint16, IntervalStats, error)
}

func EventEmitter_Create(listeners []EventEmitter) *EventEmitter {
	return &EventEmitter{
		EmitPriceChanged: func(event *QueryTopShot_MomentPurchasedListedChangedEvent) {
			for _, listener := range listeners {
				listener.EmitPriceChanged(event)
			}
		},
		EmitPurchased: func(event *QueryTopShot_MomentPurchasedListedChangedEvent) {
			for _, listener := range listeners {
				listener.EmitPurchased(event)
			}
		},
		EmitListed: func(event *QueryTopShot_MomentPurchasedListedChangedEvent) {
			for _, listener := range listeners {
				listener.EmitListed(event)
			}
		},
		EmitWithdrawn: func(event *QueryTopShot_MomentWithdrawnEvent) {
			for _, listener := range listeners {
				listener.EmitWithdrawn(event)
			}
		},
		EmitCutPercentageChanged: func(event *QueryTopShot_MomentCutPercentageChangedEvent) {
			for _, listener := range listeners {
				listener.EmitCutPercentageChanged(event)
			}
		},
		EmitIntervalStats: func(collectorId uint16, current IntervalStats, err error) {
			for _, listener := range listeners {
				listener.EmitIntervalStats(collectorId, current, err)
			}
		},
	}
}

func EventEmitter_DebugConsoleListener() EventEmitter {
	return EventEmitter{
		EmitPriceChanged: func(e *QueryTopShot_MomentPurchasedListedChangedEvent) {
			s := e.Details.MomentDetails()
			log.WithFields(log.Fields{
				"N":  s.Play()["FullName"],
				"SN": s.SerialNumber(),
				"T":  s.SetName(),
				"P":  s.Price(),
			}).Info("Changed")
		},
		EmitPurchased: func(e *QueryTopShot_MomentPurchasedListedChangedEvent) {
			s := e.Details.MomentDetails()
			log.WithFields(log.Fields{
				"N":  s.Play()["FullName"],
				"SN": s.SerialNumber(),
				"T":  s.SetName(),
				"P":  s.Price(),
			}).Info("Purchased")
		},
		EmitListed: func(e *QueryTopShot_MomentPurchasedListedChangedEvent) {
			s := e.Details.MomentDetails()
			log.WithFields(log.Fields{
				"N":  s.Play()["FullName"],
				"SN": s.SerialNumber(),
				"T":  s.SetName(),
				"P":  s.Price(),
			}).Info("Listed")
		},
		EmitWithdrawn: func(event *QueryTopShot_MomentWithdrawnEvent) {
			//TODO
		},
		EmitCutPercentageChanged: func(event *QueryTopShot_MomentCutPercentageChangedEvent) {
			//TODO
		},
		EmitIntervalStats: func(collectorId uint16, current IntervalStats, err error) {
			//TODO
		},
	}
}
