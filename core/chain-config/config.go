// Copyright 2019, Keychain Foundation Ltd.
// This file is part of the dipperin-core library.
//
// The dipperin-core library is free software: you can redistribute
// it and/or modify it under the terms of the GNU Lesser General Public License
// as published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// The dipperin-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.


package chain_config

import (
	"fmt"
	"github.com/dipperin/dipperin-core/common"
	"github.com/dipperin/dipperin-core/common/util"
	"github.com/dipperin/dipperin-core/third-party/log"
	"github.com/dipperin/dipperin-core/third-party/log/ver_halt_check_log"
	"github.com/dipperin/dipperin-core/third-party/p2p/enode"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

const (
	AppName = "dipperin"
	Version = "1.0.0"

	BootEnvTagName = "boots_env"

	StaticBootNodesFileName         = "static_boot_nodes.json"
	StaticVerifierBootNodesFileName = "static_verifier_boot_nodes.json"

	MineProtocolVersion = 1
	CsProtocolVersion   = 1

	TestServer               = "10.200.0.139"
	TestVerifierBootNodePort = "10000"

	// 20M
	MaxBlockSize = 20 * 1024 * 1024
	MaxTxSize    = 512 * 1024
)

const (
	NodeTypeOfNormal = iota
	NodeTypeOfMineMaster
	NodeTypeOfVerifier
	NodeTypeOfVerifierBoot
)

var (
	bigOne = big.NewInt(1)
)

var config = defaultChainConfig()

func defaultChainConfig() *ChainConfig {
	c := &ChainConfig{
		//DeriveShaType:         DeriveShaTypeByHash,
		SupportHardwareWallet: false,
		ChainId:               big.NewInt(1),
		Version:               uint64(0),
		// verify segment size
		SlotSize: uint64(110),
		// verifier deposit lock period
		StakeLockSlot: uint64(4),
		// the interval of the Verify section from the election section
		SlotMargin: uint64(2),
		// number of verifier
		VerifierNumber: 22,
		//VerifierNumber:4,

		// angel verifier priority
		SystemVerifierPriority: 0,

		//mine conf
		//mining maximum difficulty value
		MainPowLimit: new(big.Int).Sub(new(big.Int).Lsh(bigOne, 253), bigOne),
		//average block generation duration
		BlockGenerate: uint64(13),
		//the block number in a difficulty adjust cycle
		BlockCountOfPeriod: uint64(4096),

		//verifier boot node number
		VerifierBootNodeNumber: 4,

		BlockTimeRestriction: 15*time.Second,
	}

	switch os.Getenv(BootEnvTagName) {
	case "mercury":
		c.NetworkID = 99
	case "test":
		c.NetworkID = 1
	}
	return c
}

type ChainConfig struct {
	//DeriveShaType int
	ChainId *big.Int
	// Version
	Version uint64
	// chain network id
	NetworkID uint64

	SupportHardwareWallet bool

	// db conf
	DatabaseHandles int `toml:"-"`
	DatabaseCache   int

	// elect conf
	// verify segment size
	SlotSize uint64
	// verifier deposit lock period
	StakeLockSlot uint64
	// the interval of the Verify section from the election section
	SlotMargin uint64
	// pbft verifier number
	VerifierNumber int

	//system verifier priority
	SystemVerifierPriority uint64

	//VerifierReward
	// Block reward for successfully mining a block
	//FrontierBlockReward *big.Int
	//// Block reward for successfully mining a block upward from Byzantium
	//ByzantiumBlockReward *big.Int

	//mine conf
	//mining maximum difficulty value
	MainPowLimit *big.Int
	//average block generation duration
	BlockGenerate uint64
	//the block number in a difficulty adjust cycle
	BlockCountOfPeriod uint64

	//verifier boot node number
	VerifierBootNodeNumber int

	//timeStamp restriction
	BlockTimeRestriction time.Duration
}

func GetChainConfig() *ChainConfig {
	return config
}

// Get the operating environment：test mercury
func GetCurBootsEnv() string {
	return os.Getenv("boots_env")
}

func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := util.HomeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", AppName)
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", AppName)
		} else {
			return filepath.Join(home, "."+AppName)
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

var (
	// roughly generate a block per 10s
	//GenesisDifficulty      = common.HexToDiff("0x1e010011")
	//GenesisDifficulty      = common.HexToDiff("0x1e040011")

	// roughly generate a block per 5s
	//GenesisDifficulty      = common.HexToDiff("0x1e077011")
	// roughly generate a block per 1~3s
	//GenesisDifficulty      = common.HexToDiff("0x1e17f011")
	// Produce block very quickly
	GenesisDifficulty = common.HexToDiff("0x1e566611")
	// Produce block very quickly
	//GenesisDifficulty = common.HexToDiff("0x1effffff")
)

// verifier boot nodes
var (
	VerifierBootNodes []*enode.Node
	KBucketNodes      []*enode.Node
)

func InitBootNodes(dataDir string) {
	log.Info("the boot env is:", "env", os.Getenv(BootEnvTagName))
	// If the environment variable is set during deploy use, these environment variables are automatically taken when the startup command is used.
	switch os.Getenv(BootEnvTagName) {
	case "test":
		//log.Agent("use test boot env", "boot server", TestServer, "v boot port", TestVerifierBootNodePort)
		initTestBoots(dataDir)
	case "mercury":
		//log.Agent("use mercury boot env")
		initMercuryBoots(dataDir)
	default:
		//log.Agent("use local boot env")
		log.Info("use local boot env")
		initLocalBoots(dataDir)
	}
	for _, vb := range VerifierBootNodes {
		log.Info("VerifierBootNodes", "vb", vb.String())
	}
	for _, kn := range KBucketNodes {
		log.Info("KBucketNodes", "vb", kn.String())
	}
}

func initTestBoots(dataDir string) {
	// verifier boot node
	if VerifierBootNodes = LoadVerifierBootNodesFromFile(dataDir); len(VerifierBootNodes) == 0 {
		n, _ := enode.ParseV4(fmt.Sprintf("enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@%v:%v", TestServer, 10003))
		VerifierBootNodes = append(VerifierBootNodes, n)
		n, _ = enode.ParseV4(fmt.Sprintf("enode://199cc6526cb63866dfa5dc81aed9952f2002b677560b6f3dc2a6a34a5576216f0ca25711c5b4268444fdef5fee4a01a669af90fd5b6049b2a5272b39c466b2ac@%v:%v", TestServer, 10006))
		VerifierBootNodes = append(VerifierBootNodes, n)
		n, _ = enode.ParseV4(fmt.Sprintf("enode://71112a581231af08a63d5a9079ea8dd690efd992f2cfbf98ad43697345de564441406133247d19c754c98051c64909c40db15094770a881a373ca1ff2f20bea2@%v:%v", TestServer, 10009))
		VerifierBootNodes = append(VerifierBootNodes, n)
		n, _ = enode.ParseV4(fmt.Sprintf("enode://07f3fdca9a07b048ea7d0cb642f69004e4fa5dd390888a9bb3e9fc382697c3634280cc8d327703b872d3711462da4aca96ee805069510375e7be2aded3dc5ad6@%v:%v", TestServer, 10012))
		VerifierBootNodes = append(VerifierBootNodes, n)
	}

	// k bucket boot node. Try to read from file, use default if there isn't the file
	if KBucketNodes = LoadBootNodesFromFile(dataDir); len(KBucketNodes) == 0 {
		n, _ := enode.ParseV4(fmt.Sprintf("enode://e53903ee0001e81f9328c8d0929cedbaf9b4f5b65b536df5f5dd65e5aa650cc059976250d6fcc62685e46e035b52e22801e97b06bc84d8fc4848037c128a7b22@%v:30301", TestServer))
		KBucketNodes = append(KBucketNodes, n)
	}
}

func initLocalBoots(dataDir string) {
	// Two miners are 50030, one miner is 50027
	if VerifierBootNodes = LoadVerifierBootNodesFromFile(dataDir); len(VerifierBootNodes) == 0 {
		//n, _ := enode.ParseV4(fmt.Sprintf("enode://b832f4f2fe19dbc5604766bbb268a6d0f7ce9ce381b034b262a92f0ad8283a1b5fa058dea5269b66fbb2014a24fa7198c6dc2d8c9cbac7a348258fc20702561f@127.0.0.1:%v", TestVerifierBootNodePort))
		n, _ := enode.ParseV4(fmt.Sprintf("enode://8907151bf3e091349bb01daeae9f9a0dcfdc3eabaafadd488794c00f6a880816a535d04c111d211f6580d56a836a4063fc3c5b4195ae2be4f7c2561d028fafd5@127.0.0.1:%v", TestVerifierBootNodePort))
		VerifierBootNodes = append(VerifierBootNodes, n)
	}

	ver_halt_check_log.Info("the VerifierBootNodes is:", "VerifierBootNodes", VerifierBootNodes)

	// local boot node
	if KBucketNodes = LoadBootNodesFromFile(dataDir); len(KBucketNodes) == 0 {
		n, _ := enode.ParseV4("enode://9fb88505dca8bbed224e03e9c0c27f724bf12551003576c79aae38307d8a4b68a88ad6afa38ae87a0aadaeea4940da2b859ffb0ce6b4b94252de1acc035d2396@127.0.0.1:30301")
		KBucketNodes = append(KBucketNodes, n)
	}
}

// load from file + static nodes
func initMercuryBoots(dataDir string) {
	// The difference here is that the boot of the mercury may be manually started by the external network, so need to support both the file and the add
	//VerifierBootNodes = LoadVerifierBootNodesFromFile(dataDir)
	VerifierBootNodes = append(VerifierBootNodes, NewMercuryVBoots()...)

	// The difference here is that the boot of the mercury may be manually started by the external network, so need to support both the file and the add
	//KBucketNodes = LoadBootNodesFromFile(dataDir)
	KBucketNodes = append(KBucketNodes, mercuryKBoots()...)
}

func LoadBootNodesFromFile(dataDir string) (bootNodes []*enode.Node) {
	return LoadNodesFromFile(filepath.Join(dataDir, StaticBootNodesFileName))
}

func LoadVerifierBootNodesFromFile(dataDir string) (vBootNodes []*enode.Node) {
	return LoadNodesFromFile(filepath.Join(dataDir, StaticVerifierBootNodesFileName))
}

func LoadNodesFromFile(fileP string) (bootNodes []*enode.Node) {
	data, err := ioutil.ReadFile(fileP)
	if err != nil {
		log.Debug("load boot nodes from file failed", "err", err)
		return
	}

	var nodesStr []string
	if err = util.ParseJsonFromBytes(data, &nodesStr); err != nil {
		log.Debug("can't parse boot nodes", "err", err)
		return
	}

	for _, nStr := range nodesStr {
		if node, err := enode.ParseV4(nStr); err != nil {
			log.Debug("parse boot node failed", "err", err)
		} else {
			bootNodes = append(bootNodes, node)
		}
	}
	log.Debug("load boot nodes from file", "nodes len", len(bootNodes))
	return
}
