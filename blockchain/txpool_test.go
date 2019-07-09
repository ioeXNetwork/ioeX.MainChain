package blockchain

import (
	"bytes"
	"encoding/hex"
	"os"
	"testing"

	"github.com/ioeXNetwork/ioeX.MainChain/config"
	"github.com/ioeXNetwork/ioeX.MainChain/core"
	"github.com/ioeXNetwork/ioeX.MainChain/errors"
	"github.com/ioeXNetwork/ioeX.MainChain/log"

	"github.com/ioeXNetwork/ioeX.Utility/common"
	"github.com/stretchr/testify/assert"
)

var txPool TxPool

func TestTxPoolInit(t *testing.T) {
	log.Init(
		config.Parameters.PrintLevel,
		config.Parameters.MaxPerLogSize,
		config.Parameters.MaxLogsSize,
	)
	foundation, err := common.Uint168FromAddress("8VYXVxKKSAxkmRrfmGpQR2Kc66XhG6m3ta")
	if !assert.NoError(t, err) {
		return
	}
	FoundationAddress = *foundation

	chainStore, err := newTestChainStore()
	if err != nil {
		t.Fatal("open LedgerStore err:", err)
		os.Exit(1)
	}

	err = Init(chainStore)
	if err != nil {
		t.Fatal(err, "BlockChain generate failed")
	}

	txPool.Init()
}

func TestTxPool_AppendToTxnPool(t *testing.T) {
	tx := new(core.Transaction)
	txBytes, _ := hex.DecodeString("000403454c41010008803e6306563b26de010" +
		"000000000000000000000000000000000000000000000000000000000000000ffff" +
		"ffffffff02b037db964a231458d2d6ffd5ea18944c4f90e63d547c5d3b9874df66a" +
		"4ead0a39becdc01000000000000000012c8a2e0677227144df822b7d9246c58df68" +
		"eb11ceb037db964a231458d2d6ffd5ea18944c4f90e63d547c5d3b9874df66a4ead" +
		"0a3c1d258040000000000000000129e9cf1c5f336fcf3a6c954444ed482c5d916e5" +
		"06dd00000000")
	tx.Deserialize(bytes.NewReader(txBytes))
	errCode := txPool.AppendToTxnPool(tx)
	assert.Equal(t, errCode, errors.ErrIneffectiveCoinbase)

}
