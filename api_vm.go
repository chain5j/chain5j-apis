// Package apis
//
// @author: xwc1125
package apis

import (
	"context"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/chain5j-protocol/models"
	"github.com/chain5j/chain5j-protocol/protocol"
)

func (api *apis) GetCode(ctx context.Context, contract types.Address, blockHeight hexutil.Uint64) (*hexutil.Bytes, error) {
	panic("implement me")
}

func (api *apis) Call(ctx context.Context, hash models.VmMessage) (*hexutil.Bytes, error) {
	panic("implement me")
}

func (api *apis) EstimateGas(ctx context.Context, transaction models.Transaction) (*hexutil.Uint64, error) {
	panic("implement me")
}

func (api *apis) CompileContract(ctx context.Context, compileType protocol.CompileType, contract hexutil.Bytes) (*hexutil.Bytes, error) {
	panic("implement me")
}
