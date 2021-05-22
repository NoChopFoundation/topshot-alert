/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	log "github.com/sirupsen/logrus"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
)

type OnSaleMomentDetailsAvailableFunc func(*client.BlockEvents, MomentPurchasedListedChangedEvent, SaleMomentBatch)

type QueryTopShot_BlockInfo struct {
	Height  uint64
	BlockID flow.Identifier
}

type QueryTopShot_MomentPurchasedListedChangedEvent struct {
	BlockInfo QueryTopShot_BlockInfo
	Event     MomentPurchasedListedChangedEvent
	Details   SaleMomentBatch
	Process   *Process
}

type QueryTopShot_MomentWithdrawnEvent struct {
	BlockInfo QueryTopShot_BlockInfo
	Event     MomentWithdrawnEvent
	Process   *Process
}

type QueryTopShot_MomentCutPercentageChangedEvent struct {
	BlockInfo QueryTopShot_BlockInfo
	Event     MomentCutPercentageChangedEvent
	Process   *Process
}

func QueryTopShot_ProcessEvents(process *Process, startHeight uint64, endHeight uint64) error {
	listedBlockEvents, err := process.QueryApi.GetMomentListedEvents(GetMomentEvents_Arg{
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})
	if err != nil {
		return err
	} else {
		process.CurrentInterval.CountBlockEvents(listedBlockEvents) //TODO use acti
		go processMomentPurchasedListedChangedEvents(process, listedBlockEvents, 0,
			func(blockEvent *client.BlockEvents, event MomentPurchasedListedChangedEvent, sale SaleMomentBatch) {
				process.Emitter.EmitListed(&QueryTopShot_MomentPurchasedListedChangedEvent{
					Process: process,
					Event:   event,
					BlockInfo: QueryTopShot_BlockInfo{
						Height:  blockEvent.Height,
						BlockID: blockEvent.BlockID},
					Details: sale})
			})
	}

	purchasedBlockEvents, err := process.QueryApi.GetMomentPurchasedEvents(GetMomentEvents_Arg{
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})
	if err != nil {
		return err
	} else {
		process.CurrentInterval.CountBlockEvents(purchasedBlockEvents)
		go processMomentPurchasedListedChangedEvents(process, purchasedBlockEvents, 1,
			func(blockEvent *client.BlockEvents, event MomentPurchasedListedChangedEvent, sale SaleMomentBatch) {
				process.Emitter.EmitPurchased(&QueryTopShot_MomentPurchasedListedChangedEvent{
					Process: process,
					Event:   event,
					BlockInfo: QueryTopShot_BlockInfo{
						Height:  blockEvent.Height,
						BlockID: blockEvent.BlockID},
					Details: sale})
			})
	}

	changedBlockEvents, err := process.QueryApi.GetMomentPriceChangedEvents(GetMomentEvents_Arg{
		StartHeight: startHeight,
		EndHeight:   endHeight,
	})
	if err != nil {
		return err
	} else {
		process.CurrentInterval.CountBlockEvents(changedBlockEvents)
		go processMomentPurchasedListedChangedEvents(process, changedBlockEvents, 0,
			func(blockEvent *client.BlockEvents, event MomentPurchasedListedChangedEvent, sale SaleMomentBatch) {
				process.Emitter.EmitPriceChanged(&QueryTopShot_MomentPurchasedListedChangedEvent{
					Process: process,
					Event:   event,
					BlockInfo: QueryTopShot_BlockInfo{
						Height:  blockEvent.Height,
						BlockID: blockEvent.BlockID},
					Details: sale})
			})
	}

	return nil
}

func processMomentPurchasedListedChangedEvents(process *Process, blockEvents []client.BlockEvents, executeScriptAtOffset uint64, callBack OnSaleMomentDetailsAvailableFunc) {
	for _, blockEvent := range blockEvents {
		blockQueryArgs := []SaleMomentBatch_Query_Args{}
		for _, event := range blockEvent.Events {
			momentEvent := MomentPurchasedListedChangedEvent(event.Value)
			blockQueryArgs = append(blockQueryArgs, SaleMomentBatch_Query_Args{
				OwnerAddress: *momentEvent.Seller(),
				MomentFlowID: momentEvent.Id()})
		}

		if len(blockQueryArgs) > 0 {
			batchResults, err := SaleMomentBatch_Query(process.QueryApi, blockEvent.Height-executeScriptAtOffset, blockQueryArgs)
			if err != nil {
				log.WithFields(log.Fields{
					"blockQueryArgs": blockQueryArgs,
				}).Warnf("SaleMomentBatch_Query %v", err)
			} else {
				for _, result := range batchResults {
					saleMoment := result.MomentDetails()
					if saleMoment != nil {
						callBack(&blockEvent, findMomentPurchasedListedChangedEvent(blockEvent.Events, result.MomentID()), result)
					} else {
						log.WithFields(log.Fields{
							"result.MomentID()": result.MomentID(),
						}).Warn("Unable to get sale details")
					}
				}
			}
		}
	}
}

func findMomentPurchasedListedChangedEvent(events []flow.Event, momentIdToFind uint64) MomentPurchasedListedChangedEvent {
	for _, event := range events {
		momentEvent := MomentPurchasedListedChangedEvent(event.Value)
		if momentEvent.Id() == momentIdToFind {
			return momentEvent
		}
	}
	panic("internal error findMomentPurchasedListedChangedEvent")
}
