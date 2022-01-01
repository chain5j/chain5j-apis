// Package apis
//
// @author: xwc1125
package apis

import (
	"context"

	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/chain5j-protocol/models"
	"github.com/chain5j/chain5j-protocol/models/statetype"
)

type API struct {
	api *apis
}

func newAPI(api *apis) *API {
	return &API{
		api: api,
	}
}

func (a *API) SendRawTransaction(ctx context.Context, txType types.TxType, rawTx hexutil.Bytes) (types.Hash, error) {
	return a.api.SendRawTransaction(ctx, txType, rawTx)
}
func (a *API) GetTransaction(ctx context.Context, hash types.Hash) models.Transaction {
	return a.api.GetTransaction(ctx, hash)
}
func (a *API) GetTransactionReceipt(ctx context.Context, hash types.Hash) (models.Transaction, error) {
	return a.api.GetTransactionReceipt(ctx, hash)
}
func (a *API) GetTransactionLogs(ctx context.Context, hash types.Hash) ([]*statetype.Log, error) {
	return a.api.GetTransactionLogs(ctx, hash)
}
