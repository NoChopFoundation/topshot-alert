// Inspired by https://medium.com/@eric.ren_51534/polling-nba-top-shot-p2p-market-purchase-events-from-flow-blockchain-using-flow-go-sdk-3ec80119e75f
// The main difference is we batch multiple requests to reduce round-trip calls.
package topshot

import (
	"fmt"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
)

type SaleMomentBatch_Query_Args struct {
	OwnerAddress flow.Address
	MomentFlowID uint64
}

func SaleMomentBatch_Query(queryApi *TopshotQueryApi, blockHeight uint64, args []SaleMomentBatch_Query_Args) ([]SaleMomentBatch, error) {
	script := `
		import TopShot from 0x0b2a3299cc857e29
		import Market from 0xc1e4f4f4c4257510

        pub struct SaleMoment {
			pub var id: UInt64
			pub var playId: UInt32
			pub var play: {String: String}
			pub var setId: UInt32
			pub var setName: String
			pub var serialNumber: UInt32
			pub var price: UFix64
			init(moment: &TopShot.NFT, price: UFix64) {
			  self.id = moment.id
			  self.playId = moment.data.playID
			  self.play = TopShot.getPlayMetaData(playID: self.playId)!
			  self.setId = moment.data.setID
			  self.setName = TopShot.getSetName(setID: self.setId)!
			  self.serialNumber = moment.data.serialNumber
			  self.price = price
			}
		}

		pub struct SaleMomentBatch {
			pub var owner: Address
			pub var momentID: UInt64
			pub var saleMoment: SaleMoment?
			init(owner: Address, momentID: UInt64, saleMoment: SaleMoment?) {
				self.owner = owner
				self.momentID = momentID
				self.saleMoment = saleMoment
			}
		}

		pub fun main(ownerInputs: [Address], momentIDInputs: [UInt64]): [SaleMomentBatch] {
			var batchResults: [SaleMomentBatch] = []
			var i = 0
			while i < ownerInputs.length {
				let acct = getAccount(ownerInputs[i])
				let collectionOpt = acct.getCapability(/public/topshotSaleCollection)!.borrow<&{Market.SalePublic}>() 
				if collectionOpt != nil {
					let collectionRef = collectionOpt!
					let saleMoment = SaleMoment(moment: collectionRef.borrowMoment(id: momentIDInputs[i])!,price: collectionRef.getPrice(tokenID: momentIDInputs[i])!)
					batchResults.append(SaleMomentBatch(owner: ownerInputs[i], momentID: momentIDInputs[i], saleMoment: saleMoment))
				} else {
					batchResults.append(SaleMomentBatch(owner: ownerInputs[i], momentID: momentIDInputs[i], saleMoment: nil))
				}
    			i = i + 1
			}
			return batchResults
		}
	`
	ownerAddressInputs := []cadence.Value{}
	for _, arg := range args {
		ownerAddressInputs = append(ownerAddressInputs, cadence.BytesToAddress(arg.OwnerAddress.Bytes()))
	}
	momentFlowIDInputs := []cadence.Value{}
	for _, arg := range args {
		momentFlowIDInputs = append(momentFlowIDInputs, cadence.UInt64(arg.MomentFlowID))
	}

	res, err := queryApi.ExecuteScriptAtBlockHeight(ExecuteScriptAtBlockHeight_Arg{
		BlockHeight:   blockHeight,
		CadenceScript: script,
		ScriptArgs: []cadence.Value{
			cadence.NewArray(ownerAddressInputs),
			cadence.NewArray(momentFlowIDInputs)}})
	if err != nil {
		return []SaleMomentBatch{}, err
	}

	batchResults := []SaleMomentBatch{}
	retValues := res.(cadence.Array)
	for _, retVal := range retValues.Values {
		batchResults = append(batchResults, SaleMomentBatch(retVal.(cadence.Struct)))
	}
	return batchResults, nil
}

type SaleMomentBatch cadence.Struct

func (s SaleMomentBatch) SellerAddress() *flow.Address {
	sellerAddress := flow.BytesToAddress((s.Fields[0].(cadence.Address)).Bytes())
	return &sellerAddress
}

func (s SaleMomentBatch) MomentID() uint64 {
	return uint64(s.Fields[1].(cadence.UInt64))
}

func (s SaleMomentBatch) MomentDetails() *SaleMoment {
	optional := (s.Fields[2]).(cadence.Optional)
	if saleMomentValue, ok := optional.Value.(cadence.Struct); ok {
		saleMoment := SaleMoment(saleMomentValue)
		return &saleMoment
	}
	return nil
}

// SaleMoment from
// Inspired by https://medium.com/@eric.ren_51534/polling-nba-top-shot-p2p-market-purchase-events-from-flow-blockchain-using-flow-go-sdk-3ec80119e75f
type SaleMoment cadence.Struct

func (s SaleMoment) ID() uint64 {
	return uint64(s.Fields[0].(cadence.UInt64))
}

func (s SaleMoment) PlayID() uint32 {
	return uint32(s.Fields[1].(cadence.UInt32))
}

func (s SaleMoment) SetName() string {
	return string(s.Fields[4].(cadence.String))
}

func (s SaleMoment) SetID() uint32 {
	return uint32(s.Fields[3].(cadence.UInt32))
}

func (s SaleMoment) Play() map[string]string {
	dict := s.Fields[2].(cadence.Dictionary)
	res := map[string]string{}
	for _, kv := range dict.Pairs {
		res[string(kv.Key.(cadence.String))] = string(kv.Value.(cadence.String))
	}
	return res
}

func (s SaleMoment) SerialNumber() uint32 {
	return uint32(s.Fields[5].(cadence.UInt32))
}

func (s SaleMoment) Price() float64 {
	return float64(s.Fields[6].(cadence.UFix64).ToGoValue().(uint64)) / 1e8 // ufixed 64 have 8 digits of precision
}

func (s SaleMoment) String() string {
	playData := s.Play()
	return fmt.Sprintf("saleMoment: serialNumber: %d, setID: %d, setName: %s, playID: %d, playerName: %s",
		s.SerialNumber(), s.SetID(), s.SetName(), s.PlayID(), playData["FullName"])
}
