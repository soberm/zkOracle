package zkOracle

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Index           uint64
	Registered      bool
	BindAddress     string
	PrivateKey      string
	ContractAddress string
	Ethereum        EthereumConfig
	ZKP             ProofConfig
}

type ContractsConfig struct {
	RegistryContractAddress string
	OracleContractAddress   string
	DistKeyContractAddress  string
}

type EthereumConfig struct {
	SourceAddress string
	TargetAddress string
	PrivateKey    string
}

type ProofConfig struct {
	ProvingKey   string
	VerifyingKey string
	R1CS         string
}

func SetConfigDefaults(v *viper.Viper) {

}

func LoadConfig(v *viper.Viper, configFile string) error {
	v.SetConfigFile(configFile)
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("read in config: %w", err)
	}
	return nil
}
