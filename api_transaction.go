// Package apis
//
// @author: xwc1125
package apis

import (
	"context"
	"github.com/chain5j/chain5j-pkg/codec"
	"github.com/chain5j/chain5j-pkg/types"
	"github.com/chain5j/chain5j-pkg/util/hexutil"
	"github.com/chain5j/chain5j-protocol/models"
	"github.com/chain5j/chain5j-protocol/models/statetype"
)

// SendRawTransaction 发送签名后的交易
func (api *apis) SendRawTransaction(ctx context.Context, txType types.TxType, rawTx hexutil.Bytes) (types.Hash, error) {
	tx, err := models.NewTransaction(txType)
	if err != nil {
		return types.Hash{}, err
	}

	// 进行编码转换
	err = codec.Coder().Decode(rawTx, tx)
	if err != nil {
		return types.Hash{}, err
	}

	err = api.txPools.Add(nil, tx)
	if err != nil {
		api.log.Debug("SendRawTransaction Err", "err", err)
		return types.Hash{}, err
	}

	return tx.Hash(), nil
}

// GetTransaction 根据交易Hash获取交易
func (api *apis) GetTransaction(ctx context.Context, hash types.Hash) models.Transaction {
	tx, blockHash, blockNumber, index, _ := api.blockDB.GetTransaction(hash)
	if tx == nil {
		return nil
	}
	rtx := &RPCTransaction{
		Transaction:      tx,
		BlockHash:        blockHash,
		BlockNumber:      hexutil.Uint64(blockNumber),
		TransactionIndex: hexutil.Uint(index),
		Type:             tx.TxType(),
	}

	return rtx
}

// GetTransactionReceipt 根据交易Hash获取交易收据
func (api *apis) GetTransactionReceipt(ctx context.Context, hash types.Hash) (models.Transaction, error) {
	tx, blockHash, blockNumber, index, _ := api.blockDB.GetTransaction(hash)
	if tx == nil {
		return nil, nil
	}

	receipts, _ := api.blockDB.GetReceipts(blockHash, blockNumber)

	fields := map[string]interface{}{
		"blockHash":        blockHash,
		"blockNumber":      hexutil.Uint64(blockNumber),
		"transactionHash":  hash,
		"transactionIndex": hexutil.Uint64(index),
		//"from":             tx.From(),
		//"to":               tx.To(),
		"contractAddress": nil,
	}
	if receipts != nil {
		receipt := receipts[index]
		fields["gasUsed"] = hexutil.Uint64(receipt.GasUsed)
		fields["cumulativeGasUsed"] = hexutil.Uint64(receipt.CumulativeGasUsed)
		fields["contractAddress"] = nil
		fields["logs"] = receipt.Logs
		fields["logsBloom"] = receipt.LogsBloom
		fields["status"] = hexutil.Uint(receipt.Status)

		if receipt.Logs == nil {
			fields["logs"] = [][]*statetype.Log{}
		}
		// If the ContractAddress is 20 0x0 bytes, assume it is not a contract creation
		if receipt.ContractAddress != (types.Address{}) {
			fields["contractAddress"] = receipt.ContractAddress
		}
	}

	return nil, nil
}

// GetTransactionLogs 根据交易Hash获取状态日志
func (api *apis) GetTransactionLogs(ctx context.Context, hash types.Hash) ([]*statetype.Log, error) {
	tx, blockHash, blockNumber, index, _ := api.blockDB.GetTransaction(hash)
	if tx == nil {
		return nil, nil
	}

	receipts, _ := api.blockDB.GetReceipts(blockHash, blockNumber)
	if receipts == nil {
		return nil, nil
	}
	receipt := receipts[index]

	if receipt.Logs == nil {
		return []*statetype.Log{}, nil
	}
	return receipt.Logs, nil
}
