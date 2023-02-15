package zkOracle

type Config struct {
	Index           uint64
	BindAddress     string
	PrivateKey      string
	ContractAddress string
	Ethereum        EthereumConfig
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
