// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"testing"

	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/snow"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
)

func TestMempoolAddTx(t *testing.T) {
	require := require.New(t)
	m, err := NewMempool(&snow.Context{}, prometheus.NewRegistry(), 5_000, nil)
	require.NoError(err)

	txs := make([]*GossipAtomicTx, 0)
	for i := 0; i < 3_000; i++ {
		tx := &GossipAtomicTx{
			Tx: &Tx{
				UnsignedAtomicTx: &TestUnsignedTx{
					IDV: ids.GenerateTestID(),
				},
			},
		}

		txs = append(txs, tx)
	}
	errs := m.Add(txs...)
	for _, err := range errs {
		require.NoError(err)
	}

	for _, tx := range txs {
		require.True(m.bloom.Has(tx))
	}
}

// Add should return an error if a tx is already known
func TestMempoolAdd(t *testing.T) {
	require := require.New(t)
	m, err := NewMempool(&snow.Context{}, prometheus.NewRegistry(), 5_000, nil)
	require.NoError(err)

	tx := &GossipAtomicTx{
		Tx: &Tx{
			UnsignedAtomicTx: &TestUnsignedTx{
				IDV: ids.GenerateTestID(),
			},
		},
	}

	require.NoError(m.Add(tx)[0])
	err = m.Add(tx)[0]
	require.ErrorIs(err, errTxAlreadyKnown)
}
