// Package apis
//
// @author: xwc1125
package apis

import (
	"context"
	"github.com/chain5j/chain5j-protocol/protocol"
	"github.com/chain5j/logger"
)

var (
	_ protocol.APIs = new(apis)
)

// apis core 提供的API
type apis struct {
	log logger.Logger

	txPools     protocol.TxPools
	blockDB     protocol.DatabaseReader
	stateDB     protocol.DatabaseReader
	blockReader protocol.BlockReader

	apis []protocol.API
}

// NewApis 创建一个新的Api，提供给rpc调用
func NewApis(rootCtx context.Context, opts ...option) (protocol.APIs, error) {
	a := &apis{
		log:  logger.New("apis"),
		apis: make([]protocol.API, 0),
	}
	if err := apply(a, opts...); err != nil {
		a.log.Error("apply is error", "err", err)
		return nil, err
	}
	a.getApis()
	return a, nil
}

func (api *apis) getApis() {
	api.apis = append(api.apis, protocol.API{
		Namespace: "apps",
		Version:   "1.0",
		Service:   newAPI(api),
		Public:    true,
	})
}

func (api *apis) APIs() []protocol.API {
	return api.apis
}

func (api *apis) RegisterAPI(apis []protocol.API) {
	api.apis = append(api.apis, apis...)
}
