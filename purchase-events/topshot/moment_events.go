/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

// https://github.com/dapperlabs/nba-smart-contracts/blob/master/contracts/MarketTopShot.cdc
// pub event MomentPriceChanged(id: UInt64, newPrice: UFix64, seller: Address?)
// pub event MomentPurchased(id: UInt64, price: UFix64, seller: Address?)
// pub event MomentListed(id: UInt64, price: UFix64, seller: Address?)
// MomentPurchasedListedChangedEvent covers all three events as they have same arguments
type MomentPurchasedListedChangedEvent cadence.Event

// pub event MomentWithdrawn(id: UInt64, owner: Address?)
type MomentWithdrawnEvent cadence.Event

// pub event CutPercentageChanged(newPercent: UFix64, seller: Address?)
type MomentCutPercentageChangedEvent cadence.Event

func (evt MomentPurchasedListedChangedEvent) Id() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt MomentPurchasedListedChangedEvent) Price() float64 {
	return float64(evt.Fields[1].(cadence.UFix64).ToGoValue().(uint64)) / 1e8 // ufixed 64 have 8 digits of precision
}

func (evt MomentPurchasedListedChangedEvent) Seller() *flow.Address {
	return grabAddress(evt.Fields[2])
}

func (evt MomentWithdrawnEvent) Id() uint64 {
	return uint64(evt.Fields[0].(cadence.UInt64))
}

func (evt MomentWithdrawnEvent) Owner() *flow.Address {
	return grabAddress(evt.Fields[1])
}

func (evt MomentCutPercentageChangedEvent) NewPercent() float64 {
	return float64(evt.Fields[0].(cadence.UFix64).ToGoValue().(uint64)) / 1e8 // ufixed 64 have 8 digits of precision
}

func (evt MomentCutPercentageChangedEvent) Seller() *flow.Address {
	return grabAddress(evt.Fields[1])
}

func grabAddress(val cadence.Value) *flow.Address {
	optionalAddress := (val).(cadence.Optional)
	if cadenceAddress, ok := optionalAddress.Value.(cadence.Address); ok {
		addr := flow.BytesToAddress(cadenceAddress.Bytes())
		return &addr
	}
	return nil
}
