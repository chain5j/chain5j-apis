// Package apis
//
// @author: xwc1125
package apis

import (
	"fmt"
	"github.com/chain5j/chain5j-protocol/protocol"
)

type option func(f *apis) error

func apply(f *apis, opts ...option) error {
	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(f); err != nil {
			return fmt.Errorf("option apply err:%v", err)
		}
	}
	return nil
}
func WithBlockDB(blockDB protocol.DatabaseReader) option {
	return func(f *apis) error {
		f.blockDB = blockDB
		return nil
	}
}

func WithStateDB(stateDB protocol.DatabaseReader) option {
	return func(f *apis) error {
		f.stateDB = stateDB
		return nil
	}
}

func WithTxPools(f protocol.APIs, txPools protocol.TxPools) {
	if a, ok := f.(*apis); ok {
		a.txPools = txPools
	}
}
func WithBlockReader(f protocol.APIs, blockReader protocol.BlockReader) {
	if a, ok := f.(*apis); ok {
		a.blockReader = blockReader
	}
}
