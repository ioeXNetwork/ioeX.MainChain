package blockchain

import (
	"time"

	"github.com/ioeXNetwork/ioeX.MainChain/config"
	"github.com/ioeXNetwork/ioeX.Utility/common"
)

var (
	OriginIssuanceAmount           = 20000 * 10000 * 100000000
	InflationPerYear               = OriginIssuanceAmount * 3 / 100
	BlockGenerateInterval          = int64(config.Parameters.ChainParam.TargetTimePerBlock / time.Second)
	GeneratedBlocksPerYear         = 365 * 24 * 60 * 60 / BlockGenerateInterval
	TotalRewardAmountPerBlock      = common.Fixed64(float64(24 * 100000000))
	MinerRewardAmountPerBlock      = common.Fixed64(float64(4 * 100000000))
	FoundationRewardAmountPerBlock = common.Fixed64(float64(20 * 100000000))
)
