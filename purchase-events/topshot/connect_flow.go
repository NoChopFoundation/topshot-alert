/**
 * author: NoChopFoundation@gmail.com
 */
package topshot

import (
	"context"

	"github.com/onflow/cadence"
	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/client"
	"google.golang.org/grpc"
)

const TopshotQueryApi_MomentPurchased_Type = "A.c1e4f4f4c4257510.Market.MomentPurchased"
const TopshotQueryApi_MomentListed_Type = "A.c1e4f4f4c4257510.Market.MomentListed"
const TopshotQueryApi_PriceChanged_Type = "A.c1e4f4f4c4257510.Market.MomentPriceChanged"

type GetMomentEvents_Arg struct {
	// The block height to begin looking for events
	StartHeight uint64
	// The block height to end looking for events (inclusive)
	EndHeight uint64
}

type ExecuteScriptAtBlockHeight_Arg struct {
	BlockHeight   uint64
	CadenceScript string
	ScriptArgs    []cadence.Value
}

type ExecuteScriptAtLatestBlock_Arg struct {
	CadenceScript string
	ScriptArgs    []cadence.Value
}

type GetLatestBlockFunc func() (*flow.Block, error)
type GetMomentPurchasedEvents_Func func(GetMomentEvents_Arg) ([]client.BlockEvents, error)
type GetMomentListedEvents_Func func(GetMomentEvents_Arg) ([]client.BlockEvents, error)
type GetMomentPriceChangedEvents_Func func(GetMomentEvents_Arg) ([]client.BlockEvents, error)

type ExecuteScriptAtBlockHeight_Func func(ExecuteScriptAtBlockHeight_Arg) (cadence.Value, error)
type ExecuteScriptAtLatestBlock_Func func(ExecuteScriptAtLatestBlock_Arg) (cadence.Value, error)

type TopshotQueryApi struct {
	// private
	FlowClient *client.Client //TODO
	ctx        context.Context

	// public
	GetLatestBlock              GetLatestBlockFunc
	GetMomentPurchasedEvents    GetMomentPurchasedEvents_Func
	GetMomentListedEvents       GetMomentListedEvents_Func
	GetMomentPriceChangedEvents GetMomentPriceChangedEvents_Func
	ExecuteScriptAtBlockHeight  ExecuteScriptAtBlockHeight_Func
	ExecuteScriptAtLatestBlock  ExecuteScriptAtLatestBlock_Func
}

// Wrapper for the APIs we use in the Flow SDK to reduce common parameters
// and to provide a possible cut-point for unit testing
func Connection(ctx context.Context, config *Configuration) (*TopshotQueryApi, error) {
	// Connect to flow
	flowClient, err := client.New(config.accessNode, grpc.WithInsecure())
	if err != nil {
		return nil, err
	} else {
		err = flowClient.Ping(ctx)
		if err != nil {
			return nil, err
		} else {
			return &TopshotQueryApi{
				FlowClient: flowClient,
				ctx:        ctx,
				GetLatestBlock: func() (*flow.Block, error) {
					return flowClient.GetLatestBlock(ctx, true)
				},
				GetMomentPriceChangedEvents: func(args GetMomentEvents_Arg) ([]client.BlockEvents, error) {
					return flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
						Type:        TopshotQueryApi_PriceChanged_Type,
						StartHeight: args.StartHeight,
						EndHeight:   args.EndHeight,
					})
				},
				GetMomentPurchasedEvents: func(args GetMomentEvents_Arg) ([]client.BlockEvents, error) {
					return flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
						Type:        TopshotQueryApi_MomentPurchased_Type,
						StartHeight: args.StartHeight,
						EndHeight:   args.EndHeight,
					})
				},
				GetMomentListedEvents: func(args GetMomentEvents_Arg) ([]client.BlockEvents, error) {
					return flowClient.GetEventsForHeightRange(ctx, client.EventRangeQuery{
						Type:        TopshotQueryApi_MomentListed_Type,
						StartHeight: args.StartHeight,
						EndHeight:   args.EndHeight,
					})
				},
				ExecuteScriptAtBlockHeight: func(args ExecuteScriptAtBlockHeight_Arg) (cadence.Value, error) {
					return flowClient.ExecuteScriptAtBlockHeight(ctx, args.BlockHeight, []byte(args.CadenceScript), args.ScriptArgs)
				},
				ExecuteScriptAtLatestBlock: func(args ExecuteScriptAtLatestBlock_Arg) (cadence.Value, error) {
					return flowClient.ExecuteScriptAtLatestBlock(ctx, []byte(args.CadenceScript), args.ScriptArgs)
				},
			}, err
		}
	}
}
