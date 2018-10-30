package main

import (
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
	"os"
	"runtime"
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

	blockchain.FoundationAddress = config.Parameters.Configuration.FoundationAddress

	if blockchain.FoundationAddress == "" {
		blockchain.FoundationAddress = "8VYXVxKKSAxkmRrfmGpQR2Kc66XhG6m3ta"
	}
	runtime.GOMAXPROCS(coreNum)
}

func startConsensus() {
	servers.LocalPow = pow.NewPowService("logPow")
	if config.Parameters.PowConfiguration.AutoMining {
		log.Info("Start POW Services")
		go servers.LocalPow.Start()
	}
}

func main() {
	//var blockChain *ledger.Blockchain
	var err error
	var noder protocol.Noder
	log.Trace("Node version: ", config.Version)
	log.Info("1. BlockChain init")
	chainStore, err := blockchain.NewChainStore()
	if err != nil {
		log.Fatal("open LedgerStore err:", err)
		os.Exit(1)
	}
	defer chainStore.Close()

	err = blockchain.Init(chainStore)
	if err != nil {
		log.Fatal(err, "BlockChain generate failed")
		goto ERROR
	}

	log.Info("2. Start the P2P networks")
	noder = node.InitLocalNode()
	noder.WaitForSyncFinish()

	servers.NodeForServers = noder
	startConsensus()

	log.Info("3. --Start the RPC service")
	go httpjsonrpc.StartRPCServer()
	go httprestful.StartServer()
	go httpwebsocket.StartServer()
	if config.Parameters.HttpInfoStart {
		go httpnodeinfo.StartServer()
	}
	select {}
ERROR:
	os.Exit(1)
}
