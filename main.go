package main

import (
	"os"
	"runtime"

	"github.com/ioeXNetwork/ioeX.MainChain/blockchain"
	"github.com/ioeXNetwork/ioeX.MainChain/config"
	"github.com/ioeXNetwork/ioeX.MainChain/log"
	"github.com/ioeXNetwork/ioeX.MainChain/node"
	"github.com/ioeXNetwork/ioeX.MainChain/pow"
	"github.com/ioeXNetwork/ioeX.MainChain/protocol"
	"github.com/ioeXNetwork/ioeX.MainChain/servers"
	"github.com/ioeXNetwork/ioeX.MainChain/servers/httpjsonrpc"
	"github.com/ioeXNetwork/ioeX.MainChain/servers/httpnodeinfo"
	"github.com/ioeXNetwork/ioeX.MainChain/servers/httprestful"
	"github.com/ioeXNetwork/ioeX.MainChain/servers/httpwebsocket"
	"github.com/ioeXNetwork/ioeX.Utility/common"
)

const (
	DefaultMultiCoreNum = 4
)

func init() {
	log.Init(
		config.Parameters.PrintLevel,
		config.Parameters.MaxPerLogSize,
		config.Parameters.MaxLogsSize,
	)
	var coreNum int
	if config.Parameters.MultiCoreNum > DefaultMultiCoreNum {
		coreNum = int(config.Parameters.MultiCoreNum)
	} else {
		coreNum = DefaultMultiCoreNum
	}
	log.Debug("The Core number is ", coreNum)

	foundationAddress := config.Parameters.ChainParam.FoundationAddress1
	address, err := common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress1
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress1 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress2
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress2 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress3
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress3 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress4
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress4 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress5
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress5 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress6
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress6 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress7
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress7 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress8
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress8 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress9
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress9 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress10
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress10 = *address

	foundationAddress = config.Parameters.ChainParam.FoundationAddress11
	address, err = common.Uint168FromAddress(foundationAddress)
	if err != nil {
		log.Error(err.Error())
		os.Exit(-1)
	}
	blockchain.FoundationAddress11 = *address

	runtime.GOMAXPROCS(coreNum)
}

func startConsensus() {
	servers.LocalPow = pow.NewPowService()
	if config.Parameters.PowConfiguration.AutoMining {
		log.Info("Start POW Services")
		go servers.LocalPow.Start()
	}
}

func main() {
	//var blockChain *ledger.Blockchain
	var err error
	var noder protocol.Noder
	log.Info("Node version: ", config.Version)
	log.Info("1. BlockChain init")
	chainStore, err := blockchain.NewChainStore()
	if err != nil {
		goto ERROR
	}
	defer chainStore.Close()

	err = blockchain.Init(chainStore)
	if err != nil {
		goto ERROR
	}

	log.Info("2. Start the P2P networks")
	noder = node.InitLocalNode()

	servers.ServerNode = noder

	log.Info("3. --Start the RPC service")
	go httpjsonrpc.StartRPCServer()

	noder.WaitForSyncFinish()
	go httprestful.StartServer()
	go httpwebsocket.StartServer()
	if config.Parameters.HttpInfoStart {
		go httpnodeinfo.StartServer()
	}
	startConsensus()
	select {}
ERROR:
	log.Error(err)
	os.Exit(-1)
}
