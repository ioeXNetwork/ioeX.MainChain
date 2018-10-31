package blockchain

import (
	"time"

	"github.com/ioeXNetwork/ioeX.MainChain/config"

	"github.com/ioeXNetwork/ioeX.Utility/common"
)

var (
	OriginIssuanceAmount   = 20000 * 10000 * 100000000
	InflationPerYear       = OriginIssuanceAmount * 4 / 100
	BlockGenerateInterval  = int64(config.Parameters.ChainParam.TargetTimePerBlock / time.Second)
	GeneratedBlocksPerYear = 365 * 24 * 60 * 60 / BlockGenerateInterval
	RewardAmountPerBlock   = common.Fixed64(float64(InflationPerYear) / float64(GeneratedBlocksPerYear))
)
