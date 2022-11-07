package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/DecentralCardGame/cardchainClient/cardchain/client"
)

var config *faucetConfig
var once sync.Once

type faucetConfig struct {
	ChainHome, RPCNode, BlockchainUser, SecretKey string
}

func (f faucetConfig) verify() error {
	envs := []string{"CHAIN_HOME", "RPC_NODE", "BLOCKCHAIN_USER", "SECRET_KEY"}
	for idx, val := range []string{f.ChainHome, f.RPCNode, f.BlockchainUser, f.SecretKey} {
		if val == "" {
			return fmt.Errorf("Env var '%s' isn't set properly", envs[idx])
		}
	}
	return nil
}

func (f faucetConfig) ClientConfig() client.Config {
	return client.Config{
		f.ChainHome,
		f.RPCNode,
	}
}

func Config() faucetConfig {
	return *config
}

func FromEnv() error {
	chainHome := os.Getenv("CHAIN_HOME")
	if chainHome == "" {
		chainHome = "~/.Cardchain"
	}

	f := faucetConfig{
		chainHome,
		os.Getenv("RPC_NODE"),
		os.Getenv("BLOCKCHAIN_USER"),
		os.Getenv("SECRET_KEY"),
	}

	once.Do(func() {
		config = &f
	})

	return f.verify()
}
