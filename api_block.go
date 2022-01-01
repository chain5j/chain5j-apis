// Package apis
//
// @author: xwc1125
package apis

import (
	"context"
	"github.com/chain5j/chain5j-pkg/network/rpc"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/chain5j-protocol/models"
	"github.com/chain5j/chain5j-protocol/models/eventtype"
	"github.com/chain5j/chain5j-protocol/models/statetype"
	"math/big"
)

func (api *apis) BlockHeight(ctx context.Context) (*hexutil.Uint64, error) {
	height := hexutil.Uint64(api.blockReader.CurrentBlock().Height())
	return &height, nil
}

func (api *apis) GetBlockByHash(ctx context.Context, blockHash types.Hash) (*models.Block, error) {
	block := api.blockReader.GetBlockByHash(blockHash)

	if block == nil {
		return nil, nil
	}

	response, err := RPCMarshalBlock(block)
	if err != nil {
		return nil, err
	}
	_ = response
	return nil, nil
}

func (api *apis) GetBlockByHeight(ctx context.Context, blockHeight hexutil.Uint64) (*models.Block, error) {
	block := api.blockReader.GetBlockByNumber(uint64(blockHeight))

	if block == nil {
		return nil, nil
	}

	response, err := RPCMarshalBlock(block)
	if err != nil {
		return nil, err
	}

	_ = response
	return nil, nil
}

func (api *apis) GetBlockTransactionCountByHash(ctx context.Context, blockHash types.Hash) (*hexutil.Uint64, error) {
	block := api.blockReader.GetBlockByHash(blockHash)

	if block == nil {
		return nil, nil
	}
	n := hexutil.Uint64(block.Transactions().Len())
	return &n, nil
}

func (api *apis) GetBlockTransactionCountByHeight(ctx context.Context, blockHeight hexutil.Uint64) (*hexutil.Uint64, error) {
	panic("implement me")
}

func (api *apis) GetTransactionByBlockHashAndIndex(ctx context.Context, blockHash types.Hash, txIndex hexutil.Uint64) (models.Transaction, error) {
	panic("implement me")
}

func (api *apis) GetTransactionByBlockHeightAndIndex(ctx context.Context, blockHeight hexutil.Uint64, txIndex hexutil.Uint64) (models.Transaction, error) {
	panic("implement me")
}

// RPCTransaction rpc响应的交易
type RPCTransaction struct {
	BlockHash          types.Hash     `json:"blockHash"`
	BlockNumber        hexutil.Uint64 `json:"blockNumber"`
	TransactionIndex   hexutil.Uint   `json:"transactionIndex"`
	Type               types.TxType   `json:"type"`
	models.Transaction `json:"transaction"`
}

// ReceiptsInBlock TODO
func (api *apis) ReceiptsInBlock(ctx context.Context, blockNr rpc.BlockNumber) (statetype.Receipts, error) {
	bhash, err := api.blockDB.GetCanonicalHash(uint64(blockNr))
	if err != nil {
		return nil, err
	}
	if bhash == (types.Hash{}) {
		return nil, nil
	}

	//log().Info("read receipt", "hash", hash)
	return api.blockDB.GetReceipts(bhash, uint64(blockNr))
}

// NonceHashInfo TODO
type NonceHashInfo struct {
	Nonce *big.Int `json:"nonce"`
	Hash  string   `json:"hash"`
}

// NonceHashList TODO
type NonceHashList []NonceHashInfo

// Len TODO
func (a NonceHashList) Len() int {
	return len(a)
}

// Less TODO
func (a NonceHashList) Less(i, j int) bool {
	if a[i].Nonce.Cmp(a[j].Nonce) < 0 {
		return true
	}
	return false
}

// Swap TODO
func (a NonceHashList) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// AddressStatus 地址状态
type AddressStatus struct {
	Pending NonceHashList `json:"pending"`
	Count   uint64        `json:"count"`
}

//// Status TODO
//func (api *apis) Status() TxPoolStatus {
//	return ni
//}

func RPCMarshalBlock(b *models.Block) (map[string]interface{}, error) {
	head := b.Header() // copies the header once
	fields := map[string]interface{}{
		"height":           head.Height,
		"hash":             b.Hash(),
		"parentHash":       head.ParentHash,
		"stateRoots":       head.StateRoots,
		"consensusName":    head.Consensus.Name,
		"consensus":        hexutil.Bytes(head.Consensus.Consensus),
		"extra":            hexutil.Bytes(head.Extra),
		"size":             hexutil.Uint64(b.Size()),
		"gasLimit":         hexutil.Uint64(head.GasLimit),
		"gasUsed":          hexutil.Uint64(head.GasUsed),
		"timestamp":        hexutil.Uint64(head.Timestamp),
		"transactionsRoot": head.TxsRoot,
	}

	txs := b.Transactions()
	transactions := make([]interface{}, txs.Len())
	for i, tx := range txs.Data() {
		transactions[i] = tx
	}
	fields["transactions"] = transactions

	return fields, nil
}

// NewHeads send a notification each time a new (header) block is appended to the chain.
func (api *apis) NewHeads(ctx context.Context) (*rpc.Subscription, error) {
	notifier, supported := rpc.NotifierFromContext(ctx)
	if !supported {
		return &rpc.Subscription{}, rpc.ErrNotificationsUnsupported
	}

	rpcSub := notifier.CreateSubscription()

	go func() {
		headers := make(chan eventtype.ChainHeadEvent)
		headersSub := api.blockReader.SubscribeChainHeadEvent(headers)

		for {
			select {
			case h := <-headers:
				notifier.Notify(rpcSub.ID, h.Block.Header())
			case <-rpcSub.Err():
				headersSub.Unsubscribe()
				return
			case <-notifier.Closed():
				headersSub.Unsubscribe()
				return
			}
		}
	}()

	return rpcSub, nil
}
