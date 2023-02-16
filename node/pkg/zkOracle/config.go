package zkOracle

type Config struct {
	Index           uint64
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
